# File storage

This document describes the file storage system used by the Mimsy CMS.

## Overview

The file storage system is designed to handle media files uploaded by users. The file storage backend is extensible, allowing for different storage providers to be used. Currently, the system only supports OpenStack Swift as a storage backend.

## Setup

To set up the file storage system, you need to configure the environment variables for the storage backend. The following variables are required for the OpenStack Swift backend:

- `STORAGE`: Set to `swift` to use the OpenStack Swift backend.
- `SWIFT_USERNAME`: The username for the OpenStack Swift account.
- `SWIFT_API_KEY`: The API key for the OpenStack Swift account.
- `SWIFT_AUTH_URL`: The authentication URL for the OpenStack Swift service.
- `SWIFT_DOMAIN`: The domain for the OpenStack Swift account.
- `SWIFT_TENANT`: The tenant name for the OpenStack Swift account.
- `SWIFT_REGION`: The region where the OpenStack Swift service is hosted.
- `SWIFT_CONTAINER`: The name of the container in OpenStack Swift where files will be stored.

### Infomaniak

In order to use Infomaniak's OpenStack Swift service, you need to setup a [Public Cloud](https://www.infomaniak.com/en/hosting/public-cloud). Once you are logged into the Horizon dashboard, you can create a new container that will be used by Mimsy. You can also download your `OpenStack RC File` which contains the necessary credentials to access your OpenStack Swift account.
