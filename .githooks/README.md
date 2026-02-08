# Git Hooks

This repository includes optional Git hooks under `.githooks`.

To enable them locally:

```bash
git config core.hooksPath .githooks
```

Current hooks:

- `pre-push`: Enforces Conventional Commits (Angular) on commits being pushed.
