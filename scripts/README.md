# Scripts de Documentation Automatique

Ce répertoire contient des scripts pour générer automatiquement la documentation du projet.

## 📝 generate-changelog.sh

### Description

Génère automatiquement le fichier `CHANGELOG.md` basé sur l'historique Git en utilisant les conventions **Conventional Commits**.

### Utilisation

```bash
# Générer le changelog pour la branche courante
./scripts/generate-changelog.sh

# Spécifier une branche source et une branche de base
./scripts/generate-changelog.sh feat/ma-branche main
```

### Ou via Makefile

```bash
make changelog
```

### Format des Commits

Le script reconnaît les types de commits suivants :

| Type | Emoji | Catégorie | Exemple |
|------|-------|-----------|---------|
| `feat:` | 🚀 | Features | `feat: add new resource` |
| `fix:` | 🐛 | Bug Fixes | `fix: resolve nil pointer` |
| `test:` | ✅ | Tests | `test: add unit tests` |
| `docs:` | 📚 | Documentation | `docs: update README` |
| `refactor:` | ♻️ | Refactoring | `refactor: simplify logic` |
| `perf:` | ⚡ | Performance | `perf: optimize query` |
| `build:` | 🔧 | Build | `build: update Bazel` |
| `ci:` | 🤖 | CI/CD | `ci: add workflow` |
| `chore:` | 🔨 | Chore | `chore: update deps` |
| `style:` | 💄 | Style | `style: format code` |

### Fonctionnalités

- ✅ Catégorisation automatique par type de commit
- ✅ Hash courts pour traçabilité
- ✅ Statistiques (nombre de commits par type)
- ✅ Intégration avec COVERAGE.MD (affiche le taux de couverture)
- ✅ Liste des contributeurs
- ✅ Timestamp de génération

### Sortie

Le fichier généré suit le format [Keep a Changelog](https://keepachangelog.com/):

```markdown
# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### 🚀 Features

- comprehensive test coverage improvements (`5d3466b`)

### 🐛 Bug Fixes

- resolve critical linter issues (`5abf916`)

---

### 📊 Statistics

- **Total commits:** 30
- **Features:** 7
- **Test coverage:** 70.9%

### 👥 Contributors

- Florent <contact@making.codes>
```

## 🔄 Automatisation

### Hook Git Pre-commit

Pour générer automatiquement le changelog avant chaque commit, créez un hook :

```bash
# Créer le hook
cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash
# Auto-generate changelog if commits have changed
./scripts/generate-changelog.sh > /dev/null 2>&1
git add CHANGELOG.md
EOF

chmod +x .git/hooks/pre-commit
```

### GitHub Actions

Exemple de workflow pour générer le changelog dans CI/CD :

```yaml
name: Update Documentation

on:
  push:
    branches: [ main, develop ]

jobs:
  docs:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
        with:
          fetch-depth: 0  # Important pour l'historique Git complet

      - name: Generate Changelog
        run: |
          chmod +x ./scripts/generate-changelog.sh
          ./scripts/generate-changelog.sh

      - name: Commit changes
        run: |
          git config user.name "GitHub Actions"
          git config user.email "actions@github.com"
          git add CHANGELOG.md
          git commit -m "docs: update changelog [skip ci]" || true
          git push
```

## 📊 Makefile Integration

Le Makefile fournit des commandes pratiques :

```bash
# Générer uniquement le changelog
make changelog

# Générer le rapport de couverture
make coverage-report

# Générer toute la documentation
make docs
```

## 🎯 Bonnes Pratiques

1. **Commits conventionnels** : Utilisez toujours le format `type: description`
2. **Génération régulière** : Exécutez `make changelog` avant chaque PR
3. **Review** : Vérifiez le changelog généré pour cohérence
4. **Versioning** : Mettez à jour `[Unreleased]` en version release lors des tags

## 📚 Références

- [Conventional Commits](https://www.conventionalcommits.org/)
- [Keep a Changelog](https://keepachangelog.com/)
- [Semantic Versioning](https://semver.org/)
