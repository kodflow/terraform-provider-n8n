# Coverage Report

Rapport de couverture généré automatiquement.

**Légende:**
- 🟢 100% - Couverture complète
- 🟡 ≥90% - Bonne couverture
- 🔴 <90% - Couverture insuffisante
- ⚫ 0% - Pas de tests

---

## Coverage par Package

| Icon | Package | Coverage |
|:----:|---------|----------|
| 🔴 | **TOTAL (Global)** | **70.9%** |
| 🟢 | `github.com/kodflow/n8n/src/internal/provider` | 100.0% |
| 🟢 | `github.com/kodflow/n8n/src/internal/provider/shared/client` | 100.0% |
| 🟡 | `github.com/kodflow/n8n/src/internal/provider/workflow` | 83.0% |
| 🔴 | `github.com/kodflow/n8n/src/internal/provider/user` | 78.1% |
| 🔴 | `github.com/kodflow/n8n/src/internal/provider/execution` | 76.6% |
| 🔴 | `github.com/kodflow/n8n/src/internal/provider/variable` | 75.0% |
| 🔴 | `github.com/kodflow/n8n/src/internal/provider/tag` | 71.3% |
| 🔴 | `github.com/kodflow/n8n/src/internal/provider/credential` | 65.2% |
| 🔴 | `github.com/kodflow/n8n/src/internal/provider/sourcecontrol` | 60.0% |
| 🔴 | `github.com/kodflow/n8n/src/internal/provider/project` | 52.4% |

---

## Détail des packages incomplets

### 🟡 `github.com/kodflow/n8n/src/internal/provider/workflow` - 83.0%

**Forces:**
- Tests CRUD complets avec mocking HTTP
- Datasource Read tests implémentés
- Transfer resource tests ajoutés

**Zones à améliorer:**
| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟡 `datasource.go:Read` | 92.9% |
| 🟡 `transfer_resource.go:Create` | 92.9% |

### 🔴 `github.com/kodflow/n8n/src/internal/provider/user` - 78.1%

**Forces:**
- Tests CRUD HTTP-mocked
- Datasource Read tests
- Tests de validation email

**Zones à améliorer:**
| Fichier:Fonction | Couverture |
|------------------|------------|
| ⚫ `resource.go:NewUserResource` | 0.0% |
| ⚫ `resource.go:NewUserResourceWrapper` | 0.0% |
| ⚫ `resource.go:Metadata` | 0.0% |
| ⚫ `resource.go:Configure` | 0.0% |
| 🔴 `resource.go:Create` | 83.3% |
| 🔴 `datasource.go:Read` | 83.3% |
| ⚫ `datasources.go:Read` | 0.0% |

### 🔴 `github.com/kodflow/n8n/src/internal/provider/execution` - 76.6%

**Forces:**
- Retry resource tests complets
- Datasource tests implémentés

**Zones à améliorer:**
| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `retry_resource.go:Create` | 81.2% |
| 🔴 `datasource.go:Read` | 81.2% |
| ⚫ `datasources.go:Read` | 0.0% |

### 🔴 `github.com/kodflow/n8n/src/internal/provider/variable` - 75.0%

**Forces:**
- Tests CRUD complets
- Datasource by ID et by key tests

**Zones à améliorer:**
| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `resource.go:Create` | 81.8% |
| 🔴 `datasource.go:Read` | 81.8% |
| ⚫ `datasources.go:Read` | 0.0% |

### 🔴 `github.com/kodflow/n8n/src/internal/provider/tag` - 71.3%

**Forces:**
- Tests CRUD HTTP-mocked
- Datasource Read tests (by ID, by name)

**Zones à améliorer:**
| Fichier:Fonction | Couverture |
|------------------|------------|
| ⚫ `resource.go:NewTagResource` | 0.0% |
| ⚫ `resource.go:NewTagResourceWrapper` | 0.0% |
| ⚫ `resource.go:Metadata` | 0.0% |
| ⚫ `resource.go:Configure` | 0.0% |
| 🔴 `resource.go:Create` | 83.3% |
| 🔴 `datasource.go:Read` | 85.7% |
| ⚫ `datasources.go:Read` | 0.0% |

