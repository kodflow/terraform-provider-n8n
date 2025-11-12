# Coverage Report

Rapport de couverture généré automatiquement.

**Légende:**
- 🟢 ≥90% - Excellente couverture
- 🟡 70-89% - Bonne couverture
- 🔴 <70% - Couverture insuffisante

---

## Coverage Global

| Metric | Value |
|--------|-------|
| **Total Coverage** | **83.8%** |
| **Threshold** | **70.0%** |
| **Status** | ✅ PASSED |

---

## Coverage par Package

| Icon | Package | Coverage |
|:----:|---------|----------|
| 🟢 | `github.com/kodflow/n8n/src/internal/provider` | 100.0% |
| 🟡 | `github.com/kodflow/n8n/src/internal/provider/credential` | 77.9% |
| 🟡 | `github.com/kodflow/n8n/src/internal/provider/execution` | 81.2% |
| 🟡 | `github.com/kodflow/n8n/src/internal/provider/project` | 77.8% |
| 🟢 | `github.com/kodflow/n8n/src/internal/provider/shared` | 100.0% |
| 🟢 | `github.com/kodflow/n8n/src/internal/provider/shared/client` | 100.0% |
| 🟢 | `github.com/kodflow/n8n/src/internal/provider/sourcecontrol` | 97.0% |
| 🟡 | `github.com/kodflow/n8n/src/internal/provider/tag` | 81.1% |
| 🟡 | `github.com/kodflow/n8n/src/internal/provider/user` | 85.0% |
| 🟡 | `github.com/kodflow/n8n/src/internal/provider/variable` | 88.7% |
| 🟡 | `github.com/kodflow/n8n/src/internal/provider/workflow` | 88.2% |

---

## Coverage Détaillée par Catégorie

Cette section liste uniquement les **fonctions publiques** (exportées) pour identifier rapidement les fonctions non testées.
Les tableaux sont organisés par catégorie de ressources pour faciliter la compréhension de l'architecture du provider.

## 📦 Ressources Principales (CRUD Entities)

Gestion complète du cycle de vie des ressources n8n (Create, Read, Update, Delete).

| Function | credential | project | tag | user | variable | workflow |
|----------|:--------:|:--------:|:--------:|:--------:|:--------:|:--------:|
| `Configure` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |
| `Create` | 🟡 85.7% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |
| `Delete` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |
| `ImportState` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |
| `Metadata` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |
| `Read` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |
| `Schema` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |
| `Update` | 🟡 87.5% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |

---

## 🔧 Ressources Secondaires (Operations/Relations)

Opérations spéciales et gestion des relations entre ressources.

### Source Control Pull

| Function | sourcecontrol |
|----------|:--------:|
| `Configure` | 🟢 100.0% |
| `Create` | 🟢 100.0% |
| `Delete` | 🟢 100.0% |
| `ImportState` | 🟢 100.0% |
| `Metadata` | 🟢 100.0% |
| `Read` | 🟢 100.0% |
| `Schema` | 🟢 100.0% |
| `Update` | 🟢 100.0% |

### Execution Retry

| Function | execution |
|----------|:--------:|
| `Configure` | 🟢 100.0% |
| `Create` | 🟢 100.0% |
| `Delete` | 🔵 N/A |
| `ImportState` | 🟢 100.0% |
| `Metadata` | 🟢 100.0% |
| `Read` | 🟢 100.0% |
| `Schema` | 🟢 100.0% |
| `Update` | 🟢 100.0% |

### Transfer

| Function | credential | workflow |
|----------|:--------:|:--------:|
| `Configure` | 🟢 100.0% | 🟢 100.0% |
| `Create` | 🔴 0.0% | 🟢 100.0% |
| `Delete` | 🔵 N/A | 🔵 N/A |
| `ImportState` | 🔴 0.0% | 🟢 100.0% |
| `Metadata` | 🔴 0.0% | 🟢 100.0% |
| `Read` | 🔴 0.0% | 🟢 100.0% |
| `Schema` | 🔴 0.0% | 🟢 100.0% |
| `Update` | 🔴 0.0% | 🟢 100.0% |

### Project User

| Function | project |
|----------|:--------:|
| `Configure` | 🟢 100.0% |
| `Create` | 🔴 0.0% |
| `Delete` | 🔴 0.0% |
| `ImportState` | 🔴 0.0% |
| `Metadata` | 🟢 100.0% |
| `Read` | 🔴 0.0% |
| `Schema` | 🟢 100.0% |
| `Update` | 🔴 0.0% |

---

## 📊 Sources de Données (Data Sources)

Lecture des ressources n8n sans gestion de leur cycle de vie.

### datasource (singular)

| Function | execution | project | tag | user | variable | workflow |
|----------|:--------:|:--------:|:--------:|:--------:|:--------:|:--------:|
| `Configure` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |
| `Metadata` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |
| `Read` | 🔴 0.0% | 🔴 0.0% | 🔴 0.0% | 🔴 0.0% | 🔴 0.0% | 🔴 0.0% |
| `Schema` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |

### datasources (plural)

| Function | execution | project | tag | user | variable | workflow |
|----------|:--------:|:--------:|:--------:|:--------:|:--------:|:--------:|
| `Configure` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |
| `Metadata` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |
| `Read` | 🔴 0.0% | 🔴 0.0% | 🔴 0.0% | 🔴 0.0% | 🔴 0.0% | 🔴 0.0% |
| `Schema` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |


---

## Légende des Icônes

- 🟢 **≥90%** - Excellente couverture
- 🟡 **70-89%** - Bonne couverture
- 🟠 **1-69%** - Couverture partielle (à améliorer)
- 🔴 **0%** - Fonction non testée (implémentée mais pas de tests)
- 🔵 **N/A** - Fonction non applicable (pas dans l'API ou intentionnellement vide)

**Note:** Seules les fonctions publiques (exportées) sont listées. Les fonctions privées et constructeurs (`New*`) sont exclus.

---

*Rapport généré le: 2025-11-12*
*Threshold: 70.0%*
