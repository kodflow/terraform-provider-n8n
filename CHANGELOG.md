# Changelog

Toutes les modifications notables de ce projet seront documentées dans ce fichier.

Le format est basé sur [Keep a Changelog](https://keepachangelog.com/fr/1.0.0/),
et ce projet adhère au [Semantic Versioning](https://semver.org/lang/fr/).

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
