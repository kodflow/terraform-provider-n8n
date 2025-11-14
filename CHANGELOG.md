# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/).

---

## [Unreleased]


### 📝 Other

- update patch (`4c1baa0`)
- prepare patch (`1740fd3`)
- clear (`c67814a`)
- save (`fd769a4`)
- ok (`d88edcd`)
- bump 1.119.1 (`d3c16b8`)
- save (`c444667`)
- save (`15c0d6f`)
- rework (`450d9b4`)
- openapi (`2d59b48`)

### ✅ Tests

- remove useless empty test cases (`e7b7ec7`)
- remove copyright verification test files (`5e7b367`)
- verify copyright hook with non-test file (`2449ceb`)
- verify pre-commit hook adds copyright headers (`4e0b277`)
- add real API verification for E2E tests (`44c8c55`)

### 🚀 Features

- add copyright headers to remaining source files (`78d95da`)
- add Sustainable Use License and copyright headers (`c80c828`)
- add acceptance tests to coverage report (`c1b942b`)
- complete CI/CD pipeline with comprehensive testing and coverage (`b0089c3`)

### 🔨 Chore

- update generated documentation and OpenAPI patch (`ac6b874`)
- move git hooks to .github/hooks (`837f3a6`)
- bump n8n OpenAPI spec to 1.119.1 (`dc344a6`)
- bump n8n OpenAPI spec to 1.119.1 (`667e26a`)
- bump n8n OpenAPI spec to 1.119.1 (`4bf8906`)
- bump n8n OpenAPI spec to 1.119.1 (`015f4b7`)
- bump n8n OpenAPI spec to 1.119.1 (`c40e349`)
- standardize environment variables across workflows and docs (`2dad579`)

### 🐛 Bug Fixes

- auto-detect current branch in changelog generation (`398ea12`)
- regenerate model_workflow.go with all required fields (`26e3e39`)
- add missing make sdk step to generate Go files (`bac4bf3`)
- move dependency download after SDK generation (`2bc27a8`)
- skip automatic commit in CI environment (`4ef0aa6`)
- remove unnecessary stderr redirection in test/acceptance (`06206e1`)
- add additionalProperties to credential schema for API compatibility (`8742d1f`)
- resolve unmarshaling errors for nested workflow structures (`82c7a28`)
- add missing workflow fields and handle nullable types (`c6057cd`)
- align OpenAPI spec with actual n8n API implementation (`4930de7`)
- make shell scripts executable (`23c6808`)
- make copyright header script executable (`192d5f6`)
- force prettier to respect .prettierignore for SDK files (`160f155`)
- exclude SDK from prettier to preserve ggignore comments (`72d43c3`)
- mark example secrets with ggignore and fix GitGuardian config (`71eb244`)
- correct GitGuardian configuration syntax (`6241962`)
- exclude auto-generated docs from uncommitted changes check (`01e142d`)
- remove t.Parallel() from TestOsExitVariable to avoid race condition (`cdfc5b2`)

### 📚 Documentation

- optimize README with concise feature overview (`44133b1`)
- cleanup markdown files and update README (`8d05b47`)
- translate all French content to English (`b1200ff`)

### ♻️ Refactoring

- move openapi.yaml reset from pre-commit to make clean (`2f08633`)
- separate unit and E2E tests in Makefile (`a3844d6`)
- separate script and patch responsibilities (`eccec80`)
- remove legacy environment variable support (`6ec444a`)
- use official actions with SHA pinning and centralize in Makefile (`d9a1ecd`)
- remove redundant model unit tests (`0dcc564`)
- add blank line to separate stdlib and external imports (`b57452b`)

### 🤖 CI/CD

- remove SDK generation from workflow to avoid commit conflicts (`38dfe43`)
- run devcontainer build only on main branch push (`02b0d93`)
- add GitGuardian configuration to ignore false positives (`b8505a3`)
- add git diff output to uncommitted changes check for debugging (`d85a2e9`)

---

### 📊 Statistics

- **Total commits:** 59
- **Features:** 4
- **Bug fixes:** 18
- **Tests:** 6
- **Refactoring:** 7
- **Test coverage:** N/A

### 👥 Contributors

- Kodflow <133899878+kodflow@users.noreply.github.com>

---

*Changelog generated automatically on 2025-11-14 14:03:17*
