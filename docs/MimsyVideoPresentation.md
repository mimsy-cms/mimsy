# Mimsy video presentation

## Introduction

Reasons the project is/was "necessary"
	- Quick and easy setup phase
	- Seamless deployments with no downtime
	- Single source of truth to simplify modifications


## Technology/workflow overview

### Resource definition file (JSON)

Allows developers to define the type of page they want to use. They have endless possibilities and configurations thanks to a variety of input fields to choose from.

### GitHub CI/CD

Once updated, the definitions are then pushed to GitHub.

### SDK (TypeScript)

The homemade SDK periodically fetches the updated files from GitHub.
It then transforms the definition file into a specific database migration to go from the current state to the desired one.

### API / Backend (Go)

Acts as the middleman between the various system entities such as the SDK applying migrations or the frontend which shows a user friendly representation of the created collections and their resources.

### Admin panel (SvelteKit / TypeScript)

The admin panel is used to manage users, upload media and populate the predefined fields in the available collections.
Dynamically updated to show the recently updated definitions.

### Database / Swift storage (PostgreSQL, pgroll)

The use of database migration tools allows migrations to happen without impacting resource accessibility.

The user uploaded media is stored in an external Swift storage backend.


## Conclusion

If you are interested in an easy to implement and use, modern CMS solution Mimsy is definitely for you.
