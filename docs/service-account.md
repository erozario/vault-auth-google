# Creating a Service Account on GCP for G Suite user impersonation

## Why?

`vault-auth-google` requires Google Groups read-only capabilities in order to
list all existing groups within the organizational unit and their members.
Though this permission could be granted through OAuth, not every member inside
G Suite has the correct permissions to list groups information. The service
account is able to provide the applicaton a way to impersonate a given user and
use their permissions.

## How?

### Administrative roles on G Suite

1. Open the [G Suite Admin console](https://admin.google.com);
1. Access the **User** section, find the user and click on it;
1. Expand the **Admin roles and privileges** section;
1. Assign the **Groups Admin** role to the user.

### Service account

1. Access the [Google Cloud Plataform](https://console.cloud.google.com);
1. On the left side panel, access the **IAM & admin** section;
1. At the section, add the G Suite group admin as a *Service Account Actor*;
1. On the left side panel, access the **Service accounts** section;
1. Click on **CREATE SERVICE ACCOUNT** at the top of the page;
1. Fill the *Service account name* and *Enable G Suite Domain-wide Delegation*;
1. On the left side panel, access the **APIs & Services** section;
1. On **Credentials**, create a new  *Service account key*;
1. At the *Service account key* creation page, select the service account on
   the dropdown menu as set the *Key type* as JSON.
1. On the **Library**, enable the *Admin SDK API*;
1. Save the Service Account credendential file and the service account client
   ID;

### API scope

1. Back to the [G Suite Admin console](https://admin.google.com), access the
   **Security** section;
1. At the section, expand the **Advanced settings** and click on *Manage API
   client access*;
1. On the **Client Name** text box, insert the *Client ID* of the newly created
   service account;
1. On the **One or More API Scopes**, insert the *Google Groups Read-Only API
   scope*.
1. Authorize the API client access clicking on **Authorize**.


## References

* **Google Groups Read-Only API scope**: `https://www.googleapis.com/auth/admin.directory.group.readonly`
* [About admnistrator role](https://support.google.com/a/answer/33325?hl=en)
* [Assign administrator to a user](https://support.google.com/a/answer/172176?hl=en)
