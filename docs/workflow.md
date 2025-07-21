# Development Workflow

This document provides detailed information about our development workflow and processes.

## Overview

Our development workflow is designed to maintain high code quality while enabling fast iteration and collaboration. The process emphasizes small, focused changes with quick review cycles.

## Trunk-Based Development

We follow a trunk-based development model, where all developers work on short-lived branches that are frequently merged back to the main branch (`main`). This approach enables continuous integration and faster delivery cycles.

### Key Decisions and Benefits

Our choice of trunk-based development drives several architectural and process decisions:

#### Continuous Integration and Merge Queues
- **CI Validation**: We use continuous integration to validate that changes don't break anything before they reach `main`
- **Merge Queue Ready**: Our setup is ready to enable merge queues to ensure that `main` can always be built and deployed
- **Quality Gates**: All automated checks (linting, tests, builds) must pass before merging

#### Fast Merge Culture
- **24-Hour Target**: We target less than 24-hour turnaround time from PR creation to merge
- **Merge Time Dashboard**: We have a system developed for another context that provides a dashboard of merge times to track our performance
- **Quick Reviews**: Internal culture encourages fast reviews and merging to maintain velocity

#### GitOps with FluxCD
- **GitOps Approach**: We chose FluxCD for their excellent GitOps approach and git-based updates
- **Infrastructure as Code**: All deployments and configurations are managed through git, ensuring consistency and traceability
- **Automated Deployments**: Changes to `main` trigger automated deployments through FluxCD

#### Monorepo Strategy
- **Single Repository**: Everything is contained within a monorepo to ensure cohesive testing and versioning
- **Component Integration**: Every part of the framework is tested with all components using the same version
- **Simplified Dependencies**: Eliminates version mismatches and integration issues between components
- **Atomic Changes**: Related changes across multiple components can be made in a single commit/PR

## Issue Management

### Issue Requirements
- **All work must be tracked**: Every task, whether it's a bug fix, new feature, or maintenance work, must be assigned to an issue
- **Issue types and labels**: Use appropriate issue types (bug, feature, enhancement, chore) and labels for categorization
- **Templates**: Follow issue templates when creating new issues to ensure consistent information capture

### Discussion Period
- **Purpose**: This gives team members time to:
  - Provide input on the approach
  - Suggest alternatives or improvements
  - Identify potential conflicts with other work
  - Ask clarifying questions

## Branch Management

### Branch Creation
- **Base branch**: Always create new branches from the latest `main` branch
- **Naming convention**: `<conventional-commit-prefix>/<issue-number>/<description>`

  Examples:
  - `feat/123/add-user-authentication`
  - `fix/456/resolve-login-bug`
  - `docs/789/update-api-documentation`
  - `refactor/101/simplify-user-service`

### Conventional Commit Prefixes
Use the standard [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) prefixes:
- `feat`: New features
- `fix`: Bug fixes
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring without functional changes
- `test`: Adding or updating tests
- `chore`: Maintenance tasks, build changes, etc.

### Branch Protection
- **No direct pushes**: Direct pushes to `main` are prohibited (except for automated processes)
- **Automatic cleanup**: Branches are automatically deleted after successful merge

## Development Process

### Making Changes
- **Multiple commits allowed**: Push changes in as many commits as desired during development
- **Commit message format**: Use Conventional Commits format for all commit messages
- **Pre-commit hooks**: Automated hooks may enforce commit message conventions and run linting/formatting

### Pull Request Guidelines

#### PR Scope and Size
- **Single responsibility**: Each PR should address only one feature, bug, or improvement
- **Keep PRs small**: Smaller PRs are easier and faster to review
- **Encourage frequent merging**: Instead of stacking PRs, merge frequently to reduce complexity

#### PR Requirements (Enforced by GitHub)
1. **Reviewer approval**: At least one approving review from someone other than the PR creator
2. **Passing checks**: All automated checks must pass:
   - Linting
   - Tests
   - Builds

#### Review Process
- **Fast cycle time**: Aim for quick review and merge cycles
- **Reviewer can merge**: Reviewers are empowered to merge PRs directly without waiting for the original author
- **Merge method**: Only squash merges are supported to maintain a clean commit history

### Handling Merge Conflicts
- **Rebase only**: Use rebase (not merge commits) to resolve conflicts
- **Keep history linear**: This maintains a cleaner, more readable git history

## Documentation Standards

### Documentation Requirements
- **Document as you go**: Highly encouraged to create or update documentation during development
- **Location**: All documentation must be placed in the `docs/` folder
- **Types of documentation**: API docs, architecture decisions, setup guides, user guides, etc.

### User Stories
User stories are required for significant new features (not for minor changes like label updates).

#### Format
Add user stories to `docs/user_stories.md` using this format:

```
---

As a [type of user], I want [some goal] so that [some reason].

#<issue-number>
```

#### Example
```
---

As a project manager, I want to view a dashboard of all open issues so that I can track project progress at a glance.

#123
```

## Release Management

### Versioning and Tags
- **Ad-hoc releases**: Releases are created as needed by the team, not on a fixed schedule
- **Tagging**: Create tags for dedicated versions
- **Semantic versioning**: Follow semantic versioning principles (MAJOR.MINOR.PATCH)

### Backport Policy
Backports are only allowed for critical issues:
- **Security vulnerabilities**: Security bugs that affect previous versions
- **Release-breaking situations**: Issues that prevent normal operation
- **Data loss scenarios**: Bugs that could cause data corruption or loss

#### Backport Process
1. Fork a branch from the relevant tag
2. Apply only the minimal fix (no new features)
3. Create a new patch release
4. No new development is allowed off tags - only backports

## Emergency Procedures

### Emergency Bypass
- **Availability**: Emergency bypass of merge restrictions is technically available
- **Usage**: Highly discouraged and should only be used for critical production issues
- **Documentation**: Any emergency bypass should be documented with justification

### When to Use Emergency Bypass
- Production is down and a hotfix is needed immediately
- Security vulnerability requires immediate patching
- Data loss is occurring and requires immediate intervention

## Quality Assurance

### Automated Checks
All PRs must pass these automated checks:
- **Linting**: Code style and quality checks
- **Tests**: Unit tests, integration tests, and any other test suites
- **Builds**: Successful compilation/build process

### Review Standards
- Focus on code quality, maintainability, and correctness
- Check for proper error handling and edge cases
- Ensure documentation is updated when necessary
- Verify that the change addresses the issue requirements

## Best Practices

### Communication
- Use clear, descriptive commit messages
- Write informative PR descriptions explaining the change and its impact
- Link PRs to their corresponding issues
- Tag relevant team members for review when appropriate

### Code Quality
- Write self-documenting code with clear variable and function names
- Add comments for complex logic or business rules
- Follow established coding standards and patterns
- Write tests for new functionality

### Collaboration
- Be responsive to review feedback
- Provide constructive feedback during code reviews
- Ask questions when requirements are unclear
- Share knowledge through documentation and code comments

## Troubleshooting

### Common Issues
- **Branch naming**: Ensure branch names follow the exact convention
- **Merge conflicts**: Always use rebase to resolve conflicts
- **Failed checks**: Address all linting, test, and build failures before requesting review
- **Missing reviews**: Ensure at least one team member has approved the PR

### Getting Help
- Check existing documentation in the `docs/` folder
- Ask questions in issues or PR comments
- Reach out to team members for guidance on complex changes
