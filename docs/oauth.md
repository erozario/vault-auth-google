# Creating an OAuth Credential on GCP

## Why?

The OAuth credential provides Vault a way to authenticate through your G Suite
account and get basic information about you, like your email address. This is
credential is essencial in order to Vault verifies who you are.

## How?

1. Open the [Google Cloud Plataform](https://console.cloud.google.com) console;
1. On the left side panel, access the **APIs & Services** section;
1. At the section, on the left side panel, click on **Credentials**;
1. Configure the **Oauth consent screen**, configuring a *product name*;
1. On the **Library**, enable the *Admin SDK API*;
1. On **Credentials**, create a new  *OAuth client ID* credential;
1. Set the **Application type** as *Web application* and fill the *Name* field;
1. Save the **OAuth Client ID** and **Client Secret**.

## References

* [Setting up OAuth 2.0](https://support.google.com/cloud/answer/6158849?hl=en)
* [Enable and disable APIs](https://support.google.com/cloud/answer/6158841?hl=en&ref_topic=6262490)
