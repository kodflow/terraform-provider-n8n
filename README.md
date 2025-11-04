# Terraform Provider for n8n

Provider Terraform pour gérer les ressources n8n (workflows, credentials, etc.).

[![Bazel](https://img.shields.io/badge/Build-Bazel%209.0-43A047?logo=bazel)](https://bazel.build/)
[![Go](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go)](https://go.dev/)
[![Terraform](https://img.shields.io/badge/Terraform-Plugin%20Framework-7B42BC?logo=terraform)](https://developer.hashicorp.com/terraform/plugin/framework)

## Table des matières

- [Prérequis](#prérequis)
- [Installation](#installation)
- [Développement](#développement)
- [Build et Tests](#build-et-tests)
- [Structure du projet](#structure-du-projet)
- [Publication](#publication)
- [Contributing](#contributing)
- [License](#license)

## Prérequis

### Versions requises

- **Go 1.24.0+** (requis par terraform-plugin-framework v1.16+)
- **Bazel 9.0+** (système de build)
- **Terraform 1.0+** ou **OpenTofu 1.0+**

### DevContainer (Recommandé)

Le projet est configuré avec un DevContainer incluant tous les outils nécessaires:

- **Go 1.25.3** (compatible 1.24+)
- **Bazel 9.0.0rc1** (via Bazelisk)
- **Terraform & OpenTofu** (pré-installés)
- Extensions VS Code:
  - `golang.go` - Support Go officiel
  - `hashicorp.terraform` - Support Terraform
  - `BazelBuild.vscode-bazel` - Support Bazel

**Pour utiliser le DevContainer:**
1. Ouvrir le projet dans VS Code
2. Accepter la proposition d'ouvrir dans le conteneur
3. Attendre la construction du conteneur (première fois uniquement)

### Installation manuelle

Si vous n'utilisez pas le DevContainer:

```bash
# Installer Go 1.24+
# Voir: https://go.dev/doc/install

# Installer Bazelisk (recommandé pour gérer les versions de Bazel)
go install github.com/bazelbuild/bazelisk@latest

# Vérifier les versions
go version        # doit afficher go1.24 ou supérieur
bazel version     # doit afficher Bazel 9.0+
```

## Installation

### Via Terraform Registry (À venir)

```hcl
terraform {
  required_providers {
    n8n = {
      source  = "kodflow/n8n"
      version = "~> 0.1.0"
    }
  }
}

provider "n8n" {
  api_url = "https://your-n8n-instance.com"
  api_key = var.n8n_api_key
}
```

### Installation locale pour développement

```bash
# Compiler et installer localement
make test

# Le provider sera installé dans:
# ~/.terraform.d/plugins/registry.terraform.io/kodflow/n8n/0.0.1/<OS>_<ARCH>/
```

## Développement

### Commandes Make disponibles

Le `Makefile` fournit deux commandes essentielles:

```bash
make help    # Affiche l'aide avec toutes les commandes disponibles
make test    # Lance les tests avec Bazel
```

### Configuration Bazel

Le projet utilise **Bazel 9** avec **bzlmod** (le nouveau système de gestion des dépendances):

- **`.bazelversion`**: Version de Bazel (9.0.0rc1)
- **`MODULE.bazel`**: Dépendances et configuration bzlmod
- **`BUILD.bazel`**: Configuration de build à la racine
- **`.bazelrc`**: Options de build Bazel

**Dépendances Bazel:**
- `rules_go v0.58.3` - Règles Go pour Bazel (avec support Bazel 9)
- `gazelle v0.46.0` - Générateur automatique de BUILD files
- `rules_proto v7.1.0` - Support Protocol Buffers
- `bazel_features v1.33.0` - Détection de fonctionnalités Bazel

### Architecture du projet

```
.
├── .bazelrc              # Configuration Bazel
├── .bazelversion         # Version Bazel (9.0.0rc1)
├── MODULE.bazel          # Dépendances bzlmod
├── BUILD.bazel           # Configuration build racine
├── go.mod                # Dépendances Go
├── Makefile              # Commandes de build
├── .devcontainer/        # Configuration DevContainer
│   ├── Dockerfile        # Image de développement
│   └── devcontainer.json # Configuration VS Code
├── src/                  # Code source du provider
│   ├── main.go           # Point d'entrée
│   ├── BUILD.bazel       # Configuration build src
│   └── internal/
│       └── provider/     # Implémentation du provider
│           ├── provider.go
│           ├── provider_test.go
│           └── BUILD.bazel
└── .github/
    └── workflows/        # CI/CD GitHub Actions
        └── release.yml   # Workflow de release automatique
```

## Build et Tests

### Lancer les tests

```bash
# Via Make (recommandé)
make test

# Directement avec Bazel
bazel test //src/...

# Tests avec verbose timeout warnings
bazel test --test_verbose_timeout_warnings //src/...
```

### Build du provider

```bash
# Build avec Bazel
bazel build //src:terraform-provider-n8n

# Le binaire sera disponible dans:
# bazel-bin/src/terraform-provider-n8n
```

### Nettoyage

```bash
# Nettoyer les artifacts Bazel
bazel clean

# Nettoyage complet (cache inclus)
bazel clean --expunge
```

## Structure du projet

### Code source

Le provider est structuré selon les best practices Terraform Plugin Framework:

```
src/
├── main.go                    # Point d'entrée du provider
└── internal/
    └── provider/
        ├── provider.go        # Implémentation du provider principal
        ├── provider_test.go   # Tests du provider
        └── BUILD.bazel        # Configuration build
```

### Configuration Terraform

Pour utiliser le provider en développement local:

```hcl
terraform {
  required_providers {
    n8n = {
      source  = "registry.terraform.io/kodflow/n8n"
      version = "0.0.1"
    }
  }
}

provider "n8n" {
  # Configuration du provider
}
```

## Publication

### Workflow de release

Le projet utilise **GoReleaser** via GitHub Actions pour automatiser les releases:

1. **Créer un tag**:
   ```bash
   git tag -a v0.1.0 -m "Release v0.1.0"
   git push origin v0.1.0
   ```

2. **GitHub Actions** déclenche automatiquement:
   - Compilation cross-platform (Linux, macOS, Windows, FreeBSD)
   - Génération des checksums SHA256
   - Signature GPG des checksums
   - Création d'une release GitHub avec les artifacts

### Configuration GPG

Pour signer les releases, configurer une clé GPG:

```bash
# Générer une clé GPG
gpg --full-generate-key
# Choisir: RSA and RSA, 4096 bits, pas d'expiration

# Exporter la clé privée
gpg --armor --export-secret-keys YOUR_EMAIL > private-key.asc

# Exporter la clé publique
gpg --armor --export YOUR_EMAIL
```

Ajouter les secrets GitHub (Settings > Secrets and variables > Actions):
- `GPG_PRIVATE_KEY`: Contenu de `private-key.asc`
- `GPG_PASSPHRASE`: Passphrase de la clé GPG

### Artifacts de release

GoReleaser génère automatiquement:

```
terraform-provider-n8n_0.1.0_darwin_amd64.zip
terraform-provider-n8n_0.1.0_darwin_arm64.zip
terraform-provider-n8n_0.1.0_linux_amd64.zip
terraform-provider-n8n_0.1.0_linux_arm64.zip
terraform-provider-n8n_0.1.0_windows_amd64.zip
terraform-provider-n8n_0.1.0_SHA256SUMS
terraform-provider-n8n_0.1.0_SHA256SUMS.sig
```

### Inscription sur le Registry

#### Terraform Registry (officiel)

1. Se connecter sur [registry.terraform.io](https://registry.terraform.io)
2. Aller dans "Publish" > "Provider"
3. Connecter le repository GitHub
4. Ajouter la clé GPG publique
5. Le registry détectera automatiquement les releases

#### OpenTofu Registry

OpenTofu utilise le même format. Suivre la documentation sur [github.com/opentofu/registry](https://github.com/opentofu/registry).

## Contributing

### Prérequis

1. Fork le repository
2. Cloner votre fork
3. Ouvrir dans VS Code avec DevContainer (recommandé)
4. Créer une branche pour vos modifications

### Workflow de contribution

```bash
# Créer une branche
git checkout -b feature/ma-fonctionnalite

# Faire vos modifications
# ...

# Tester
make test

# Commit et push
git add .
git commit -m "feat: ajout de ma fonctionnalité"
git push origin feature/ma-fonctionnalite

# Créer une Pull Request sur GitHub
```

### Standards de code

- **Go**: Suivre les conventions Go standards (`gofmt`, `golint`)
- **Commits**: Utiliser [Conventional Commits](https://www.conventionalcommits.org/)
  - `feat:` - Nouvelle fonctionnalité
  - `fix:` - Correction de bug
  - `docs:` - Documentation
  - `refactor:` - Refactoring
  - `test:` - Ajout de tests
  - `chore:` - Tâches de maintenance

### Tests

Tous les changements doivent inclure des tests:

```bash
# Lancer les tests
make test

# Les tests doivent passer avant de créer une PR
```

## Dépendances

### Principales

- `github.com/hashicorp/terraform-plugin-framework` v1.16.1 - Framework pour providers Terraform
- `github.com/hashicorp/terraform-plugin-docs` v0.24.0 - Génération de documentation

### Build

- Bazel 9.0.0rc1 - Système de build
- Go 1.24.0 - Langage de programmation

Voir `go.mod` pour la liste complète des dépendances.

## CI/CD

### GitHub Actions

- **`.github/workflows/release.yml`**: Publication automatique sur les tags `v*`

### Bazel

Le build Bazel assure:
- ✅ Builds reproductibles
- ✅ Cache distribué
- ✅ Compilation incrémentale
- ✅ Tests parallèles
- ✅ Support multi-plateforme

## Troubleshooting

### Bazel ne compile pas

```bash
# Nettoyer le cache
bazel clean --expunge

# Vérifier la version
bazel version  # Doit afficher 9.0.0rc1 ou supérieur

# Vérifier .bazelversion
cat .bazelversion
```

### Tests échouent

```bash
# Lancer les tests avec plus de détails
bazel test --test_output=all //src/...

# Vérifier les logs
bazel test --test_verbose_timeout_warnings //src/...
```

### DevContainer ne démarre pas

```bash
# Rebuild le container
CMD/CTRL + Shift + P > "Dev Containers: Rebuild Container"

# Vérifier les logs
CMD/CTRL + Shift + P > "Dev Containers: Show Log"
```

## License

MPL-2.0

---

**Développé avec ❤️ par [KodFlow](https://github.com/kodflow)**
