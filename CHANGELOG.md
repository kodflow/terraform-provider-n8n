# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/).

---

## [Unreleased]


### ✅ Tests

- improve resource test coverage to 98.4% (`fa10123`)
- improve datasource and datasources test coverage to 100% (`462e068`)
- improve coverage with enhanced tests and shared utilities (`d19fa35`)
- achieve 94.3% coverage for credential package (`03d8697`)
- comprehensive test coverage improvements (`3344c2f`)

### 🚀 Features

- convert project model tests to table-driven (datasource/datasources) (`cb458fe`)
- convert datasource model tests to table-driven (tag/variable/user) (`52f6e50`)
- convert project resource model tests to table-driven format (`4881441`)
- convert project user resource model tests to table-driven (`b2d30a2`)
- convert execution item model tests to table-driven format (`af3febf`)
- convert execution datasources model tests to table-driven (`d5c16d1`)
- convert execution retry resource model tests to table-driven (`ac4ef21`)
- convert execution datasource tests to table-driven format (`a3deeaf`)
- convert execution datasources tests to table-driven format (`81b13ce`)
- convert variable datasources tests to table-driven format (`317e43c`)
- convert tag datasources tests to table-driven format with error cases (`1109915`)
- convert tag datasource tests to table-driven format with error cases (`4cd7c8a`)
- convert tag resource tests to table-driven format with error cases (`42530a7`)
- convert execution models datasource tests to table-driven (`9294118`)
- convert credential tests to table-driven format (`d58c8d0`)
- add comprehensive unit tests - 78.1% to 97.7% coverage (`17b8e12`)
- add comprehensive unit tests - 71.3% to 98.8% coverage (`ae4f64e`)
- improve test coverage from 60% to 90.8% (`da7d597`)
- add automatic documentation generation system (`034cf4f`)
- integrate ktn-linter with golangci-lint and VSCode (`b828914`)
- implement complete SDK coverage with all datasources and resources (`bbec535`)
- achieve 100% test coverage and add OpenAPI download command (`cabbb71`)
- integrate ktn-linter and enhance code quality standards (`ad76c4a`)
- add comprehensive code formatting support (`bd70c79`)
- add automatic semantic versioning and translate to English (`33dc716`)
- migration vers Bazel 9 et structure Terraform provider (`807372a`)

### 🔨 Chore

- ignore repomix configuration and output files (`99cb808`)
- apply formatting and gazelle updates (`a140452`)
- remove temporary build and analysis scripts (`9eab33c`)
- upgrade Go version from 1.24 to 1.25.3 (`6cd7cec`)
- add CodeRabbit configuration for automated PR reviews (`4331f4f`)

### 🐛 Bug Fixes

- correct workflow test mocks to match SDK requirements (`33a6035`)
- remove duplicate test definition in src/BUILD.bazel (`314f1b3`)
- enhance commit-msg hook to block AI attribution (`b99ff30`)
- allow tests to fail in coverage generation (`e6cbf19`)
- correct coverage calculation in generate-coverage.sh (`0d74bd2`)
- remove executable permission from README.md (`4607bcb`)
- add hook permissions fix script for GUI clients (`46ccbe2`)
- resolve critical linter issues across provider codebase (`ec19363`)
- update golangci-lint config to v2 format and integrate ktn-linter (`9f4f36d`)
- correct Docker tag format in devcontainer workflow (`36e97fd`)

### 📚 Documentation

- update coverage report with workflow package at 96.8% (`318ad40`)
- add strict testing and quality standards (`eac2a36`)
- improve AI assistant configuration with stricter directives (`dd4f419`)
- add AI assistant configuration document (`9d0c27d`)
- add coverage report (`f200924`)
- enhance code documentation and test coverage (`f41642a`)

### ♻️ Refactoring

- rename n8n_lib to src_lib for gazelle compatibility (`84cc3dc`)
- merge test files and clean up per linter requirements (`1467db5`)
- optimize code for Go 1.25.4 and improve pre-push hook (`160308b`)
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

- **Total commits:** 71
- **Features:** 26
- **Bug fixes:** 10
- **Tests:** 5
- **Refactoring:** 15
- **Test coverage:** N/A

### 👥 Contributors

- Florent <contact@making.codes>

---

*Changelog generated automatically on 2025-11-10 15:17:44*
