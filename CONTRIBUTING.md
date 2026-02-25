# Contributing to poof

We love your input! We want to make contributing to `poof` as easy and transparent as possible, whether it's:

- Reporting a bug
- Discussing the current state of the code
- Submitting a fix
- Proposing new features
- Becoming a maintainer

## Development Process

We use [Task](https://taskfile.dev/) to manage our development workflow.

1. **Clone the repository**: `git clone https://github.com/chege/poof.git`
2. **Install dependencies**: `go mod download`
3. **Run tests**: `go test ./...`
4. **Make changes**: Hack away!
5. **Verify changes**: `task ready` (runs fmt, lint, and tests).
6. **Submit a PR**: Open a pull request against the `main` branch.

## Standards

- Follow [Conventional Commits](https://www.conventionalcommits.org/) for commit messages.
- If you want local enforcement, enable the pre-push hook: `git config core.hooksPath .githooks`.
- Ensure all tests pass.
- Add tests for new features or bug fixes.
- Keep the code clean and well-documented.

## Code of Conduct

Please be respectful and professional in all interactions within this project.

---
*Thank you for helping make `poof` better!*
