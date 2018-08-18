# Vault Auth Google [![Build Status](https://travis-ci.com/erozario/vault-auth-google.svg?branch=master)](https://travis-ci.com/erozario/vault-auth-google) [![Go Report Card](https://goreportcard.com/badge/github.com/erozario/vault-auth-google)](https://goreportcard.com/report/github.com/erozario/vault-auth-google)

A [Hashicorp Vault](https://github.com/hashicorp/vault) plugin that enables
authentication and policy bounding via Google Accounts and Google Groups (G
Suite only).


## Table of Contents

 - [Compatibility Matrix](#compatibility-matrix)
 - [Getting](#getting)
    - [Binary](#binary)
    - [From source](#from-source)
 - [Google API requirements](#google-api-requirements)
 - [Installation](#installation)
 - [Configuration & Usage](#configuration--usage)
    - [Parameters](#binary)
    - [Local flow vs. Web-based flow](#local-flow-vs-web-based-flow)
    - [Bare-minimum settings](#bare-minimum-settings)
    - [How to...](#how-to)
 - [Contributing](#contributing)
 - [License](#license)


## Compatibility Matrix

| Plugin Version | Vault 0.9.x  | Vault 0.10.x |
|----------------|--------------|--------------|
| 1.0.0          | :thumbsdown: | :thumbsup:   |
| 0.1.0          | :thumbsdown: | :thumbsup:   |


## Getting

You can get the plugin binary via the release page or building yourself. After
getting the binary, move it into the [plugin
directory](https://www.vaultproject.io/guides/operations/plugin-backends.html)
set at the Vault's configuration file.

### Binary

The pre-compiled binary is available at the [release](https://github.com/erozario/vault-auth-google/releases) page.


### From source

Clone this repository and build via `make`:

```sh
make all
```

The make recipe requires [dep](https://github.com/golang/dep) in order to get
the project's dependencies.

Alternatively, you can get the dependencies via `go get`.


## Google API requirements

`vault-auth-google` requires an OAuth2 credential in order to authenticate
Google Accounts into Vault. For G Suite users with group bounding, the plugin
also requires a domain-wide service account for user impersonation.

Though it requires the `Admin SDK` API enabled, `vault-auth-google` only make
use of the [admin.directory.group.readonly](https://developers.google.com/admin-sdk/directory/v1/guides/authorizing)
function.

For more information on how to create OAuth2 credentials and service account
keys, check the docs:

* [Creating an OAuth Credential on GCP](docs/oauth.md)
* [Creating a Service Account on GCP for G Suite user impersonation](docs/service-account.md)


## Installation

Third-party auth methods, such as this, cannot be enabled the same way as
native [auth methods](https://www.vaultproject.io/docs/auth/index.html). To
install it, add the plugin into the [plugin
catalog](https://www.vaultproject.io/docs/internals/plugins.html#plugin-catalog).

* After [getting the plugin](#getting) binary, calculate the SHA256 of it and stores it
  inside an environment variable for easy reference.

```sh
SHASUM=$(shasum -a 256 "/path/to/vault-auth-google" | cut -d " " -f1)
```

* Register the binary and its SHASUM in Vault's plugin catalog

```sh
vault write sys/plugins/catalog/vault-auth-google \
  sha_256="$SHASUM" \
  command="vault-auth-google"
```

* Finally, mount the `vault-auth-google` as an auth method.

```sh
vault auth enable \
    -path="google" \
    -plugin-name="vault-auth-google" plugin
```


## Configuration & Usage

The configuration of this auth method depends on the kind of the Google Account
used and the OAuth2 flow. Both Gmail and G Suite accounts requires OAuth2
credentials for user authentication, but the [kind of flow can
differ](https://developers.google.com/identity/protocols/OAuth2).

The plugin can be configured via parameters. Each parameter have a different
effect on the plugin and serves different purposes.


### Parameters

`vault-auth-google` can be configured via five parameters, two of which are
required:

 - _(string)_ `client_id` __*__: The Google OAuth2 Client ID.
 - _(string)_ `client_secret` __*__: The Google OAuth2 Client secret.
 - _(boolean)_ `fetch_groups`: Should the plugin bound policies to groups? **true** if yes, **false** otherwise.
 - _(string)_ `redirect_url`: The URL that Google will redirect after the
     OAuth2 flow. This URL should also be added at the credentials authorized URIs.
 - _(string_ `delegation_user`: The Google user that delegates the API permission.
 - _(string)_ `service_acc_key`: The content of the Service Account private key.

__* Required parameters__


### Local flow vs. Web-based flow

The flow can be made on the [local
machine](https://developers.google.com/identity/protocols/OAuth2InstalledApp)
or via [web
services](https://developers.google.com/identity/protocols/OAuth2WebServer).

Local authentication flows generates a temporary, one-time usage Google token.
The Google OAuth token have to be written in Vault, so it can then generate
it's own access token with policy builtin on.

The difference between a local flow and a web-based flow on this plugin mostly
relies on the `redirect_url` parameter. This parameter is optional and unset by
default. When unset, the plugin assumes a local flow. Since Vault cannot handle
a GET callback from Google, the token has to be feeded manually. E.g.:

* Gets a Google OAuth flow URL from Vault and opens it with Firefox.

```sh
firefox $(vault read -field=url auth/google/code_url)
```

* Writes the Google code on Vault for a token generation.

```sh
vault write auth/google/login code=<GOOGLE-OAUTH2-CODE> role=default
```

Alternatively, when `redirect_url` is setted, the plugin assumes a web-based
flow and uses the given URL as  the Redirect URI used by the Google OAuth2
credential. It is important to notice that this URL has to be a web application
ready to handle GET requests, since Google will make the request on it, passing
the authorization code as a GET parameter.


### Bare-minimum settings

Below lies a simple, bare-minimum configuration that enables
`vault-auth-google` to work properly alongside your Vault server.

```hcl
listener "tcp" {
        address = "0.0.0.0:8200"
}

storage "file" {
        path = "/path/to/data"
}

plugin_directory = "/path/to/vault/plugin/dir"
api_addr = "https://vault.domain.com:8200"
```

Beyond the required parameters, such as `listener` and `storage`, attention
should be given to `plugin_directory` and `api_addr`.

* The [`plugin_directory`](https://www.vaultproject.io/docs/configuration/index.html#plugin_directory)
  parameter tells Vault in which directory rests all plugins used by it. Have
  you downloaded or compiled the plugin, the binary should be located at the
  absolute path defined in this parameter.

* The [`api_addr`](https://www.vaultproject.io/docs/configuration/index.html#api_addr)
  parameter specifies the address of the Vault server listener, for client
  redirection (i.e., when the main Vault Server receives a request that points
  to another server, in which listener should the client be redirect to?).
  Since `vault-auth-google` is, in practice, another Vault server, this
  parameter should be setted; otherwise, the plugin will not be able to respond
  and will fail in every request made upon him. In general this should be set
  as a full URL that points to the value of the listener address.


### How to...

* [Configure `vault-auth-google` and use it for Gmail accounts](docs/gmail.md)
* [Configure `vault-auth-google` and use it for G Suite accounts (with group bounding)](docs/gsuite.md)


## Contributing

If you wish to contribute to this project, either by fixing a bug or suggesting
new functionalities, check the documentation on how to setup a [local
environment for development](docs/local-dev.md).


## License

This project is licensed under the [Mozilla Public License](LICENSE).
