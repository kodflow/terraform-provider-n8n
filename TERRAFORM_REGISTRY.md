# Publishing to Terraform Registry

This guide explains how to publish the n8n Terraform provider to the Terraform Registry.

## Prerequisites

Before publishing, ensure you have:

1. A GitHub account with the `kodflow/n8n` repository
2. Admin access to the repository
3. GPG key configured for signing releases

## Step 1: Generate and Configure GPG Key

The Terraform Registry requires all releases to be GPG signed.

**See [GPG_SETUP.md](GPG_SETUP.md) for detailed instructions.**

Quick summary:

```bash
# 1. Generate key manually (automated generation may fail in some environments)
gpg --full-gen-key
# Select: ECC (sign only), Curve 25519, no expiration, no passphrase

# 2. Configure git signing
make gpg/configure

# 3. Export keys for GitHub and Terraform Registry
make gpg/export

# 4. Test
make gpg/test
```

## Step 2: Configure GitHub Secrets

After running `make gpg/export`, you'll have keys in `.gpg-export/`:

1. Go to: https://github.com/kodflow/n8n/settings/secrets/actions

2. Add these secrets:
   - `GPG_PRIVATE_KEY`: Content of `.gpg-export/private-key.asc`
   - `GPG_PASSPHRASE`: Leave empty (key has no passphrase)

3. **IMPORTANT**: Delete the `.gpg-export/` directory after adding secrets:
   ```bash
   make gpg/clean
   ```

## Step 3: Upload Public Key to Terraform Registry

1. Go to: https://registry.terraform.io/settings/gpg-keys

2. Click "Add a GPG key"

3. Paste the content of `.gpg-export/public-key.asc` (or regenerate with `make gpg/export/public`)

4. Click "Add GPG key"

## Step 4: Create First Release

### Verify Current Status

```bash
# Check current branch and tags
git status
git tag -l

# Ensure you're on main with latest changes
git checkout main
git pull
```

### Create and Push Tag

The release workflow triggers automatically when you push a version tag:

```bash
# Create annotated, signed tag
git tag -s v1.0.0 -m "First stable release"

# Push tag to trigger release
git push origin v1.0.0
```

### What Happens Next

The GitHub Actions workflow (`.github/workflows/release.yml`) will automatically:

1. Build binaries for multiple platforms:
   - Linux (amd64, 386, arm, arm64)
   - macOS (amd64, arm64)
   - Windows (amd64, 386)
   - FreeBSD (amd64, 386, arm, arm64)

2. Generate SHA256 checksums

3. Sign checksums with GPG

4. Create GitHub Release with all artifacts

5. Include `terraform-registry-manifest.json` in the release

## Step 5: Publish to Terraform Registry

1. Go to: https://registry.terraform.io/publish/provider

2. Sign in with your GitHub account

3. Select `kodflow/n8n` repository

4. Accept the terms of service

5. Click "Publish"

The Terraform Registry will:
- Detect the `v1.0.0` tag
- Verify GPG signatures
- Parse `terraform-registry-manifest.json`
- Index documentation from `docs/`
- Make the provider available at `registry.terraform.io/kodflow/n8n`

## Step 6: Verify Publication

After a few minutes, check:

1. Provider page: https://registry.terraform.io/providers/kodflow/n8n

2. Documentation is visible and correctly formatted

3. All versions are listed

4. Download counts start tracking

## Step 7: Update README

Once published, update the README to remove "Coming Soon":

```hcl
terraform {
  required_providers {
    n8n = {
      source  = "kodflow/n8n"
      version = "~> 1.0"
    }
  }
}
```

## Troubleshooting

### GPG Signing Fails

```bash
# Check GPG key info
make gpg/info

# Test signing
make gpg/test

# View fingerprint for GoReleaser
make gpg/fingerprint
```

### Release Workflow Fails

Check GitHub Actions logs:
- https://github.com/kodflow/n8n/actions

Common issues:
- Missing `GPG_PRIVATE_KEY` secret
- Incorrect `GPG_PASSPHRASE`
- Invalid tag format (must be `vX.Y.Z`)

### Registry Rejects Release

The Terraform Registry requires:
- Valid semantic version tag (`v1.0.0`, not `1.0.0`)
- GPG-signed checksums file
- `terraform-registry-manifest.json` in release assets
- Documentation in `docs/` directory

Verify release assets contain:
- `terraform-provider-n8n_X.Y.Z_SHA256SUMS`
- `terraform-provider-n8n_X.Y.Z_SHA256SUMS.sig`
- `terraform-registry-manifest.json`
- Binary archives for each platform

## Subsequent Releases

For future releases, the process is simpler:

1. Merge features to `main`
2. Create and push new tag:
   ```bash
   git tag -s v1.1.0 -m "Release v1.1.0"
   git push origin v1.1.0
   ```
3. The Terraform Registry automatically detects and publishes new versions

## Using Semantic Release (Alternative)

If you prefer automated versioning:

1. Merge PR to `main` with conventional commits
2. Semantic release workflow runs automatically
3. Tag is created based on commit types (feat, fix, etc.)
4. Release workflow triggers on new tag
5. Terraform Registry auto-updates

## License Note

The provider uses the **Sustainable Use License 1.0**, which may be listed as "Custom" or "Other" on the Terraform Registry. This is expected and doesn't prevent publication.

## Support

For issues with publication:
- Terraform Registry: https://discuss.hashicorp.com/c/terraform-providers
- Provider Issues: https://github.com/kodflow/n8n/issues