### 🔴 `github.com/kodflow/n8n/src/internal/provider/credential` - 65.2%

**Forces:**
- Tests CRUD de base
- Transfer resource tests
- Tests de rotation (partiel)

**Zones à améliorer:**
| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟡 `helpers.go:createNewCredential` | 75.0% |
| 🔴 `helpers.go:scanAffectedWorkflows` | 60.0% |
| ⚫ `helpers.go:updateAffectedWorkflows` | 12.5% |
| 🟡 `resource.go:Create` | 90.5% |
| 🟡 `transfer_resource.go:Create` | 92.9% |
| 🔴 `transfer_resource.go:Read` | 80.0% |

**Problèmes majeurs:**
- La fonction `updateAffectedWorkflows` (12.5%) nécessite des tests de rollback
- Scénarios de rotation incomplète non testés

### 🔴 `github.com/kodflow/n8n/src/internal/provider/sourcecontrol` - 60.0%

**État actuel:**
- Tests existants pour pull resource
- Pas de tests CRUD HTTP-mocked

**Zones à améliorer:**
| Fichier:Fonction | Couverture |
|------------------|------------|
| ⚫ `pull_resource.go:Create` | 0.0% |

### 🔴 `github.com/kodflow/n8n/src/internal/provider/project` - 52.4%

**Forces:**
- Tests CRUD HTTP-mocked de base
- Datasource Read tests

**Zones à améliorer:**
| Fichier:Fonction | Couverture |
|------------------|------------|
| ⚫ `resource.go:NewProjectResource` | 0.0% |
| ⚫ `resource.go:NewProjectResourceWrapper` | 0.0% |
| ⚫ `resource.go:Metadata` | 0.0% |
| ⚫ `resource.go:Configure` | 0.0% |
| 🔴 `resource.go:Create` | 81.8% |
| 🟡 `datasource.go:Read` | 76.0% |
| ⚫ `datasources.go:Read` | 0.0% |
| ⚫ `user_resource.go:NewProjectUserResource` | 0.0% |
| ⚫ `user_resource.go:NewProjectUserResourceWrapper` | 0.0% |
| ⚫ `user_resource.go:Metadata` | 0.0% |
| ⚫ `user_resource.go:Schema` | 0.0% |
| ⚫ `user_resource.go:Configure` | 0.0% |
| ⚫ `user_resource.go:Create` | 0.0% |

**Problèmes majeurs:**
- User resource non testée (0%)
- Datasources multiples non testés (0%)

---

## Améliorations apportées

**Progression globale:** 32.6% → 70.9% (+38.3 points)

### Tests ajoutés

1. **Infrastructure HTTP Mocking**
   - Helper `setupTestClient()` avec httptest
   - Configuration SDK pour tests
   - Gestion correcte des chemins API

2. **Tests CRUD Resources** (7 packages)
   - Project, Variable, Tag, User, Credential, Workflow, Execution
   - Scénarios succès et échec
   - Gestion des états Terraform

3. **Tests DataSources** (6 packages)
   - Tag, User, Variable, Project, Workflow, Execution
   - Tests par ID, par nom, not found

4. **Tests Transfer Resources**
   - Credential Transfer (71 tests)
   - Workflow Transfer (46 tests)

### Pour atteindre 90%

**Actions prioritaires:**

1. **Fonctions constructeurs** (0% actuellement)
   - Ajouter tests unitaires simples pour New* et Wrapper
   - Estimé: +5% de couverture

2. **Datasources multiples** (0% actuellement)
   - Implémenter tests pour datasources.go
   - Estimé: +3% de couverture

3. **Project User Resource** (0% actuellement)
   - Tests CRUD complets avec HTTP mocking
   - Estimé: +4% de couverture

4. **Source Control** (60% actuellement)
   - Tests Create/Read/Delete avec mocking
   - Estimé: +2% de couverture

5. **Credential rotation rollback** (12.5% actuellement)
   - Tests scénarios d'erreur complexes
   - Estimé: +3% de couverture

6. **Edge cases et error paths**
   - Nil checks, validations, error branches
   - Estimé: +3% de couverture

**Total estimé avec ces actions: ~90%**

---

*Rapport généré le: 2025-11-08*
