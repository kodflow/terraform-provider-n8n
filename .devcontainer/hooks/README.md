# Git Hooks

This directory contains the source files for Git hooks that are embedded in the Docker image.

## Architecture

During Docker build, these hooks are copied to `$HOME/.git-hooks/` with proper executable permissions (755):

```dockerfile
COPY --chown=vscode:vscode --chmod=755 hooks/ /home/vscode/.git-hooks/
```

Git is then configured to use this directory during devcontainer post-create setup.

## Setup

Hooks are automatically configured during devcontainer setup. To manually reconfigure:

```bash
./scripts/install-hooks.sh
```

Or directly:

```bash
git config core.hooksPath $HOME/.git-hooks
```

## Available Hooks

### pre-commit

Auto-generates documentation files before each commit:

- `CHANGELOG.md` - Generated from git history
- Test coverage report

### prepare-commit-msg

Suggests conventional commit message format in your editor.

### commit-msg

Validates commit messages using commitlint.

Required format: `<type>: <description>`

Valid types: `feat`, `fix`, `docs`, `test`, `refactor`, `perf`, `build`, `ci`, `chore`, `revert`

### pre-push

Prevents pushing commits with AI mentions or Co-Authored-By tags to maintain clean git history.

## Bypassing Hooks

If needed (not recommended):

```bash
git commit --no-verify
git push --no-verify
```

## Advantages Over Traditional Hooks

✅ Version controlled in `.devcontainer/hooks/` ✅ Permissions baked into Docker image - no chmod issues ✅ Easy to edit (just rebuild to apply changes) ✅
Works perfectly with GUI clients (GitKraken, SourceTree, etc.) ✅ Cross-platform compatible (macOS/Windows/Linux) ✅ Automatically shared with team via Docker
image ✅ Single source of truth

## Modifying Hooks

To modify hooks:

1. Edit files in this directory (`.devcontainer/hooks/`)
2. Rebuild the devcontainer
3. Hooks will be embedded in the new image with correct permissions

## Requirements

- commitlint: `npm install -g @commitlint/cli @commitlint/config-conventional`
- Go tools for testing (auto-installed in devcontainer)
