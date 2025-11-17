# N8N Nodes Synchronization System

## Overview

Ce système permet de synchroniser automatiquement tous les nodes n8n depuis le repository officiel, de détecter les changements, et de générer automatiquement des exemples Terraform.

## 🎯 Objectifs

- **Automatiser** la découverte de tous les nodes n8n
- **Synchroniser** avec chaque nouvelle version de n8n
- **Générer** automatiquement des exemples Terraform
- **Détecter** les ajouts, suppressions, et modifications de nodes
- **Maintenir** un référentiel complet et à jour

## 📊 Statistiques Actuelles

- **296 nodes** catalogués
- **3 catégories** :
  - Core: 5 nodes (Code, If, Merge, Set, Switch)
  - Trigger: 25 nodes (Webhook, Schedule, Email, etc.)
  - Integration: 266 nodes (GitHub, Slack, PostgreSQL, etc.)

## 🚀 Utilisation

### Synchronisation Complète

```bash
make nodes
```

Cette commande exécute :
1. Fetch du repository n8n officiel
2. Parsing de tous les nodes
3. Génération du changelog (si changements)
4. Génération des exemples Terraform
5. Affichage des statistiques

### Commandes Individuelles

```bash
# Récupérer le repository n8n
make nodes/fetch

# Parser les nodes
make nodes/parse

# Générer le changelog
make nodes/diff

# Afficher les statistiques
make nodes/stats

# Générer les exemples
make nodes/generate

# Nettoyer le cache
make nodes/clean
```

## 📁 Structure des Fichiers

```
/workspace/
├── data/
│   ├── n8n-nodes-registry.json    # Registre complet (296 nodes)
│   ├── n8n-nodes-metadata.json    # Métadonnées et stats
│   ├── n8n-nodes-version.txt      # Version n8n trackée
│   └── n8n-nodes-changelog.md     # Changelog auto-généré
├── scripts/nodes/
│   ├── sync-n8n-nodes.sh          # Script principal
│   ├── parse-nodes.js             # Parser TypeScript -> JSON
│   ├── generate-diff.js           # Génération changelog
│   └── generate-examples.js       # Génération exemples TF
├── examples/nodes/
│   ├── core/                      # 5 nodes Core
│   ├── trigger/                   # 25 trigger nodes
│   ├── integration/               # 266 intégrations
│   └── INDEX.md                   # Index complet
└── .n8n-repo-cache/               # Cache du repo (gitignored)
```

## 📋 Format du Registry JSON

```json
{
  "version": "v1.119.2",
  "last_sync": "2025-11-17T14:00:00Z",
  "total_nodes": 296,
  "nodes": [
    {
      "name": "Webhook",
      "type": "n8n-nodes-base.webhook",
      "category": "Trigger",
      "group": "trigger",
      "versions": [1, 2],
      "latest_version": 2,
      "description": "Wait for a webhook call",
      "inputs": [],
      "outputs": ["main"],
      "file": "packages/nodes-base/nodes/Webhook/Webhook.node.ts",
      "resources": {
        "primaryDocumentation": [...]
      }
    }
  ]
}
```

## 🔄 Workflow de Synchronisation

1. **Fetch** : Clone/update du repository n8n (shallow clone, branch master)
2. **Parse** : Parcours de `packages/nodes-base/nodes/`
3. **Extract** : Lecture des fichiers `.node.ts` et extraction des métadonnées
4. **Generate** : Création du registry JSON avec tous les nodes
5. **Diff** : Comparaison avec la version précédente
6. **Changelog** : Génération automatique du changelog
7. **Examples** : Génération d'exemples Terraform pour chaque catégorie

## 📝 Exemples Générés

Chaque catégorie de nodes a son dossier avec :

- `main.tf` - Exemples de nodes de la catégorie
- `variables.tf` - Variables Terraform
- `README.md` - Documentation complète listant tous les nodes

### Exemple : Core Nodes

```terraform
resource "n8n_workflow_node" "code" {
  name     = "Code"
  type     = "code"
  position = [250, 300]

  parameters = jsonencode({
    mode   = "runOnceForAllItems"
    jsCode = "return items;"
  })
}
```

## 🔍 Détection des Changements

Le système détecte automatiquement :

- ✅ **Nouveaux nodes** ajoutés
- ❌ **Nodes supprimés** (deprecated)
- 🔄 **Modifications** (version, description, inputs/outputs)

Le changelog est automatiquement généré dans `data/n8n-nodes-changelog.md`.

## 🎯 Use Cases

### 1. Mettre à jour après une nouvelle release n8n

```bash
make nodes
git add data/
git commit -m "chore(nodes): sync with n8n v1.120.0"
```

### 2. Vérifier si de nouveaux nodes sont disponibles

```bash
make nodes/fetch nodes/parse nodes/diff
cat data/n8n-nodes-changelog.md
```

### 3. Générer des exemples pour une catégorie spécifique

```bash
make nodes/parse
node scripts/nodes/generate-examples.js data/ examples/nodes/
```

## 🛠️ Développement

### Ajouter un nouveau type de parsing

Modifier `scripts/nodes/parse-nodes.js` :

```javascript
// Ajouter une nouvelle extraction
const customMatch = content.match(/customField:\s*['"]([^'"]+)['"]/);
nodeInfo.customField = customMatch ? customMatch[1] : null;
```

### Personnaliser la génération d'exemples

Modifier `scripts/nodes/generate-examples.js` pour ajuster le format Terraform généré.

## 📊 Métriques et Monitoring

Le système track :
- Nombre total de nodes
- Répartition par catégorie
- Nombre de versions par node
- Dernière synchronisation

Voir les stats avec :
```bash
make nodes/stats
```

## 🔐 Sécurité

- Le cache `.n8n-repo-cache/` est gitignored
- Pas de credentials stockés
- Clone shallow (historique minimal)
- Lecture seule du repository officiel

## 🚧 Limitations Connues

1. **Parsing TypeScript limité** : Utilise regex, pas un vrai parser TS
2. **Paramètres incomplets** : Les paramètres des nodes ne sont pas entièrement extraits
3. **Credentials non parsés** : Les credentials des nodes ne sont pas documentés

## 🔮 Futures Améliorations

- [ ] Parser TypeScript complet avec AST
- [ ] Extraction complète des paramètres de chaque node
- [ ] Génération de types Go constants
- [ ] Tests acceptance auto-générés
- [ ] CI/CD GitHub Actions pour auto-sync
- [ ] Documentation auto-générée pour chaque node

## 📚 Ressources

- [Repository n8n officiel](https://github.com/n8n-io/n8n)
- [Documentation n8n](https://docs.n8n.io)
- [NODES_SYNC_PROGRESS.md](/workspace/NODES_SYNC_PROGRESS.md) - Suivi détaillé

## 🤝 Contribution

Pour ajouter de nouvelles fonctionnalités au système de sync :

1. Modifier les scripts dans `scripts/nodes/`
2. Tester avec `make nodes`
3. Mettre à jour cette documentation
4. Soumettre une PR

---

**Dernière mise à jour** : 17 Novembre 2025
**Version n8n trackée** : v1.119.2
**Nodes catalogués** : 296
