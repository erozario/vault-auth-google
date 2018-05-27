# Vault Auth Google.

A Vault plugin for authenticating and receiving policies via Google Accounts

### Compatibility Matrix

| plugin version | vault 0.9.x | vault 0.10.x |
|----------------|-------------|--------------|
| 0.1.0          | N           | Y            |


## Download and Install

#### Binary Distributions

* Download plugin from the [releases](https://github.com/erozario/vault-auth-google/releases) page.

* Move the compiled plugin into Vault's configured plugin_directory:

```shell
mv vault-auth-google /etc/vault/plugins/vault-auth-google
```

* Calculate and register the SHA256 sum of the plugin in Vault's plugin catalog:

```shell
SHASUM=$(shasum -a 256 "<vault-plugin-path>/vault-auth-google" | cut -d " " -f1)
vault write sys/plugins/catalog/vault-auth-google \
  sha_256="$SHASUM" \
  command="vault-auth-google"
```

* Mount the auth method:

```shell
vault auth enable \
    -path="google" \
    -plugin-name="vault-auth-google" plugin
```

#### Install From Source

* Clone this repo and execute make all

```shell
make all
```

* Move the compiled plugin into Vault's configured plugin_directory:

```shell
mv vault-auth-google /etc/vault/plugins/vault-auth-google
```

* Calculate and register the SHA256 sum of the plugin in Vault's plugin catalog:

```shell
SHASUM=$(shasum -a 256 "<vault-plugin-path>/vault-auth-google" | cut -d " " -f1)
vault write sys/plugins/catalog/vault-auth-google \
  sha_256="$SHASUM" \
  command="vault-auth-google"
```

* Mount the auth method:

```shell
vault auth enable \
    -path="google" \
    -plugin-name="vault-auth-google" plugin
```

## Configure plugin

* Create an OAuth client ID in [the Google Cloud Console](https://console.cloud.google.com/apis/credentials), of type "Other".


* Configure the auth method:

> Gmail Example:

```shell
vault write auth/google/config \
    client_id=<GOOGLE_CLIENT_ID> \
    client_secret=<GOOGLE_CLIENT_SECRET>
```

* Create a role for a given User Email, mapping to a set of policies:

```shell
vault write auth/google/role/default \
    bound_emails=eduardoagrj@gmail.com \
    policies=default
```

> Google Organizations with Google Groups Example:

You need to enable the [Admin SDK](https://developers.google.com/admin-sdk/directory/v1/guides/authorizing) to be able to list the email groups

```shell
vault write auth/google/config \
    client_id=<GOOGLE_CLIENT_ID> \
    client_secret=<GOOGLE_CLIENT_SECRET> \
    fetch_groups=true
```

* Create a role for a Google group mapping to a set of policies:

```shell
vault write auth/google/role/default \
    bound_domain=<DOMAIN> \
    bound_groups=sec@<DOMAIN>,infra@<DOMAIN> \
    policies=default
```

* Login using Google credentials

```shell
firefox $(vault read -field=url auth/google/code_url)
vault write auth/google/login code=$GOOGLE_CODE role=default
```

## Setup in Local Envinroment

* Clone this repo

```shell
git clone git@github.com:erozario/vault-auth-google.git
```

* Create a temporary directory to compile the plugin into and to use as the plugin directory for Vault:

```shell
mkdir -p /tmp/vault-plugins
```

* Compile the plugin into the temporary directory:

```shell
go build -o /tmp/vault-plugins/vault-auth-google
```

* Create a configuration file to point Vault at this plugin directory:

```shell
tee /tmp/vault.hcl <<EOF
plugin_directory = "/tmp/vault-plugins"
EOF
```

* Start a Vault server in development mode with the configuration:

```shell
vault server -dev -dev-root-token-id="root" -config=/tmp/vault.hcl &
```

* Leave this running and open a new tab or terminal window. Authenticate to Vault:

```shell
export VAULT_ADDR='http://127.0.0.1:8200'
vault login root
```

* Calculate and register the SHA256 sum of the plugin in Vault's plugin catalog:

```shell
SHASUM=$(shasum -a 256 "/tmp/vault-plugins/vault-auth-google" | cut -d " " -f1)
vault write sys/plugins/catalog/vault-auth-google \
  sha_256="$SHASUM" \
  command="vault-auth-google"
```

* Mount the auth method:

```shell
vault auth enable \
    -path="google" \
    -plugin-name="vault-auth-google" plugin
```
## License

This code is licensed under the MPLv2 license.