# How to configure `vault-auth-google` and use it for G Suite accounts

Setting `vault-auth-google` up on G Suite accounts requires a larger set of
parameters and configurations, but extends the capabilities of policy bounding,
since you can bound policies not only to email addresses, but on Google Groups
and its members. This allow a more granular set of permissions: each group is
bound to a specific policy; to grant or revoke permissions on a given user, the
G Suite Admin only have to add ou remove that user to the group.


## Configuration

### Configuring the auth method with group bounding and web-based OAuth2 flow

To enable `vault-auth-google` to bound policies to Google groups and its
members, you must create a domain-wide Service Account on GCP for user
impersonation on G Suite.

Instructions on how to create an OAuth2 credential is available
[here](oauth.md). On how to create a Service Account key, check
[here](service-account.md).

The following parameters are required:

 - _(string)_ `client_id`: The Google OAuth Client ID.
 - _(string)_ `client_secret`: The Google OAuth Client secret.
 - _(boolean)_ `fetch_groups`: Should the plugin bound policies to groups? **true** if yes, **false** otherwise.
 - _(string_ `delegation_user`: The Google user that delegates the API permission.
 - _(string)_ `service_acc_key`: The content of the Service Account private key.

After creating the OAuth credential and the Service Account key, follow the steps:

- Read the Service Account key content and stores it inside an environment
  variable for easy reference.

```sh
SERVICE_ACCOUNT_KEY=$(cat /path/to/the/key.json)
```

- Write the parameters onto the plugin's configuration.

```sh
vault write auth/google/config \
    client_id="<GOOGLE_CLIENT_ID>" \
    client_secret="<GOOGLE_CLIENT_SECRET>" \
    redirect_url="https://domain.com/redirect" \
    fetch_groups=true \
    delegation_user="user@domain.com" \
    service_acc_key=$SERVICE_ACCOUNT_KEY
```

Bare in mind that the delegation user must be the same one referenced at the
moment the Service Account were created. The `delegation_user` parameter tells
`vault-auth-google` which user must be impersonated in order to Vault be able
to list your organization's groups. The `service_acc_key` will have the Service
Account private key's content, authorizing the impersonation.


## Usage

After configuring the auth method, a role, bounding a given email or group to a
policy, is required. The following parameters are expected to create a role:

 - _(string)_ `bound_domain`: A domain name bounding a G Suite organization to
     a given policy. When setting a bounded domain, the plugin expects that
     all email addresses (users or groups) are part of the domain.
 - _(string)_ `bound_emails`: A list of email addresses bounding users to a
     given policy.
 - _(string)_ `bound_groups`: A list of Google groups bounding its members to a
     given policy.
 - _(string)_ `policies`: The list of policies associated with the role.

### Creating a role bounding a policy to a G Suite group

The following snippet creates a role named `default`, bounding the G Suite
groups `infra@domain.com` and `dev@domain.com` to the policies
`read-only-repos` and `read-only-machines`.

```sh
vault write auth/google/role/default \
    bound_domain="domain.com" \
    bound_groups="infra@domain.com, dev@domain.com" \
    policies="read-only-repos, read-only-machines"
```
