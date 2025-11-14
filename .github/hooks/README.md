# Git Hooks

This directory contains Git hooks that enforce project standards and automate workflows.

## Installation

Hooks are automatically configured during devcontainer setup. The devcontainer configures Git to use this directory:

```bash
git config core.hooksPath .github/hooks
```

## Available Hooks

### pre-commit

Auto-generates documentation files before each commit:

- Adds copyright headers to source files
- Generates `CHANGELOG.md` from git history
- Generates `COVERAGE.MD` from test coverage
- Automatically unstages `sdk/n8nsdk/api/openapi.yaml` (auto-generated file)

### prepare-commit-msg

Provides a conventional commit message template in your editor.

### commit-msg

**Strictly validates** commit messages using commitlint.

Required format: `<type>: <description>`

Valid types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Test additions/changes
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `build`: Build system changes
- `ci`: CI/CD changes
- `chore`: Maintenance tasks
- `revert`: Revert previous commit

Also blocks:
- AI attribution (Co-Authored-By: Claude, GPT, etc.)
- Generated with messages

### pre-push

Prevents pushing commits with AI mentions or Co-Authored-By tags to maintain clean git history.

## Manual Configuration

If not using the devcontainer:

```bash
git config core.hooksPath .github/hooks
chmod +x .github/hooks/*
```

## Bypassing Hooks

Only in exceptional cases (not recommended):

```bash
git commit --no-verify  # Skip pre-commit and commit-msg
git push --no-verify    # Skip pre-push
```

## Advantages

- ✅ Version controlled in repository
- ✅ Shared across team automatically
- ✅ No manual installation required
- ✅ Works with GUI clients (GitKraken, SourceTree, etc.)
- ✅ Cross-platform compatible
- ✅ Easy to modify and review

## Modifying Hooks

To modify hooks:

1. Edit files in `.github/hooks/`
2. Commit changes
3. Hooks are immediately active for all team members

No rebuild or manual installation required!
