# GPG Setup for Terraform Registry

## Quick Setup (Manual)

The automated GPG generation (`make gpg/generate`) may not work in some environments due to GPG agent restrictions. Here's the manual process:

### Step 1: Generate GPG Key

Run this command:

```bash
gpg --full-gen-key
```

When prompted, select:
1. **Kind of key**: `(9) ECC (sign only)`
2. **Elliptic curve**: `(1) Curve 25519`
3. **Valid for**: `0` (key does not expire)
4. **Real name**: `Kodflow`
5. **Email**: `133899878+kodflow@users.noreply.github.com`
6. **Comment**: (press ENTER to skip)
7. **Passphrase**: Press ENTER twice (no passphrase - required for GitHub Actions)

### Step 2: Configure Git

```bash
make gpg/configure
```

This will automatically configure git to use your new GPG key for signing commits and tags.

### Step 3: Export Keys for GitHub

```bash
make gpg/export
```

This will:
- Export your private and public keys to `.gpg-export/`
- Display instructions for adding to GitHub Secrets
- Show the public key for Terraform Registry

### Step 4: Add to GitHub Secrets

1. Go to: https://github.com/kodflow/n8n/settings/secrets/actions

2. Add these secrets:
   - `GPG_PRIVATE_KEY`: Copy content from `.gpg-export/private-key.asc`
   - `GPG_PASSPHRASE`: Leave empty (no passphrase)

### Step 5: Upload to Terraform Registry

1. Go to: https://registry.terraform.io/settings/gpg-keys
2. Click "Add a GPG key"
3. Copy content from `.gpg-export/public-key.asc`
4. Click "Add GPG key"

### Step 6: Clean Up

```bash
make gpg/clean
```

This removes the sensitive `.gpg-export/` directory from your filesystem.

## Verification

Check your setup:

```bash
make gpg/info        # Show GPG configuration
make gpg/test        # Test signing
make gpg/fingerprint # Show fingerprint
```

## All GPG Commands

```bash
make gpg/generate        # Generate GPG key (may not work in all environments)
make gpg/configure       # Configure git to use GPG
make gpg/export          # Export keys for GitHub + Terraform Registry
make gpg/export/public   # Export only public key
make gpg/info            # Show current configuration
make gpg/test            # Test GPG signing
make gpg/fingerprint     # Show key fingerprint
make gpg/delete          # Delete GPG key
make gpg/clean           # Remove .gpg-export/ directory
make gpg/setup           # Full automated setup (generate + configure + export)
```

## Troubleshooting

### "gpg: agent_genkey failed: Forbidden"

This error occurs when GPG agent doesn't allow key generation without a passphrase in batch mode. Use the manual method above (`gpg --full-gen-key`).

### "No GPG key found"

Run `make gpg/info` to check if a key exists. If not, generate one manually as described above.

### Git not signing commits

Run:
```bash
git config --global --list | grep gpg
```

If you don't see `commit.gpgsign=true` and `tag.gpgsign=true`, run:
```bash
make gpg/configure
```

## For GitHub Actions

The release workflow (`.github/workflows/release.yml`) uses the GPG key automatically via:
- `GPG_PRIVATE_KEY` secret
- `GPG_PASSPHRASE` secret (empty)
- `GPG_FINGERPRINT` environment variable (auto-detected)

No additional configuration needed once secrets are added.
