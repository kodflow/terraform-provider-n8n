# Git Hooks Setup

## Architecture

Git hooks are stored in **`.devcontainer/hooks/`** and embedded in the Docker image with proper executable permissions.

During Docker image build, hooks are copied to `$HOME/.git-hooks/` with `755` permissions via:
```dockerfile
COPY --chown=vscode:vscode --chmod=755 hooks/ /home/vscode/.git-hooks/
```

During devcontainer setup, Git is configured to use this directory via:
```bash
git config core.hooksPath $HOME/.git-hooks
```

### Benefits of This Approach

✅ **Version controlled** - Hooks source is in `.devcontainer/hooks/` tracked in git
✅ **Single source of truth** - Edit hooks in `.devcontainer/hooks/`, rebuild to apply
✅ **Team synchronization** - Everyone gets the same hooks with correct permissions
✅ **No permission issues** - Permissions are baked into the Docker image
✅ **GUI compatible** - Works perfectly with GitKraken, SourceTree, etc.
✅ **Cross-platform** - No chmod needed, works on macOS/Windows/Linux

## Automatic Setup

Git hooks are automatically configured when rebuilding the devcontainer via `.devcontainer/post-create.sh`.

## For GitKraken and GUI Clients Users

Git hooks are embedded in the Docker image with correct permissions and work automatically with GitKraken and other GUI clients. No manual permission fixes needed!

The hooks are configured to run from `$HOME/.git-hooks/` which is part of the Docker image, ensuring consistent behavior across all environments.

## Available Hooks

All hooks are located in `$HOME/.git-hooks/` (source: `.devcontainer/hooks/`):

1. **pre-commit** - Auto-generates CHANGELOG.md and runs coverage
2. **prepare-commit-msg** - Suggests conventional commit format
3. **commit-msg** - Validates commit message with commitlint
4. **pre-push** - Blocks push if AI mentions or Co-Authored-By detected

## Manual Configuration

If needed, you can manually reconfigure hooks:

```bash
# Inside the devcontainer
./scripts/install-hooks.sh
```

Or directly:
```bash
git config core.hooksPath $HOME/.git-hooks
```

**Note:** No need to chmod - permissions are already set in the Docker image.

## Bypassing Hooks

**Not recommended**, but if you need to bypass hooks temporarily:

```bash
git commit --no-verify   # Skip pre-commit and commit-msg hooks
git push --no-verify     # Skip pre-push hook
```

## Troubleshooting

### Hooks not running

1. Check git config:
```bash
git config core.hooksPath
# Should output: /home/vscode/.git-hooks
```

2. Verify hooks exist:
```bash
ls -la $HOME/.git-hooks/
# All hooks should show -rwxr-xr-x permissions
```

3. Reconfigure if needed:
```bash
./scripts/install-hooks.sh
```

### commitlint not found

Rebuild the devcontainer - commitlint is installed automatically in post-create.sh.

### Modifying Hooks

To modify hooks:
1. Edit files in `.devcontainer/hooks/`
2. Rebuild the devcontainer to apply changes
3. The new hooks will be embedded in the image with correct permissions
