# PDG - Final Delivery

## Overview of the issue and the solution
Headless CMSes are becoming increasingly popular, and many developers are looking to integrate them into their projects. Unfortunately, the current status quo of existing projects is not ideal for the developer experience. Current solutions like Payload and WordPress support many different frontend technologies instead of focusing on a single stack. This is why Mimsy was started: to experiment and understand whether making an opinionated CMS is the way to go.

## Source code

The source code is available on GitHub at [https://github.com/mimsy-cms/mimsy](https://github.com/mimsy-cms/mimsy).

A demonstration of how the CMS is used is available at [https://github.com/mimsy-cms/slides-demo](https://github.com/mimsy-cms/slides-demo).

## How to launch the project locally

#### Install nix
To encourage reproducibility, the project ships with a prepared nix configuration file (flake.nix).

First, you need to have [nix](https://nixos.org/download/#download-nix) installed. We recommend using [Determinate Nix](https://docs.determinate.systems/determinate-nix/) to install nix more easily.

#### Environment file
To support the development environment, a `.env` file is required. You can create one by copying the `.env.example` file inside the `api` directory and filling in the required values.

You are going to need a GitHub Application project to enable the project syncing. You can create one by going to [https://github.com/settings/apps/new](https://github.com/settings/apps/new) and filling in the required values.

To sync, the application needs to be installed on the GitHub organization or user account that owns the repository, with the permission to read and write to the repository.

You can then export your PEM key for the project inside of the settings, and provide it inside of the environment file, along with the application ID.

#### Start the project

Once you have a shell with nix installed, you can launch the project locally by running the following command:

```bash
nix develop # Opens the shell with all required dependencies
dev # Starts the project
```

The `dev` command starts all the required dependencies (Postgres) along with all components of the deployment.

You can then access the CMS at [http://localhost:5173](http://localhost:5173).

More information can be found in the [README](../../README.md).

## Continuous integration
The CI/CD pipeline is fully automated and runs on GitHub Actions and Flux.

Once you have made your changes, create a PR with them. The pipeline will automatically build and test your changes before allowing for a merge to happen.

Once the changes are merged, the pipeline will automatically deploy the changes to the production environment. It happens in three steps:

1. The changes are built and pushed to the `ghcr.io` registry.
2. The new image is found by the flux controller, and a commit is made to the `main` branch.
3. The changes are deployed to the production environment.

For more details, check the [devops documentation](../devops.md).

## Contribution Guidelines

As Mimsy is an open source project, we welcome contributions from the community. If you would like to contribute, please follow these guidelines:

1. **Create an issue describing the changes you want to make.** As Mimsy is opinionated, we encourage you to discuss your ideas with us before making any changes, as otherwise we might not accept your changes.
2. Fork the repository and create a new branch for your changes.
3. Make your changes and ensure that they pass all tests.
4. Submit a pull request with your changes.

We will review your pull request and provide feedback if necessary. Once your changes have been approved, we will merge them into the main branch.

If you want to get started quickly, check the issue tracker for open issues with the label "good first issue". These issues are usually well-suited for beginners.

For more details, check:
- [Contributing Guidelines](../../CONTRIBUTING.md)
- [README](../../README.md)
- [Workflow documentation](../workflow.md)

## Video

The Video is available at the following URL:
TODO
