package google

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const (
	configPath                                = "config"
	clientIDConfigPropertyName                = "client_id"
	clientSecretConfigPropertyName            = "client_secret"
	clientOAuthRedirectUrlPropertyName        = "redirect_url"
	clientFetchGroupsConfigPropertyName       = "fetch_groups"
	clientServiceAccountKeyConfigPropertyName = "service_acc_key"
	clientDelegationUserConfigPropertyName    = "delegation_user"
	configEntry                               = "config"
)

func (b *backend) pathConfigWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var (
		clientID       = data.Get(clientIDConfigPropertyName).(string)
		clientSecret   = data.Get(clientSecretConfigPropertyName).(string)
		redirectUrl    = data.Get(clientOAuthRedirectUrlPropertyName).(string)
		fetchGroups    = data.Get(clientFetchGroupsConfigPropertyName).(bool)
		serviceAccount = data.Get(clientServiceAccountKeyConfigPropertyName).(string)
		delegationUser = data.Get(clientDelegationUserConfigPropertyName).(string)
	)

	entry, err := logical.StorageEntryJSON(configEntry, config{
		ClientID:       clientID,
		ClientSecret:   clientSecret,
		RedirectUrl:    redirectUrl,
		FetchGroups:    fetchGroups,
		ServiceAccount: serviceAccount,
		DelegationUser: delegationUser,
	})
	if err != nil {
		return nil, err
	}

	return nil, req.Storage.Put(ctx, entry)
}

func (b *backend) pathConfigRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, nil
	}

	configMap := map[string]interface{}{
		clientIDConfigPropertyName:                config.ClientID,
		clientSecretConfigPropertyName:            config.ClientSecret,
		clientOAuthRedirectUrlPropertyName:        config.RedirectUrl,
		clientFetchGroupsConfigPropertyName:       config.FetchGroups,
		clientServiceAccountKeyConfigPropertyName: config.ServiceAccount,
		clientDelegationUserConfigPropertyName:    config.DelegationUser,
	}

	return &logical.Response{
		Data: configMap,
	}, nil
}

// Config returns the configuration for this backend.
func (b *backend) config(ctx context.Context, s logical.Storage) (*config, error) {
	entry, err := s.Get(ctx, configEntry)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result config
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, fmt.Errorf("error reading configuration: %s", err)
	}

	return &result, nil
}

type config struct {
	ClientID       string `json:"client_id"`
	ClientSecret   string `json:"client_secret"`
	RedirectUrl    string `json:"redirect_url"`
	FetchGroups    bool   `json:"fetch_groups"`
	ServiceAccount string `json:"service_acc_key"`
	DelegationUser string `json:"delegation_user"`
}

func (c *config) oauth2Config() *oauth2.Config {
	oauthRedirectUrl := c.RedirectUrl
	if len(strings.TrimSpace(oauthRedirectUrl)) == 0 {
		oauthRedirectUrl = "urn:ietf:wg:oauth:2.0:oob"
	}

	config := &oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  oauthRedirectUrl,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
	}
	return config
}
