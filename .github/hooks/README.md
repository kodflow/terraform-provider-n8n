# Git Hooks

This directory contains Git hooks that enforce project standards and automate workflows.

## Installation

Hooks are automatically configured during devcontainer setup. The devcontainer configures Git to use this directory:

```bash
git config core.hooksPath .github/hooks
```

## Available Hooks

### pre-commit

**Blocks direct commits to main branch** and auto-generates documentation:

- **🔒 Prevents commits directly to `main` branch** (use feature branches + PR)
- Adds copyright headers to source files
- Generates `COVERAGE.MD` from test coverage
- Automatically unstages `sdk/n8nsdk/api/openapi.yaml` (auto-generated file)

### prepare-commit-msg

Provides a conventional commit message template in your editor.

### commit-msg

**Strictly validates** commit messages and **enforces GPG signatures**.

#### GPG Signature Enforcement

**ALL commits MUST be signed with a GPG key.** The hook verifies:

1. GPG signing is enabled (`commit.gpgsign=true`)
2. A signing key is configured (`user.signingkey`)
3. The GPG key exists in your keyring

If any check fails, the commit will be **REJECTED** with instructions to fix the issue.

**Setup GPG signing:**

```bash
# 1. List your GPG keys
gpg --list-secret-keys --keyid-format=long

# 2. Enable GPG signing
git config --global commit.gpgsign true

# 3. Set your signing key
git config --global user.signingkey YOUR_GPG_KEY_ID
```

#### Commit Message Validation

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

### post-commit

**Verifies** that the commit was actually signed with GPG after creation.

This hook checks the GPG signature status:

- ✅ **G** (Good): Valid signature from trusted key
- ❌ **B** (Bad): Invalid signature
- ⚠️ **U** (Untrusted): Valid signature but key not trusted
- ❌ **X** (Expired): Signature with expired key
- ❌ **R** (Revoked): Signature with revoked key
- ❌ **N** (No signature): Commit not signed (should not happen if `commit-msg` hook works)

The hook displays the signer name and key ID after successful signing.

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
