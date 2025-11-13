# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/).

---

## [Unreleased]


### ✅ Tests

- achieve 100% coverage for all datasource Read methods (`1c85a71`)
- improve datasource Read coverage from 0% to ~95% (`28de8a4`)
- achieve 100% coverage for transfer and project user resources (`0858404`)
- improve coverage and fix linting issues (`e331ca1`)
- improve secondary resources coverage from 0% to 100% (`1c6d6d8`)
- improve execution retry_resource coverage to 100% (except Delete) (`220d73d`)
- improve CRUD coverage to 83.6% with execute*Logic pattern (`298bfce`)
- improve Create and Update test coverage with helper extraction (`200234b`)
- improve credential resource coverage with extractable helpers (`2605765`)
- improve project Update method coverage to 100% (`40a5e3e`)
- improve Update method test coverage to 100% for variable (`0c9a124`)
- add Update method coverage tests (`d9c6c52`)
- improve CRUD method coverage (`429898d`)
- achieve 100% coverage for Delete method (`249516e`)
- achieve 100% coverage for pull_resource CRUD methods (`d23a600`)
- complete coverage to 85.3% and improve reporting (`145785b`)
- complete workflow CRUD functions to 100% coverage (`fc5aa03`)
- rename coverage test files to match linter requirements (`c2e5d3c`)
- add comprehensive coverage tests for credential and user modules (`49a168a`)
- refactor tests to table-driven format with comprehensive coverage (`7ca6910`)
- refactor all tests to table-driven format with error cases (`1ff7283`)
- add tests for private helper functions in project package (`d589824`)
- convert project user_resource_test to table-driven format (`4ae9963`)
- convert project datasource_test to table-driven format (`64b4723`)
- convert variable helpers_test to table-driven format (`5ccc2be`)
- convert credential workflow_backup_test to table-driven format (`c8c2f31`)
- convert tag helpers_test to table-driven format (`5b966d5`)
- convert user helpers_test to table-driven format (`17667d0`)
- convert credential models transfer_resource_test to table-driven (`fdab61f`)
- convert credential models resource_test to table-driven format (`adaa98b`)
- convert shared provider_test to table-driven format (`268ab0e`)
- add wantErr to tag datasource_internal_test stubs (`4dda43d`)
- add wantErr to project/sourcecontrol internal test stubs (`feb6b2b`)
- convert project helpers_test to table-driven format (`2a3d2bd`)
- add wantErr to execution/project internal test stubs (`4cd2e88`)
- convert workflow resource_internal_test to table-driven format (`b026f0b`)
- add wantErr fields to user internal test stubs (`82b60c6`)
- add wantErr fields to execution internal test stubs (`2c4ab0b`)
- convert credential resource_test validation+usecases (`405ccca`)
- add wantErr to internal test stubs (4 files) (`3b2086c`)
- convert constants_test to table-driven (`43cde9e`)
- add wantErr field to shared pointers_test (`70c0426`)
- convert execution helpers_test to table-driven (`23730a2`)
- improve resource test coverage to 98.4% (`fa10123`)
- improve datasource and datasources test coverage to 100% (`462e068`)
- improve coverage with enhanced tests and shared utilities (`d19fa35`)
- achieve 94.3% coverage for credential package (`03d8697`)
- comprehensive test coverage improvements (`3344c2f`)

### 🚀 Features

- add environment variable support for provider configuration (`840d689`)
- reorganize report by semantic categories (`c8a22d0`)
- improve coverage report with public function details (`2881515`)
- convert workflow external tests to table-driven format (`81c9f30`)
- convert sourcecontrol/models/resource_test.go to table-driven (`747010e`)
- convert project/models/item_test.go to table-driven format (`24cd91b`)
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

- remove TODO comment from credential test (`b595868`)
- ignore repomix configuration and output files (`99cb808`)
- apply formatting and gazelle updates (`a140452`)
- remove temporary build and analysis scripts (`9eab33c`)
- upgrade Go version from 1.24 to 1.25.3 (`6cd7cec`)
- add CodeRabbit configuration for automated PR reviews (`4331f4f`)

### 🐛 Bug Fixes

- add sudo for Bazelisk installation (`369cfbe`)
- make acceptance tests non-blocking in pipeline (`e193d10`)
- improve secondary resource titles formatting (`81d3f78`)
- detect empty functions in coverage report (🔴 → 🔵) (`29aa9cf`)
- correct test function names for private functions (`e98f32e`)
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

- extract ElementsAs to improve testability (`9972ff2`)
- extract transfer Create logic for testability (`a7a7f7b`)
- move CRUD tests to external test file (`d6eec56`)
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

- update CI/CD pipeline with complete test suite and E2E testing (`e33a9fc`)
- add comprehensive CI workflow and commit validation system (`a3a0c9d`)

---

### 📊 Statistics

- **Total commits:** 130
- **Features:** 32
- **Bug fixes:** 15
- **Tests:** 48
- **Refactoring:** 18
- **Test coverage:** N/A

### 👥 Contributors

- Florent <contact@making.codes>

---

*Changelog generated automatically on 2025-11-13 11:01:30*
