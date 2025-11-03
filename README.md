# DevContainer Minimal Template

Template minimaliste pour démarrer rapidement vos projets avec un environnement DevContainer propre et léger.

## Fonctionnalités

- **Ubuntu 24.04 LTS** comme base
- **User vscode** (UID/GID 1000:1000) avec sudo
- **Zsh + Oh My Zsh + Powerlevel10k** pré-installé et configuré
- **Outils essentiels** : git, curl, wget, jq, yq, build-essential
- **MCP (Model Context Protocol)** : Configuration et scripts d'initialisation inclus
- **Script d'initialisation** : Configuration automatique à la création du container
- **Persistance** via volumes Docker
- **Aucune feature externe** : tout est dans le Dockerfile
- **GitHub Actions** : Workflow de build automatisé pour le devcontainer

## Ce qui n'est PAS inclus

Ce template est **volontairement minimaliste**. Il ne contient pas :

- ❌ Langages de programmation (Go, Node.js, Python, etc.)
- ❌ CLIs spécifiques (GitHub CLI, Claude CLI, etc.)
- ❌ Docker-in-Docker
- ❌ Bases de données

**Pourquoi ?** Pour garder l'image légère et vous laisser installer uniquement ce dont vous avez besoin.

## Installation rapide

### Via GitHub

```bash
# Utiliser ce repository comme template
gh repo create mon-projet --template .repository --public
cd mon-projet
code .
```

### Localement

```bash
# Copier le template
cp -r .repository mon-projet
cd mon-projet
rm -rf .git
git init
code .
```

Acceptez l'ouverture dans le DevContainer lorsque VS Code vous le propose.

## Personnalisation

### Ajouter des langages/outils

**Option 1 : Dans le Dockerfile** (recommandé pour les outils systèmes)

Éditez `.devcontainer/Dockerfile` :

```dockerfile
# Ajouter des packages apt
RUN apt-get update && apt-get install -y \
    python3 \
    python3-pip \
    && apt-get clean

# Installer Node.js
RUN curl -fsSL https://deb.nodesource.com/setup_lts.x | bash - \
    && apt-get install -y nodejs
```

**Option 2 : Avec les DevContainer Features** (pour les langages standards)

Ajoutez dans `.devcontainer/devcontainer.json` :

```json
"features": {
  "ghcr.io/devcontainers/features/go:1": {
    "version": "latest"
  },
  "ghcr.io/devcontainers/features/node:1": {
    "version": "lts"
  }
}
```

Voir : <https://containers.dev/features>

**Option 3 : Installation manuelle** (pour les outils utilisateur)

Installez après l'ouverture du container :

```bash
# Exemple
curl -sSL https://example.com/install.sh | sh
```

### Ajouter des extensions VS Code

Éditez `.devcontainer/devcontainer.json` dans `customizations.vscode.extensions`.

### Variables d'environnement

Créez un fichier `.env` à la racine pour vos variables d'environnement.

### Personnaliser Powerlevel10k

Pour configurer le prompt Powerlevel10k :

```bash
# Lancer l'assistant de configuration interactif
p10k configure
```

Cela créera un fichier `~/.p10k.zsh` avec votre configuration personnalisée. Ce fichier sera automatiquement chargé au démarrage du shell.

## Structure des volumes

Les volumes Docker persistent entre les rebuilds :

### Volumes spécifiques au projet

- `{nom-du-projet}-local-bin` : Binaires locaux installés

### Volumes partagés

- `vscode-extensions` : Extensions VS Code
- `vscode-insiders-extensions` : Extensions VS Code Insiders
- `zsh-history` : Historique Zsh

Vous pouvez ajouter vos propres volumes dans `.devcontainer/devcontainer.json`.

## Commandes utiles

### Rebuild du container

```bash
# Depuis VS Code
Cmd+Shift+P > "Dev Containers: Rebuild Container"

# Ou depuis le terminal
docker compose -f .devcontainer/docker-compose.yml down
docker compose -f .devcontainer/docker-compose.yml build --no-cache
docker compose -f .devcontainer/docker-compose.yml up -d
```

### Nettoyer les volumes

```bash
# Supprimer tous les volumes (⚠️ perte de données)
docker compose -f .devcontainer/docker-compose.yml down -v
```

### Voir les logs

```bash
docker compose -f .devcontainer/docker-compose.yml logs -f devcontainer
```

## Configuration MCP (Model Context Protocol)

Le template inclut une configuration MCP pour faciliter l'intégration avec des outils d'IA.

### Script de configuration

Le script `.devcontainer/setup-mcp.sh` est exécuté automatiquement au démarrage du container et permet de :

- Configurer les serveurs MCP
- Initialiser les variables d'environnement nécessaires
- Préparer l'environnement pour l'utilisation des outils MCP

### Variables d'environnement

Copiez `.devcontainer/.env.example` vers `.devcontainer/.env` et configurez vos variables :

```bash
cp .devcontainer/.env.example .devcontainer/.env
```

Éditez `.devcontainer/.env` selon vos besoins pour ajouter vos clés API et configurations.

### Script d'initialisation

Le script `.devcontainer/init.sh` s'exécute automatiquement à la création du container via le `postCreateCommand` et permet de :

- Effectuer des configurations post-création
- Installer des dépendances supplémentaires
- Personnaliser l'environnement de développement

## Dépannage

### Problèmes de permissions

```bash
# Depuis le container
sudo chown -R vscode:vscode $HOME
```

### Rebuild complet

```bash
# Supprimer le container et les volumes
docker compose -f .devcontainer/docker-compose.yml down -v
docker system prune -a

# Rouvrir dans VS Code
code .
```

## Philosophie

Ce template suit le principe **"moins c'est plus"** :

- ✅ Démarrage rapide
- ✅ Faible consommation de ressources
- ✅ Facile à personnaliser
- ✅ Pas de dépendances inutiles

Ajoutez seulement ce dont vous avez besoin, quand vous en avez besoin.

## License

Libre d'utilisation pour vos projets personnels et professionnels.
