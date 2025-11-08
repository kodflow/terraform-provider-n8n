# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/).

---

## [Unreleased]


### ✅ Tests

- comprehensive test coverage improvements (`3344c2f`)

### 🚀 Features

- add automatic documentation generation system (`034cf4f`)
- integrate ktn-linter with golangci-lint and VSCode (`b828914`)
- implement complete SDK coverage with all datasources and resources (`bbec535`)
- achieve 100% test coverage and add OpenAPI download command (`cabbb71`)
- integrate ktn-linter and enhance code quality standards (`ad76c4a`)
- add comprehensive code formatting support (`bd70c79`)
- add automatic semantic versioning and translate to English (`33dc716`)
- migration vers Bazel 9 et structure Terraform provider (`807372a`)

### 🔨 Chore

- remove temporary build and analysis scripts (`9eab33c`)
- upgrade Go version from 1.24 to 1.25.3 (`6cd7cec`)
- add CodeRabbit configuration for automated PR reviews (`4331f4f`)

### 🐛 Bug Fixes

- add hook permissions fix script for GUI clients (`46ccbe2`)
- resolve critical linter issues across provider codebase (`ec19363`)
- update golangci-lint config to v2 format and integrate ktn-linter (`9f4f36d`)
- correct Docker tag format in devcontainer workflow (`36e97fd`)

### 📚 Documentation

- add coverage report (`f200924`)
- enhance code documentation and test coverage (`f41642a`)

### ♻️ Refactoring

- improve Makefile formatting and documentation targets (`353658c`)
- reorganize models into domain subdirectories (`8ae34eb`)
- reorganize models into dedicated subdirectories with simplified naming (`de9445d`)
- fix all ktn-linter violations and improve code organization (`edb1c95`)
- improve struct naming consistency across domains (`6b5b584`)
- migrate to DDD architecture and improve code quality (`1a6d759`)
- fix all ktn-linter errors and improve code quality (`0cc97c8`)
- improve code quality and documentation across provider (`1099c8c`)
- consolidate SDK build pipeline and enhance Makefile formatting (`0a0288d`)
- simplify OpenAPI patch workflow (`44353ca`)
- improve Makefile output formatting and organization (`5b8346d`)
- remove redundant commit convention section from PR template (`aedabcd`)

### 🔧 Build

- filter KTN-STRUCT-005 false positives in Makefile (`a50c21f`)
- add ktn-linter filter script to exclude false positives (`9a2f112`)
- integrate ktn-linter with golangci-lint for strict code quality enforcement (`a7febe1`)

### 🤖 CI/CD

- add comprehensive CI workflow and commit validation system (`a3a0c9d`)

---

### 📊 Statistics

- **Total commits:** 34
- **Features:** 8
- **Bug fixes:** 4
- **Tests:** 1
- **Refactoring:** 12
- **Test coverage:** N/A

### 👥 Contributors

- Florent <contact@making.codes>

---

*Changelog generated automatically on 2025-11-08 15:32:05*
