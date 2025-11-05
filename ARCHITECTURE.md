# 🏗️ Architecture Complète du Provider N8N Terraform

## 📁 Structure du Projet

```
src/internal/provider/
├── types/                          # Types partagés (pas de cycle d'import)
│   ├── client.go                   # N8nClient (wrapper SDK)
│   └── model.go                    # N8nProviderModel (config)
│
├── resources/                      # Resources Terraform
│   ├── resources.go                # Enregistrement de toutes les resources
│   ├── workflow.go                 # n8n_workflow
│   ├── credential.go               # n8n_credential
│   ├── tag.go                      # n8n_tag
│   ├── variable.go                 # n8n_variable
│   └── project.go                  # n8n_project
│
├── datasources/                    # Data Sources Terraform
│   ├── datasources.go              # Enregistrement de toutes les datasources
│   ├── workflow.go                 # data.n8n_workflow (single)
│   └── workflows.go                # data.n8n_workflows (list)
│
└── provider.go                     # Provider principal

```

## 🎯 Séparation des Responsabilités

### 1. **types/** - Types Partagés (Core)

**Rôle**: Définit les types de base utilisés par tous les autres packages.

#### `types/client.go`
```go
package types

type N8nClient struct {
    APIClient *n8nsdk.APIClient  // SDK généré
    BaseURL   string
    APIKey    string
}

func NewN8nClient(baseURL, apiKey string) *N8nClient
```

**Responsabilités**:
- Wrapper autour du SDK généré
- Configuration du client HTTP
- Ajout des headers d'authentification

#### `types/model.go`
```go
package types

type N8nProviderModel struct {
    APIKey  types.String `tfsdk:"api_key"`
    BaseURL types.String `tfsdk:"base_url"`
}
```

**Responsabilités**:
- Modèle de configuration du provider
- Définition du schema Terraform

**Pourquoi un package séparé ?**
- ✅ Évite les cycles d'imports
- ✅ Types réutilisables par resources/ et datasources/
- ✅ Séparation claire des responsabilités

---

### 2. **resources/** - Resources Terraform

**Rôle**: Implémente toutes les resources Terraform (CRUD).

#### `resources/resources.go` - Registry
```go
package resources

func Resources() []func() resource.Resource {
    return []func() resource.Resource{
        NewWorkflowResource,
        NewCredentialResource,
        NewTagResource,
        NewVariableResource,
        NewProjectResource,
    }
}
```

**Responsabilités**:
- **Point d'entrée unique** pour toutes les resources
- Enregistrement centralisé
- Facilite l'ajout de nouvelles resources

#### `resources/workflow.go` - Example Resource
```go
package resources

type WorkflowResource struct {
    client *providertypes.N8nClient  // Import depuis types/
}

func NewWorkflowResource() resource.Resource
func (r *WorkflowResource) Create(ctx, req, resp)  // Utilise SDK
func (r *WorkflowResource) Read(ctx, req, resp)
func (r *WorkflowResource) Update(ctx, req, resp)
func (r *WorkflowResource) Delete(ctx, req, resp)
```

**Pattern pour chaque resource**:
1. **Struct** avec client
2. **Factory function** `NewXxxResource()`
3. **Metadata** - Nom de la resource (`n8n_workflow`)
4. **Schema** - Attributs Terraform
5. **Configure** - Récupère le client du provider
6. **CRUD** - Utilise `client.APIClient.XxxAPI.MethodName()`

**Resources créées**:
- ✅ `workflow.go` → `n8n_workflow`
- ✅ `credential.go` → `n8n_credential`
- ✅ `tag.go` → `n8n_tag`
- ✅ `variable.go` → `n8n_variable`
- ✅ `project.go` → `n8n_project`

---

### 3. **datasources/** - Data Sources Terraform

**Rôle**: Implémente toutes les data sources (lecture seule).

#### `datasources/datasources.go` - Registry
```go
package datasources

func DataSources() []func() datasource.DataSource {
    return []func() datasource.DataSource{
        NewWorkflowDataSource,
        NewWorkflowsDataSource,
    }
}
```

