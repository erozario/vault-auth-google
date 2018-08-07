# How to configure `vault-auth-google` and use it for Gmail accounts

Setting `vault-auth-google` up for Gmail accounts is quite simple and
straightforward, since the only kind of policy bound available is the email
address bounding.


## Configuration

The configuration expects the following two required parameters:

 - _(string)_ `client_id`: The Google OAuth Client ID.
 - _(string)_ `client_secret`: The Google OAuth Client secret.

Instructions on how to create an OAuth2 credential is available [here](oauth.md).

### Configuring the auth method with local OAuth2 flow

```sh
vault write auth/google/config \
    client_id=<GOOGLE_CLIENT_ID> \
    client_secret=<GOOGLE_CLIENT_SECRET>
```


### Configuring the auth method with web-based OAuth2 flow

To configure a web-based flow, set the `redirect_url` parameter with the
desired callback URL.

```sh
vault write auth/google/config \
    client_id="<GOOGLE_CLIENT_ID>" \
    client_secret="<GOOGLE_CLIENT_SECRET>" \
    redirect_url="https://domain.com/redirect"
```


## Usage

After configuring the auth method, a role, bounding a given email to a policy,
is required. The following parameters are expected to create a role:

 - _(string)_ `bound_emails`: A list of email addresses bounding users to a
     given policy.
 - _(string)_ `policies`: The list of policies associated with the role.

### Creating a role bounding a policy to a Gmail account

The following snippet creates a role named `default`, bounding the email
`user@gmail.com`  to the policy `my-policy`. For multiple association,
delimiter the email addresses by comma.

```sh
vault write auth/google/role/default \
    bound_emails="user@gmail.com" \
    policies="my-policy"
```
