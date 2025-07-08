# Architecture description

This document describes the high level architecture components used to build the mimsy CMS.

## Summary

![Architecture](./imgs/architecture.png)

Mimsy is a headless CMS that can be used through its SDK. Mimsy uses the resource definition file managed by the user to define the resources that are present in the CMS.

## Technologies

### Frontend

**Sveltekit**

The frontend is written in Sveltekit. The frontend allows users to create resources in collections using a web based UI.

### Backend

**Go**

Go was chosen as the backend language for its simplicity in building networked applications. The Go backend is responsible for serving the frontend, handling authentication, and providing an API used by the admin UI.

### Database

**PostgreSQL**

PostgreSQL enables Mimsy to perform zero downtime migrations relatively easily using [pgroll](https://pgroll.com/). Mimsy can use the commit hash to determine the schema version, which the user can use to perform queries.

### SDK

**Typescript**

The SDK allows the user to define their collections and blocks. These resources are then used to create, read, update, and delete resources in the CMS. The SDK is designed to be used in a typesafe manner, leveraging Typescript's type system to ensure that the resources are correctly defined and used.

### Resource definition

The users of Mimsy must define their resources in a resource definition file. This file is used to define the collections and blocks that are present in the CMS. The resource definition file is used to infer the types of the resources which are used in the SDK.
