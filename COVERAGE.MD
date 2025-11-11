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
| **Total Coverage** | **76.4%** |
| **Threshold** | **70.0%** |
| **Status** | ✅ PASSED |

---

## Coverage par Package

| Icon | Package | Coverage |
|:----:|---------|----------|
| 🟢 | `github.com/kodflow/n8n/src/internal/provider` | 100.0% |
| 🔴 | `github.com/kodflow/n8n/src/internal/provider/credential` | 66.5% |
| 🟡 | `github.com/kodflow/n8n/src/internal/provider/execution` | 79.2% |
| 🔴 | `github.com/kodflow/n8n/src/internal/provider/project` | 67.2% |
| 🟢 | `github.com/kodflow/n8n/src/internal/provider/shared` | 100.0% |
| 🟢 | `github.com/kodflow/n8n/src/internal/provider/shared/client` | 100.0% |
| 🟢 | `github.com/kodflow/n8n/src/internal/provider/sourcecontrol` | 97.0% |
| 🟡 | `github.com/kodflow/n8n/src/internal/provider/tag` | 79.9% |
| 🟡 | `github.com/kodflow/n8n/src/internal/provider/user` | 73.5% |
| 🟡 | `github.com/kodflow/n8n/src/internal/provider/variable` | 74.0% |
| 🟡 | `github.com/kodflow/n8n/src/internal/provider/workflow` | 87.4% |

---

## Coverage Détaillée par Fichier

Cette section liste uniquement les **fonctions publiques** (exportées) pour identifier rapidement les fonctions non testées.
Les tableaux sont organisés par type de fichier pour faciliter la comparaison entre packages.

### 📄 datasource.go

| Function | execution | project | provider | tag | user | variable | workflow |
|----------|:--------:|:--------:|:--------:|:--------:|:--------:|:--------:|:--------:|
| `Configure` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |
| `Metadata` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |
| `Read` | 🔴 0.0% | 🔴 0.0% | 🔴 0.0% | 🔴 0.0% | 🔴 0.0% | 🔴 0.0% | 🔴 0.0% |
| `Schema` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |

### 📄 datasources.go

| Function | execution | project | provider | tag | user | variable | workflow |
|----------|:--------:|:--------:|:--------:|:--------:|:--------:|:--------:|:--------:|
| `Configure` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |
| `Metadata` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |
| `Read` | 🔴 0.0% | 🔴 0.0% | 🔴 0.0% | 🔴 0.0% | 🔴 0.0% | 🔴 0.0% | 🔴 0.0% |
| `Schema` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |

### 📄 pointers.go

| Function | provider | shared |
|----------|:--------:|:--------:|
| `Bool` | 🟢 100.0% | 🟢 100.0% |
| `Float32` | 🟢 100.0% | 🟢 100.0% |
| `Int` | 🟢 100.0% | 🟢 100.0% |
| `Int32` | 🟢 100.0% | 🟢 100.0% |
| `Ptr` | 🟢 100.0% | 🟢 100.0% |
| `String` | 🟢 100.0% | 🟢 100.0% |

### 📄 provider.go

| Function | provider |
|----------|:--------:|
| `Configure` | 🟢 100.0% |
| `DataSources` | 🟢 100.0% |
| `Metadata` | 🟢 100.0% |
| `Resources` | 🟢 100.0% |
| `Schema` | 🟢 100.0% |
| `ValidateProvider` | 🟢 100.0% |

### 📄 pull_resource.go

| Function | provider | sourcecontrol |
|----------|:--------:|:--------:|
| `Configure` | 🟢 100.0% | 🟢 100.0% |
| `Create` | 🟢 100.0% | 🟢 100.0% |
| `Delete` | 🟢 100.0% | 🟢 100.0% |
| `ImportState` | 🟢 100.0% | 🟢 100.0% |
| `Metadata` | 🟢 100.0% | 🟢 100.0% |
| `Read` | 🟢 100.0% | 🟢 100.0% |
| `Schema` | 🟢 100.0% | 🟢 100.0% |
| `Update` | 🟢 100.0% | 🟢 100.0% |

