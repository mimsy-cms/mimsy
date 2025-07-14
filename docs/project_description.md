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

### Core CMS Functionality

#### Resource Definition Management
- The system shall allow developers to define collections using TypeScript interfaces and configurations.
- The system shall support field types including `text`, `number`, `boolean`, `date`, `relationships`, and `media`.
- The system shall validate resource definitions at compile time using TypeScript type checking.
- The system shall generate schema definitions for multiple versions from TypeScript resource definitions.

#### CLI Tool Operations
- The CLI shall discover and process TypeScript resource definitions from project files.
- The CLI shall store migration decisions in a `.mimsy` folder for future reference.
- The CLI shall generate migration scripts compatible with pgroll for PostgreSQL.
- The CLI shall validate resource definitions before processing.
- The CLI shall support both incremental and full schema regeneration.

#### Admin Panel Interface
- The admin panel shall dynamically generate forms based on resource definitions.
- The admin panel shall provide CRUD operations (Create, Read, Update, Delete) for all defined collections.
- The admin panel shall support media upload and management.
- The admin panel shall provide user authentication and authorization.
- The admin panel shall display validation errors and success messages.
- The admin panel shall NOT allow schema modifications through the UI.

#### Backend API Services
- The backend shall expose REST API endpoints dynamically based on resource definitions.
- The backend shall handle authentication for both the admin panel and the external API access.
- The backend shall serve the admin panel static files.
- The backend shall provide filtering, sorting, and pagination for collection queries.
- The backend shall support relationship queries between collections.
- The backend shall validate incoming data against resource definitions.

#### SDK Integration
- The SDK shall provide type-safe resource access for SvelteKit applications.
- The SDK shall support querying collections with filtering and sorting.
- The SDK shall handle caching and request optimization.
- The SDK shall provide TypeScript interfaces for all defined resources.
- The SDK shall support both development and production environments.

### Migration and Versioning

#### Live Migration Support
- The system shall support zero-downtime schema migrations using pgroll.
- The system shall maintain multiple schema versions simultaneously during migrations.
- The system shall use commit hashes to identify and query specific schema versions.
- The system shall provide rollback capabilities for failed migrations.
- The system shall automatically clean up old schema versions after successful migrations.

#### Version Control Integration
- The system shall track schema changes through Git commits.
- The system shall support branching and merging of schema changes.
- The system shall detect conflicts in concurrent schema modifications.



## Non-functional requirements

### Performance Requirements

#### Response Time
- The API responses shall complete within 200ms for "simple" queries under normal load.
- The admin panel page loads shall complete within 1 second.
- The schema migrations shall complete within 5 minutes for datasets up to 1 million records.
- The CLI processing shall complete within 30 seconds for projects with up to 100 collections.

#### Throughput
- The system shall handle at least 1,000 concurrent API requests.
- The database shall support at least 10,000 read operations per second.
- The admin panel shall support at least 100 concurrent users.

### Scalability Requirements

#### Horizontal Scaling
- The backend services shall be stateless to support horizontal scaling.
- The system shall support container orchestration platforms (e.g. Kubernetes).

#### Data Volume
- The system shall support databases with up to 100GB of content data.
- The system shall handle collections with up to 1 million records.
- The system shall support file storage up to 10GB.

### Reliability Requirements

#### Availability
- The system shall maintain 99% uptime during business hours.
- The schema migrations shall not cause service downtime.
- The system shall support graceful degradation during partial failures.

#### Data Integrity
- The system shall ensure ACID compliance for all database operations.
- The system shall provide automatic backups with point-in-time recovery.
- The system shall validate all data against schema definitions before persistence.

### Security Requirements

#### Authentication & Authorization
- The system shall support role-based access control.
- The admin panel shall require authentication for all operations.
- The system shall enforce password complexity requirements.
  - SPECIFICS??

#### Data Protection
- The system shall encrypt sensitive data at rest.
- The system shall use HTTPS for all client-server communications.
- The system shall sanitize all user inputs to prevent injection attacks.
- The system shall provide audit logging for all administrative actions.

### Usability Requirements

#### Developer Experience
- The resource definitions shall provide support in TypeScript-compatible editors.
- The error messages shall be clear and actionable.
- The documentation shall include code examples for common use cases.
- The CLI shall provide helpful error messages and suggestions.

#### Admin Interface
- The admin panel shall be responsive and work on mobile devices.
- The admin panel shall provide intuitive navigation for content management.
- The admin panel shall provide contextual help and tooltips.

### Maintainability Requirements

#### Code Quality
- The system shall use consistent coding standards across all components.
- The system shall provide comprehensive API documentation.
- The system shall support automated testing for all major functionalities.

#### Deployment
- The system shall support containerized deployment.
- The system shall provide database migration scripts for version upgrades.
- The system shall support environment-specific configurations.
- The system shall provide monitoring and logging capabilities.

### Compatibility Requirements

#### Platform Support
- The backend shall run on Linux, macOS, and Windows.
- The admin panel shall work on modern browsers (Chrome, Firefox, Safari, Edge).
- The SDK shall be compatible with Node.js 18+ and SvelteKit 1.0+.
- The system shall support PostgreSQL 14+ with pgroll compatibility.

#### Interoperability
- The resource definition format shall be language-agnostic and JSON-serializable.
- The API shall follow REST conventions for third-party integration.
- The system shall provide OpenAPI specifications for all endpoints.
- The system shall support standard authentication protocols (JWT, OAuth2).

