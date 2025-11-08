# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/).

---

## [Unreleased]


### ✅ Tests

- comprehensive test coverage improvements (`5d3466b`)

### 🚀 Features

- integrate ktn-linter with golangci-lint and VSCode (`c64ba58`)
- implement complete SDK coverage with all datasources and resources (`389623d`)
- achieve 100% test coverage and add OpenAPI download command (`f1477c4`)
- integrate ktn-linter and enhance code quality standards (`3306966`)
- add comprehensive code formatting support (`855943e`)
- add automatic semantic versioning and translate to English (`33dc716`)
- migration vers Bazel 9 et structure Terraform provider (`807372a`)

### 🔨 Chore

- remove temporary build and analysis scripts (`2706153`)
- upgrade Go version from 1.24 to 1.25.3 (`35756d9`)
- add CodeRabbit configuration for automated PR reviews (`4331f4f`)

### 🐛 Bug Fixes

- resolve critical linter issues across provider codebase (`5abf916`)
- update golangci-lint config to v2 format and integrate ktn-linter (`9f4f36d`)
- correct Docker tag format in devcontainer workflow (`36e97fd`)

### 📚 Documentation

- add coverage report (`2dd13bb`)
- enhance code documentation and test coverage (`868bd49`)

### ♻️ Refactoring

- reorganize models into domain subdirectories (`9868fed`)
- reorganize models into dedicated subdirectories with simplified naming (`d407dcd`)
- fix all ktn-linter violations and improve code organization (`442ec00`)
- improve struct naming consistency across domains (`d7137bf`)
- migrate to DDD architecture and improve code quality (`77c7ed3`)
- fix all ktn-linter errors and improve code quality (`9a56942`)
- improve code quality and documentation across provider (`a532143`)
- consolidate SDK build pipeline and enhance Makefile formatting (`9076518`)
- simplify OpenAPI patch workflow (`2cb717f`)
- improve Makefile output formatting and organization (`5d6b1f4`)
- remove redundant commit convention section from PR template (`aedabcd`)

### 🔧 Build

- filter KTN-STRUCT-005 false positives in Makefile (`a628fcf`)
- add ktn-linter filter script to exclude false positives (`b76b279`)
- integrate ktn-linter with golangci-lint for strict code quality enforcement (`a7febe1`)

---

### 📊 Statistics

- **Total commits:** 30
- **Features:** 7
- **Bug fixes:** 3
- **Tests:** 1
- **Refactoring:** 11
- **Test coverage:** 70.9%

### 👥 Contributors

- Florent <contact@making.codes>

---

*Changelog generated automatically on 2025-11-08 13:40:46*