### 📄 resource.go

| Function | credential | project | provider | tag | user | variable | workflow |
|----------|:--------:|:--------:|:--------:|:--------:|:--------:|:--------:|:--------:|
| `Configure` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |
| `Create` | 🔴 0.0% | 🟡 81.8% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🔴 0.0% | 🟢 100.0% |
| `Delete` | 🟡 81.8% | 🟡 80.0% | 🟢 100.0% | 🟢 100.0% | 🔴 0.0% | 🟡 80.0% | 🟢 100.0% |
| `ImportState` | 🟢 100.0% | 🟢 100.0% | 🔴 0.0% | 🟢 100.0% | 🔴 0.0% | 🟢 100.0% | 🔴 0.0% |
| `Metadata` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |
| `Read` | 🟢 100.0% | 🟡 88.9% | 🟢 100.0% | 🟢 100.0% | 🔴 0.0% | 🔴 0.0% | 🟢 100.0% |
| `Schema` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |
| `Update` | 🔴 0.0% | 🟡 81.8% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🔴 0.0% | 🟢 100.0% |

### 📄 retry_resource.go

| Function | execution | provider |
|----------|:--------:|:--------:|
| `Configure` | 🟢 100.0% | 🟢 100.0% |
| `Create` | 🟢 100.0% | 🟢 100.0% |
| `Delete` | 🔴 0.0% | 🔴 0.0% |
| `ImportState` | 🔴 0.0% | 🔴 0.0% |
| `Metadata` | 🟢 100.0% | 🟢 100.0% |
| `Read` | 🟡 80.0% | 🟡 80.0% |
| `Schema` | 🟢 100.0% | 🟢 100.0% |
| `Update` | 🔴 0.0% | 🔴 0.0% |

### 📄 transfer_resource.go

| Function | credential | provider | workflow |
|----------|:--------:|:--------:|:--------:|
| `Configure` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |
| `Create` | 🔴 0.0% | 🟢 100.0% | 🟢 100.0% |
| `Delete` | 🔴 0.0% | 🔴 0.0% | 🔴 0.0% |
| `ImportState` | 🔴 0.0% | 🟢 100.0% | 🟢 100.0% |
| `Metadata` | 🔴 0.0% | 🟢 100.0% | 🟢 100.0% |
| `Read` | 🔴 0.0% | 🟢 100.0% | 🟢 100.0% |
| `Schema` | 🔴 0.0% | 🟢 100.0% | 🟢 100.0% |
| `Update` | 🔴 0.0% | 🟢 100.0% | 🟢 100.0% |

### 📄 user_resource.go

| Function | project | provider |
|----------|:--------:|:--------:|
| `Configure` | 🟢 100.0% | 🟢 100.0% |
| `Create` | 🔴 0.0% | 🔴 0.0% |
| `Delete` | 🔴 0.0% | 🔴 0.0% |
| `ImportState` | 🔴 0.0% | 🔴 0.0% |
| `Metadata` | 🟢 100.0% | 🟢 100.0% |
| `Read` | 🔴 0.0% | 🔴 0.0% |
| `Schema` | 🟢 100.0% | 🟢 100.0% |
| `Update` | 🔴 0.0% | 🔴 0.0% |


---

## Légende des Icônes

- 🟢 **≥90%** - Excellente couverture
- 🟡 **70-89%** - Bonne couverture
- 🟠 **1-69%** - Couverture partielle (à améliorer)
- 🔴 **0%** - Fonction non testée

**Note:** Seules les fonctions publiques (exportées) sont listées. Les fonctions privées et constructeurs (`New*`) sont exclus.

---

*Rapport généré le: 2025-11-11*
*Threshold: 70.0%*