**Responsabilités**:
- **Point d'entrée unique** pour toutes les datasources
- Enregistrement centralisé
- Facilite l'ajout de nouvelles datasources

#### `datasources/workflow.go` - Single Item
```go
package datasources

type WorkflowDataSource struct {
    client *providertypes.N8nClient
}

func NewWorkflowDataSource() datasource.DataSource
func (d *WorkflowDataSource) Read(ctx, req, resp)  // GET /workflows/{id}
```

**Usage Terraform**:
```hcl
data "n8n_workflow" "existing" {
  id = "workflow-123"
}
```

#### `datasources/workflows.go` - List
```go
package datasources

type WorkflowsDataSource struct {
    client *providertypes.N8nClient
}

func NewWorkflowsDataSource() datasource.DataSource
func (d *WorkflowsDataSource) Read(ctx, req, resp)  // GET /workflows
```

**Usage Terraform**:
```hcl
data "n8n_workflows" "all_active" {
  active = true
}
```

**Pattern datasource vs resource**:
- ✅ **Single** (`workflow`) → Récupère UN élément par ID
- ✅ **Plural** (`workflows`) → Liste avec filtres optionnels
- ✅ **Read-only** → Pas de Create/Update/Delete

---

### 4. **provider.go** - Provider Principal

**Rôle**: Orchestre tout le provider.

```go
package provider

import (
    "github.com/kodflow/n8n/src/internal/provider/datasources"
    "github.com/kodflow/n8n/src/internal/provider/resources"
    providertypes "github.com/kodflow/n8n/src/internal/provider/types"
)

type N8nProvider struct {
    version string
}

// Metadata - Nom et version du provider
func (p *N8nProvider) Metadata(ctx, req, resp) {
    resp.TypeName = "n8n"
    resp.Version = p.version
}

// Schema - Configuration provider (api_key, base_url)
func (p *N8nProvider) Schema(ctx, req, resp) { ... }

// Configure - Crée le client SDK et le partage
func (p *N8nProvider) Configure(ctx, req, resp) {
    config := &providertypes.N8nProviderModel{}

    client := providertypes.NewN8nClient(
        config.BaseURL.ValueString(),
        config.APIKey.ValueString(),
    )

    resp.ResourceData = client      // Partagé avec resources
    resp.DataSourceData = client    // Partagé avec datasources
}

// Resources - Délègue à resources.Resources()
func (p *N8nProvider) Resources(ctx) []func() resource.Resource {
    return resources.Resources()
}

// DataSources - Délègue à datasources.DataSources()
func (p *N8nProvider) DataSources(ctx) []func() datasource.DataSource {
    return datasources.DataSources()
}
```

**Responsabilités**:
- ✅ **Configuration globale** du provider
- ✅ **Création du client SDK** partagé
- ✅ **Enregistrement** resources et datasources
- ✅ **Point d'entrée** pour Terraform

---

## 🔄 Flux de Données

### 1. Initialisation du Provider

```
User HCL Config
      ↓
provider.Configure()
      ↓
providertypes.NewN8nClient(baseURL, apiKey)
      ↓
n8nsdk.NewConfiguration()
      ↓
n8nsdk.NewAPIClient(cfg)
      ↓
Client partagé → resp.ResourceData
                → resp.DataSourceData
```

### 2. Utilisation d'une Resource

```
terraform apply
      ↓
provider.Resources() → resources.Resources()
      ↓
resources.NewWorkflowResource()
      ↓
WorkflowResource.Configure(client)
      ↓
WorkflowResource.Create()
      ↓
client.APIClient.WorkflowAPI.WorkflowsPost(ctx)
      ↓
SDK HTTP Request → n8n API
```

### 3. Utilisation d'une DataSource

