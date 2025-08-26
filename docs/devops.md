# DevOps

This document outlines the DevOps strategy used in the Mimsy project, covering the processes used for development, testing, and deployment.

## Overview

Mimsy uses a GitOps approach to manage infrastructure and deployments. This means that all changes to the infrastructure are made through Git pull requests, and the desired state of the system is stored in Git repositories. This allows for version control, collaboration, and automated deployments.

We use [Flux](https://fluxcd.io/) to continuously deploy Mimsy on a Kubernetes cluster. You can find the files related to Flux inside the [flux directory](../flux). In addition to Flux, we are using [SOPS](https://getsops.io/) to manage secrets inside the repository.

The Mimsy application that is deployed on Kubernetes serves as a validation ground for changes, we expect users of Mimsy to deploy the CMS themselves.

## Develop

The Mimsy project is structured as a monorepository, where each directory inside the repository corresponds to a different app or library. You can find more about each project's inside their respective `README.md` files.

## Test

Each project inside the monorepository has its own set of tests, these tests run automatically when a new pull request is created and files related to a given project are modified.

## Contribute

Once you have made your changes, you can submit a pull request to the repository. Make sure to follow the [contribution guidelines](./CONTRIBUTING.md) when submitting your pull request.

Refer to the [workflow.md](./workflow.md) document for details on the development process.

## Deploy

When a pull request is merged on the main branch, the [build.yaml](../.github/workflows/build.yaml) Github workflow is triggered. This workflow tests, builds and publishes the [API](../api) and [Admin panel](../web) as docker containers to the Github container registry.

Flux, which is running on Kubernetes, continuously monitors the container registry for new tags. When a new tag is detected (uses timestamps), Flux will automatically update the flux deployment with the new Docker image, which will then be used by the Kubernetes cluster once the reconciliation process is complete.
