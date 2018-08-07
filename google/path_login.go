package google

import (
	"errors"
	"fmt"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	directory "google.golang.org/api/admin/directory/v1"
	goauth "google.golang.org/api/oauth2/v2"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const (
	loginPath                   = "login"
	googleAuthCodeParameterName = "code"
	roleParameterName           = "role"
)

func (b *backend) pathLogin(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	code := data.Get(googleAuthCodeParameterName).(string)
	roleName := data.Get(roleParameterName).(string)
	role, err := b.role(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("role '%s' not found", roleName)), nil
	}

	config, err := b.config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return logical.ErrorResponse("missing config"), nil
	}

	googleConfig := config.oauth2Config()
	token, err := googleConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, err
	}

	user, groups, err := b.authenticate(config, token)
	if err != nil {
		return nil, err
	}

	policies, err := b.authorise(req.Storage, role, user, groups)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	encodedToken, err := encodeToken(token)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Auth: &logical.Auth{
			InternalData: map[string]interface{}{
				"token": encodedToken,
				"role":  roleName,
			},
			Policies: policies,
			Metadata: map[string]string{
				"username": user.Email,
				"domain":   user.Hd,
			},
			DisplayName: user.Email,
			LeaseOptions: logical.LeaseOptions{
				TTL:       role.TTL,
				Renewable: true,
			},
		},
	}, nil
}

func (b *backend) authRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	encodedToken, ok := req.Auth.InternalData["token"].(string)
	if !ok {
		return nil, errors.New("no refresh token from previous login")
	}

	roleName, ok := req.Auth.InternalData["role"].(string)
	if !ok {
		return nil, errors.New("no role name from previous login")
	}

	role, err := b.role(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, fmt.Errorf("role '%s' not found", roleName)
	}

	config, err := b.config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return logical.ErrorResponse("missing config"), nil
	}

	token, err := decodeToken(encodedToken)
	if err != nil {
		return nil, err
	}

	user, groups, err := b.authenticate(config, token)
	if err != nil {
		return nil, err
	}

	policies, err := b.authorise(req.Storage, role, user, groups)
	if err != nil {
		return nil, err
	}

	if !strSliceEquals(policies, req.Auth.Policies) {
		return logical.ErrorResponse(fmt.Sprintf("policies do not match. new policies: %s. old policies: %s.", policies, req.Auth.Policies)), nil
	}

	return framework.LeaseExtend(role.TTL, role.MaxTTL, b.System())(ctx, req, d)
}

func (b *backend) authenticate(config *config, token *oauth2.Token) (*goauth.Userinfoplus, []string, error) {
	client := config.oauth2Config().Client(context.Background(), token)
	userService, err := goauth.New(client)
	if err != nil {
		return nil, nil, err
	}

	user, err := goauth.NewUserinfoV2MeService(userService).Get().Do()
	if err != nil {
		return nil, nil, err
	}

	groups := []string{}
	if config.FetchGroups {
		scope := "https://www.googleapis.com/auth/admin.directory.group.readonly"

		serviceAccountCredential, err := google.JWTConfigFromJSON([]byte(config.ServiceAccount), scope)
		if err != nil {
			return nil, nil, err
		}

		serviceAccountCredential.Subject = config.DelegationUser
		saClient, err := directory.New(serviceAccountCredential.Client(context.Background()))
		if err != nil {
			return nil, nil, err
		}

		response, err := saClient.Groups.List().UserKey(user.Email).Do()
		for _, g := range response.Groups {
			groups = append(groups, g.Email)
		}
	}

	return user, groups, nil
}

func (b *backend) authorise(storage logical.Storage, role *role, user *goauth.Userinfoplus, groups []string) ([]string, error) {
	if user.Hd != role.BoundDomain && role.BoundDomain != "" {
		return nil, fmt.Errorf("user %s is not part of required domain %s, found %s", user.Email, role.BoundDomain, user.Hd)
	}

	// Is this user in one of the bound groups for this role?
	isGroupMember := strSliceHasIntersection(groups, role.BoundGroups)
	isUserMember := strSliceHasIntersection([]string{user.Email}, role.BoundEmails)

	if !isGroupMember && !isUserMember {
		return nil, fmt.Errorf("user is not allowed to use this role")
	}

	return role.Policies, nil
}
