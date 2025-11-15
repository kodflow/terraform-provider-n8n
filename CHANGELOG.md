# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/).

---

## [Unreleased]


### 📝 Other

- Revert "feat(build): add GPG key management commands for Terraform Registry" (`b3e6efa`)

### 🚀 Features

- add GPG signing for releases and Terraform Registry support (`1ad5cef`)
- enforce GPG signature on all commits (`163846a`)
- add GPG verification and status display (`2f7474a`)
- add GPG key mounting and auto-configuration (`49819b4`)
- add GPG key management commands for Terraform Registry (`ed3964b`)
- add Terraform Registry publication support (`459b9b1`)
- Complete CI/CD pipeline with comprehensive testing and coverage (#4) (`532f02f`)

### 🔨 Chore

- remove deprecated IntelliCode extension (`5f383bc`)
- 1.1.7 (`bf98ec6`)
- 1.1.7 (`58acea8`)
- 1.1.7 (`1cafd2b`)
- 1.1.7 (`8e1f7eb`)
- 1.1.7 (`7fd93c0`)
- 1.1.7 (`d1b86ba`)
- 1.1.7 (`0dd7164`)
- 1.1.6 (`6dcd237`)
- 1.1.5 (`d734f98`)
- 1.1.4 (`edaa354`)
- 1.1.3 (`ae95bdd`)
- 1.1.2 (`8ae9b5e`)
- 1.1.1 (`f97551c`)
- 1.1.0 (`62c7bb7`)
- 1.1.0 (`547f997`)
- 1.0.0 [skip ci] (`f642aec`)

### 🐛 Bug Fixes

- complete project migration and translate French comments (`2b147d5`)
- correct asset download patterns for attestation (`6db8a52`)
- use SEMANTIC_RELEASE_TOKEN for release creation (`ec5b6e6`)
- download release assets before attestation (`f714aef`)
- publish releases directly instead of draft (`82d13a7`)
- correct svu command syntax (`367b233`)
- update svu to v3.3.0 (`2cc194d`)
- use @semantic-release/github to create releases (`b2bce94`)
- remove @semantic-release/git plugin - core handles tags (`112317b`)
- remove empty message option from semantic-release git config (`7cfae9f`)
- use @semantic-release/git with empty assets to create tags only (`6c9c8a4`)
- add @semantic-release/github plugin to create tags (`30b79c0`)
- enable debug mode for semantic-release (`238380d`)
- disable cancel-in-progress to allow release completion (`7abff0a`)
- add release summary step to show created tag (`a32e809`)
- add debug output for GPG import status (`dc3bd37`)
- add timeout and concurrency controls to prevent stuck workflows (`98e551e`)
- enable GPG signing for tags when key import succeeds (`c4fb62f`)
- remove deprecated archives.format configuration (`388b010`)
- specify main path and re-enable GPG signing (`fd4292a`)
- remove GPG references from release template (`1d511c1`)
- enable persist-credentials and disable semantic-release GitHub plugin (`063c486`)
- temporarily disable GPG signing in GoReleaser (`37ef43e`)
- add continue-on-error to GPG import in release workflow (`8749856`)
- disable GPG signing temporarily to unblock semantic-release (`b598036`)
- correct GitHub Actions syntax error in semver workflow (`099c384`)
- add fallback for GPG import failure in semantic-release (`f34c13b`)
- make GPG_PASSPHRASE optional for keys without passphrase (`c75dce4`)
- configure GPG agent for automated key generation (`2e111dc`)
- add manual GPG setup guide and improve scripts (`c0e29a5`)
- remove [skip ci] from semantic-release to allow release workflow (`551df3f`)

### 📚 Documentation

- add clarifying comments to GPG import step (`5c43f59`)
- add comprehensive Terraform Registry publication guide (`cecd7f9`)

### ♻️ Refactoring

- update Terraform registry references to full provider name (`72b96f9`)
- migrate project name from n8n to terraform-provider-n8n (`cf7141f`)
- replace semantic-release with svu for versioning (`f3c43be`)
- merge workflows into unified release pipeline (`da8b57a`)
- remove release commits for clean PR-based workflow (`9a66c45`)

---

### 📊 Statistics

- **Total commits:** 63
- **Features:** 7
- **Bug fixes:** 31
- **Tests:** 1
- **Refactoring:** 5
- **Test coverage:** N/A

### 👥 Contributors

- Kodflow <133899878+kodflow@users.noreply.github.com>
- semantic-release-bot <semantic-release-bot@martynus.net>

---

*Changelog generated automatically on 2025-11-15 15:26:02*
