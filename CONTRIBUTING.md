# Contributing to Mimsy

Thank you for your interest in contributing to Mimsy! This document provides guidelines for contributing to the project.

## Development Workflow

### Quick Start
1. All work must be assigned to an issue
2. Allow 24 hours for discussion before starting work
3. Create a branch following naming conventions
4. Make your changes and push commits
5. Open a Pull Request with proper requirements
6. Get approval and merge

### Issue Management
- All tasks must be assigned to an issue (bugs, features, chores, etc.)
- Use issue templates and labels appropriately
- For new features, consider adding user stories to the issue

### Branch Management
- Create new branches off the latest `main`
- **Branch naming convention**: `<conventional-commit-prefix>/<issue-number>/<description>`
  - Examples: `feat/123/add-user-authentication`, `fix/456/resolve-login-bug`
  - Use [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) prefixes
- No direct pushes to `main` (except for automations)
- Branches are automatically deleted after merge

### Pull Requests
- Keep PRs small and scoped to a single feature/fix for easier review
- PRs must meet these requirements to merge:
  - At least 1 approving reviewer (not the creator)
  - All checks must pass (linting, tests, builds)
- Only squash merges are supported to maintain clean history
- Reviewers can merge PRs directly - no need to wait for the original author
- Aim for low review + merge cycle time
- Merge conflicts should be resolved using rebase only

### Commit Messages
- Use [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) format
- Pre-commit hooks may be added to enforce this convention

### Documentation
- Highly encouraged to create/update documentation as you develop
- All documentation should be placed in the `docs/` folder
- For new features requiring user stories, add them to `docs/user_stories.md`

### User Stories
- Required for significant new features (not minor changes like label updates)
- Format in `docs/user_stories.md`:
  ```
  ---

  [User story description]

  #<issue-number>
  ```

### Releases and Versioning
- Releases are created ad-hoc by the team
- Tags are created for dedicated versions
- Backports are only allowed for:
  - Security bugs
  - Release-breaking situations
  - Data loss edge cases
- No new development off tags - only backports

### Emergency Procedures
- Emergency bypass of merge restrictions is available but highly discouraged
- Should only be used for critical production issues

## Getting Help

For detailed workflow information, see [`docs/workflow.md`](docs/workflow.md).

If you have questions about the contribution process, please open an issue or reach out to the team.
