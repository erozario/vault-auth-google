# Setting up a local development environment

You'll first need Go properly installed on your machine. You can follow Go's
documentation on [how to get started](https://golang.org/doc).

1. Clone this repository inside Go's workspace. The workspace path is generally
   mapped into the `GOPATH` environment variable.

```sh
git clone git@github.com:erozario/vault-auth-google.git $GOPATH/src
```


1. Create a temporary directory to compile the plugin into and to use as the
   plugin directory for Vault.

```sh
mkdir -p /tmp/vault-plugins
```


1. Compile the plugin into the temporary directory.

```sh
cd $GOPATH/src/vault-auth-google && go build -o /tmp/vault-plugins/vault-auth-google
```


1. Create a configuration file that sets the temporary directory as the Vault's
   plugin directory.

```sh
tee /tmp/vault.hcl <<EOF
plugin_directory = "/tmp/vault-plugins"
EOF
```


1. Start the Vault server in development mode, referencing the configuration
   file.

```sh
vault server -dev -dev-root-token-id="root" -config=/tmp/vault.hcl &
```


1. Leave this running and open a new tab or terminal window. Authenticate to Vault:

```sh
export VAULT_ADDR='http://127.0.0.1:8200'
vault login root
```
