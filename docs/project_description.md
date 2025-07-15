# Project description - Mimsy

## Goals

The Mimsy project aims to create a modern and developer-oriented headless CMS for resource management.

### Resource Management

- **Code-First Schema Definition**: Enable developers to define collections, fields, and types entirely in configuration files.
- **Type Safety**: Provide compile-time validation through an in-house SDK to leverage TypeScripts type system.
- **Simplified Resource Creation**: Reduce the complexity and initial setup costs compared to highly customizable solutions like Payload CMS.

### Seamless SvelteKit Integration

- **First-Class Ecosystem Support**: Provide integration with Svelte/SvelteKit applications as the primary development framework.
- **SDK-Driven Development**: Offer a homemade SDK that allows efficient resource access and management within SvelteKit applications.
- **Dynamic API Generation**: Automatically generate REST API endpoints based on resource definitions without manual configuration.

### Zero-Downtime Operations

- **Live Migration Support**: Enable schema changes without service interruption using pgroll for PostgreSQL.
- **Kubernetes-Compatible Deployments**: Support deployment patterns where old and new versions run simultaneously for validation.
- **Version-Controlled Schema Evolution**: Use Git commit hashes to track and manage different schema versions.

### Operational Excellence

- **Source of Truth Clarity**: Eliminate confusion by making code the single source of truth for all schema modifications.
- **Simplified Technology Stack**: Minimize language complexity and overlap by using Go for the backend and TypeScript for the frontend/SDK.

## Functional requirements

1. The CMS must allow developers to define collections and their fields in a definition file.
2. Generate the SDK from the definition file, either using TypeScript's reflection, or code generation
3. The CMS must automatically generate REST API endpoints based on the defined collections.
4. The CMS must support generating and running schema migrations using the definition file.
5. The CMS must perform zero-downtime migrations for schema changes.
6. The admin interface must be accessible only to authenticated users.
7. The CMS must provide user account management features.
8. Users of the CMS must be able to version their definition file inside Git.
9. The admin interface must allow users to create, read, update, and delete resources defined in any collection.
10. The admin interface must display forms for each collection based on the defined fields for the collection.
11. The definition file must be the single source of truth for the schema.

## Non-functional requirements

1. The CMS backend must be stateless and must support horizontal scaling.
2. During active schema migrations, the CMS backend must support handling requests that specify a particular schema version.
3. The actions performed on resources inside collections must adhere to the ACID properties.

## User Stories

| Name                                                                                | Role                | Goal                                                                                                                                 | Benefit                                                                          |
| ----------------------------------------------------------------------------------- | ------------------- | ------------------------------------------------------------------------------------------------------------------------------------ | -------------------------------------------------------------------------------- |
| [Definition file](https://github.com/mimsy-cms/mimsy/issues/22)                     | As a developer      | I want to define collections and fields in a configuration file                                                                      | In order to version control my schema.                                           |
| [Typesafe resource management](https://github.com/mimsy-cms/mimsy/issues/23)        | As a developer      | I want to use a TypeScript SDK for type-safe resource management                                                                     | In order to avoid runtime errors.                                                |
| [Content editing](https://github.com/mimsy-cms/mimsy/issues/24)                     | As a content editor | I want to edit the content of the resources of my application                                                                        | In order to make changes to the content displayed in my application.             |
| [Automated migrations](https://github.com/mimsy-cms/mimsy/issues/25)                | As a developer      | I want Mimsy to automatically generate migrations between schema versions                                                            | In order to remove the burden of manual migration.                               |
| [Transparent schema migrations](https://github.com/mimsy-cms/mimsy/issues/26)       | As a developer      | I want the resources inside Mimsy to always be accessible even when performing schema updates                                        | In order to remove downtime during schema updates.                               |
| [Editor access management](https://github.com/mimsy-cms/mimsy/issues/27)            | As a developer      | I want to give access to the resources in Mimsy to content editors                                                                   | In order to enable them to manage content without developer assistance.          |
| [SDK resource management](https://github.com/mimsy-cms/mimsy/issues/28)             | As a developer      | I want to manage resources in collections using the SDK                                                                              | In order to have typesafety guarantees when working with resources.              |
| [Collection management UI generation](https://github.com/mimsy-cms/mimsy/issues/29) | As a developer      | I want the content editor to have access to an up-to-date resource management interface the moment the collections schema is updated | In order to remove friction and delay between the developers and content editors |
