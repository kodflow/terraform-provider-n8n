## Description

<!-- Clearly describe the changes made by this PR -->

## Type of Change

<!-- Check the appropriate box by replacing [ ] with [x] -->

- [ ] 🐛 **fix**: Bug fix (patch version)
- [ ] 🚀 **feat**: New feature (minor version)
- [ ] ⚡ **perf**: Performance improvement (patch version)
- [ ] ♻️ **refactor**: Code refactoring (patch version)
- [ ] 🔧 **build**: Build system changes (patch version)
- [ ] 📚 **docs**: Documentation only (no release)
- [ ] ✅ **test**: Adding or modifying tests (no release)
- [ ] 🧹 **chore**: Maintenance (no release)
- [ ] 👷 **ci**: CI/CD changes (no release)
- [ ] 💥 **BREAKING CHANGE**: Breaking compatibility change (major version)

## Checklist

<!-- Verify that all the following points are met -->

- [ ] My code follows the project conventions
- [ ] I have performed a self-review of my code
- [ ] I have commented the code in areas that are difficult to understand
- [ ] I have updated the documentation if necessary
- [ ] My changes do not generate new warnings
- [ ] I have added tests that prove my fix works or my feature works
- [ ] Unit tests pass locally (`make test`)
- [ ] My PR title follows the format: `<type>: <short description>`

## Tests Performed

<!-- Describe the tests you performed -->

```bash
# Example test commands
make test
make build
cd sample/ && terraform init && terraform plan
```

## Commit Convention

This PR follows [Conventional Commits](https://www.conventionalcommits.org/).

The **PR title** must be in the format: `<type>(<scope>): <description>`

### Valid title examples:
- `feat: add n8n workflows support`
- `fix: correct credentials parsing error`
- `feat(workflows)!: change workflows API` (breaking change)
- `docs: update README with examples`

### Impact on versioning:
- `fix:`, `perf:`, `refactor:`, `build:` → **Patch** (0.1.0 → 0.1.1)
- `feat:` → **Minor** (0.1.0 → 0.2.0)
- `BREAKING CHANGE:` or `!` → **Major** (0.1.0 → 1.0.0)
- `docs:`, `test:`, `chore:`, `ci:` → **No release**

## Additional Notes

<!-- Any other relevant information -->
