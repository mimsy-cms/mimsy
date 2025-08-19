# Authentication

The authentication system in Mimsy uses session tokens to manage user sessions. These tokens are stored in the database. The session token can be passed in the `Authorization` header of API requests as a Bearer token, or as a cookie named `session`.

The API requires authentication for `POST`, `PUT`, and `DELETE` endpoints. For `GET` requests, no authentication is required. This behavior might change in the future if/when private collections are implemented.