```
data "n8n_workflows" "all"
      ↓
provider.DataSources() → datasources.DataSources()
      ↓
datasources.NewWorkflowsDataSource()
      ↓
WorkflowsDataSource.Configure(client)
      ↓
WorkflowsDataSource.Read()
      ↓
client.APIClient.WorkflowAPI.WorkflowsGet(ctx)
      ↓
SDK HTTP Request → n8n API
```

---

## ✅ Avantages de cette Architecture

### 1. **Séparation des Responsabilités**
```
types/       → Types partagés (pas de logique)
resources/   → CRUD resources (écriture)
datasources/ → Lecture seule (queries)
provider.go  → Orchestration
```

### 2. **Pas de Cycles d'Import**
```
provider.go  → imports → types/, resources/, datasources/
resources/   → imports → types/  (PAS provider!)
datasources/ → imports → types/  (PAS provider!)
types/       → imports → RIEN du provider
```

### 3. **Scalabilité**
Pour ajouter une nouvelle resource:
```bash
# 1. Créer le fichier
src/internal/provider/resources/user.go

# 2. Ajouter à resources.go
func Resources() []func() resource.Resource {
    return []func() resource.Resource{
        NewWorkflowResource,
        NewUserResource,  // ← Nouvelle ligne !
    }
}
```

### 4. **Testabilité**
Chaque package peut être testé indépendamment:
```go
// resources/workflow_test.go
func TestWorkflowResource(t *testing.T) {
    // Mock providertypes.N8nClient
    // Test CRUD operations
}
```

### 5. **Maintenance**
- ✅ **1 fichier = 1 resource** → Facile à trouver
- ✅ **Registry centralisé** → Vue d'ensemble
- ✅ **Types partagés** → DRY (Don't Repeat Yourself)
- ✅ **SDK encapsulé** → Changements SDK isolés dans types/

---

## 📦 Fichiers Créés

### Types (2 fichiers)
```
types/client.go  → N8nClient wrapper
types/model.go   → N8nProviderModel config
```

### Resources (6 fichiers)
```
resources/resources.go   → Registry
resources/workflow.go    → n8n_workflow
resources/credential.go  → n8n_credential
resources/tag.go         → n8n_tag
resources/variable.go    → n8n_variable
resources/project.go     → n8n_project
```

### DataSources (3 fichiers)
```
datasources/datasources.go  → Registry
datasources/workflow.go      → data.n8n_workflow
datasources/workflows.go     → data.n8n_workflows
```

### Provider (1 fichier)
```
provider.go  → Orchestration
```

**Total**: **12 fichiers** organisés dans **4 packages**

---

## 🚀 Prochaines Étapes

### 1. Correction des Erreurs de Compilation
Les resources utilisent des méthodes SDK qui n'existent pas exactement:
- credential.go: `CredentialsIdGet` → Vérifier les méthodes réelles
- project.go: Signatures incorrectes
- variable.go: Retours de fonction incorrects

### 2. Ajout de Plus de Resources
```
resources/user.go        → n8n_user
resources/execution.go   → n8n_execution (read-only?)
```

### 3. Ajout de Plus de DataSources
```
datasources/credential.go   → data.n8n_credential
datasources/credentials.go  → data.n8n_credentials
datasources/tag.go          → data.n8n_tag
datasources/tags.go         → data.n8n_tags
```

### 4. Tests
```
resources/workflow_test.go
datasources/workflow_test.go
types/client_test.go
```

### 5. Documentation
```
docs/resources/workflow.md
docs/datasources/workflows.md
```

---

## 🎉 Résumé

**Architecture créée** ✅:
- ✅ Structure propre et scalable
- ✅ Séparation claire des responsabilités
- ✅ Pas de cycles d'import
- ✅ 1 fichier par resource/datasource
- ✅ Registries centralisés
- ✅ Types partagés isolés
- ✅ SDK encapsulé proprement

**À corriger** ⚠️:
- Erreurs de compilation (méthodes SDK)
- Vérifier les signatures de fonctions
- Ajuster les types de retour

**État actuel**: Architecture complète, erreurs de détail à corriger.
