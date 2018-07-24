# Vault Auth Google

A Vault plugin for authenticating and receiving policies via Google Accounts
(Gmail and G Suite).


## Table of Contents

 - [Compatibility Matrix](#compatibility-matrix)
 - [Getting](#getting)
    - [Binary](#binary)
    - [From source](#from-source)
 - [Google API requirements](#google-api-requirements)
 - [Installation](#installation)
 - [Configuration](#configuration)
 - [Usage](#usage)
    - [Creating roles](#creating-roles)
       - [Gmail](#gmail)
       - [G Suite](#g-suite)
    - [Google's OAuth URL](#googles-oauth-url)
 - [Local environment setup](#local-environment-setup)
 - [License](#license)


## Compatibility Matrix

| Plugin Version | Vault 0.9.x  | Vault 0.10.x |
|----------------|--------------|--------------|
| 0.1.0          | :thumbsdown: | :thumbsup:   |


## Getting

After downloading the binary, move it into the [plugin
directory](https://www.vaultproject.io/guides/operations/plugin-backends.html)
configured at the Vault's configuration file. You can get the plugin binary via
the release page or building yourself.


### Binary

The pre-compiled binary is available at the [release](https://github.com/erozario/vault-auth-google/releases) page.


### From source

Clone this repository and build via `make`:

```sh
make all
```

The make recipe requires [dep](https://github.com/golang/dep) in order to get the project's dependencies.

Alternatively, you can get the dependencies via `go get <dependency>`.


## Google API requirements

In order to authenticate with your Google or G Suite account,
`vault-auth-google` requires an OAuth client ID and secret. You can generate
the client ID and secret at the Google Cloud Console, on the [credential
section](https://console.cloud.google.com/apis/credentials).

For a local oauth flow, set the credential type as *Other*. For an online
authentication flow (that requires redirection), set as *Web application*.

For bounding policies with G Suite email groups, the plugin requires access to
the [Admin SDK](https://developers.google.com/admin-sdk/directory/v1/guides/authorizing).


## Installation

* Calculate and register the SHA256 sum of the plugin in Vault's plugin catalog:

```sh
SHASUM=$(shasum -a 256 "/path/to/vault-auth-google" | cut -d " " -f1)
vault write sys/plugins/catalog/vault-auth-google \
  sha_256="$SHASUM" \
  command="vault-auth-google"
```

* Mount the auth method:

```sh
vault auth enable \
    -path="google" \
    -plugin-name="vault-auth-google" plugin
```


## Configuration

The plugin is set to receive four parameters, two of which are required:

 - _(string)_ `client_id` __*__: The Google OAuth client id.
 - _(string)_ `client_secret` __*__: The Google OAuth client secret.
 - _(boolean)_ `fetch_groups`: Should the plugin bound policies to groups? `true` if yes, `false` otherwise.
 - _(string)_ `redirect_url`: The URL that google will redirect after the OAuth
     flow. This URL should also be added at the credentials authorized URIs..

__* Required parameters__


Configuring the auth method:

```sh
vault write auth/google/config \
    client_id=<GOOGLE_CLIENT_ID> \
    client_secret=<GOOGLE_CLIENT_SECRET> \
    redirect_url="https://domain.com/redirect" \
    fetch_groups=true
```


## Usage

### Creating roles

 - _(string)_ `bound_emails`: The list of emails bound to the role.
 - _(string)_ `bound_domain`: The domain name bound to the role. When
     defining a bounded domain, the plugin expects that the emails are part of
     the domain.
 - _(string)_ `bound_groups`: The google group email bound the to role. When
     bounding a group to a role, every user within the group will be assigned to
     the policy.
 - _(string)_ `policies`: The policy name associated with the role.

#### Gmail

Creating a role to a Gmail account:
```sh
vault write auth/google/role/default \
    bound_emails=user@gmail.com \
    policies=default
```

#### G Suite

Creating a role to a G Suite account, bounding a groups:
```sh
vault write auth/google/role/default \
    bound_domain=<DOMAIN> \
    bound_groups=sec@<DOMAIN>,infra@<DOMAIN> \
    policies=default
```


### Google's OAuth URL

* Login using Google credentials

```sh
firefox $(vault read -field=url auth/google/code_url)
vault write auth/google/login code=$GOOGLE_CODE role=default
```


## Local environment setup

* Clone this repo

```sh
git clone git@github.com:erozario/vault-auth-google.git
```

* Create a temporary directory to compile the plugin into and to use as the plugin directory for Vault:

```sh
mkdir -p /tmp/vault-plugins
```

* Compile the plugin into the temporary directory:

```sh
go build -o /tmp/vault-plugins/vault-auth-google
```

* Create a configuration file to point Vault at this plugin directory:

```sh
tee /tmp/vault.hcl <<EOF
plugin_directory = "/tmp/vault-plugins"
EOF
```

* Start a Vault server in development mode with the configuration:

```sh
vault server -dev -dev-root-token-id="root" -config=/tmp/vault.hcl &
```

* Leave this running and open a new tab or terminal window. Authenticate to Vault:

```sh
export VAULT_ADDR='http://127.0.0.1:8200'
vault login root
```


## License

This code is licensed under the Mozilla Public License.
