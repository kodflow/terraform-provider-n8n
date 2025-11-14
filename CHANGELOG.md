# Changelog

Toutes les modifications notables de ce projet seront documentées dans ce fichier.

Le format est basé sur [Keep a Changelog](https://keepachangelog.com/fr/1.0.0/),
et ce projet adhère au [Semantic Versioning](https://semver.org/lang/fr/).

## [1.1.0](https://github.com/kodflow/n8n/compare/v1.0.0...v1.1.0) (2025-11-14)

### 🚀 Features

* **devcontainer:** add GPG key mounting and auto-configuration ([49819b4](https://github.com/kodflow/n8n/commit/49819b48bcb5dc9f2df2314453dfc10b5c5f6380))
* **devcontainer:** add GPG verification and status display ([2f7474a](https://github.com/kodflow/n8n/commit/2f7474acbf2329af7685a6110b43078aba06cdb3))
* **hooks:** enforce GPG signature on all commits ([163846a](https://github.com/kodflow/n8n/commit/163846aad5d005b5fa906b882950729107da59fc))
* **provider:** add Terraform Registry publication support ([459b9b1](https://github.com/kodflow/n8n/commit/459b9b15494a54cd85e189a37e6b4ce1d0b671f2))
* **release:** add GPG signing for releases and Terraform Registry support ([1ad5cef](https://github.com/kodflow/n8n/commit/1ad5cef6fb40d45a0a7679c684917ba922690c9c))

### 🐛 Bug Fixes

* **ci:** remove [skip ci] from semantic-release to allow release workflow ([551df3f](https://github.com/kodflow/n8n/commit/551df3fbd2fa22142f19f08b5ed3b7663ae6dd4b))
* **devcontainer:** configure GPG agent for automated key generation ([2e111dc](https://github.com/kodflow/n8n/commit/2e111dc8c7b8e35f9381212e2caf6893de3dbd9f))
* **gpg:** add manual GPG setup guide and improve scripts ([c0e29a5](https://github.com/kodflow/n8n/commit/c0e29a5862b2c4a049f75b64a4d32615c7fd0a3a))
* **workflows:** make GPG_PASSPHRASE optional for keys without passphrase ([c75dce4](https://github.com/kodflow/n8n/commit/c75dce47a31f649c5d96495a0fca5ac156aa7207))

### 📚 Documentation

* **provider:** add comprehensive Terraform Registry publication guide ([cecd7f9](https://github.com/kodflow/n8n/commit/cecd7f9ec171b057a141cb427ca4a94352bf7b36))

## 1.0.0 (2025-11-14)

### 🚀 Features

* Complete CI/CD pipeline with comprehensive testing and coverage ([#4](https://github.com/kodflow/n8n/issues/4)) ([532f02f](https://github.com/kodflow/n8n/commit/532f02f25555a832c118cbae15c57b7ed575282c))

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

### 🐛 Bug Fixes

- configure GPG agent for automated key generation (`2e111dc`)
- add manual GPG setup guide and improve scripts (`c0e29a5`)
- remove [skip ci] from semantic-release to allow release workflow (`551df3f`)

### 📚 Documentation

- add comprehensive Terraform Registry publication guide (`cecd7f9`)

---

### 📊 Statistics

- **Total commits:** 11
- **Features:** 6
- **Bug fixes:** 3
- **Tests:** 0
- **Refactoring:** 0
- **Test coverage:** N/A

### 👥 Contributors

- Kodflow <133899878+kodflow@users.noreply.github.com>

---

*Changelog generated automatically on 2025-11-14 22:24:51*
