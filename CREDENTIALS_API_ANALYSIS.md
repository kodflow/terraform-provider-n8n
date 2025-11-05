# 🔍 Analyse Complète de l'API Credentials n8n

## 📋 Résumé Exécutif

L'API publique n8n actuelle (v1.1.1) a des **limitations importantes** pour les credentials:

| Endpoint | Méthode | Status | Notes |
|----------|---------|--------|-------|
| `/credentials` | POST | ✅ Disponible | Crée un credential |
| `/credentials/{id}` | GET | ❌ **NON DISPONIBLE** | Pas d'endpoint pour récupérer |
| `/credentials` | GET (LIST) | ❌ **NON DISPONIBLE** | Pas d'endpoint pour lister |
| `/credentials/{id}` | PUT | ⏳ **EN ATTENTE** | PR #18082 ouvert (pas mergé) |
| `/credentials/{id}` | DELETE | ✅ Disponible | Supprime un credential |
| `/credentials/{id}/transfer` | POST | ✅ Disponible | Transfère vers un projet |
| `/credential-types/{type}` | GET | ✅ Disponible | Récupère le schéma |

**Conclusion**: Impossible de faire CRUD complet car GET et UPDATE ne sont pas disponibles

---

## 🔎 Recherche Effectuée

### 1. Documentation Officielle

- **URL**: https://docs.n8n.io/api/api-reference/
- **Swagger UI**: Disponible à `/api/v1/api-docs` sur instances self-hosted
- **Résultat**: Documentation confirmée limitée à POST/DELETE

### 2. Code Source GitHub

- **Repository**: https://github.com/n8n-io/n8n
- **Fichiers Analysés**:
  - `packages/cli/src/public-api/v1/handlers/credentials/credentials.handler.ts`
  - `packages/cli/src/public-api/v1/handlers/credentials/credentials.service.ts`
  - `packages/cli/src/public-api/v1/handlers/credentials/spec/paths/credentials.yml`
  - `packages/cli/src/public-api/v1/handlers/credentials/spec/paths/credentials.id.yml`

### 3. Pull Requests

- **PR #18082**: "feat (public-api): update credentials"
  - **Status**: ⏳ **OUVERT** (pas encore mergé)
  - **Auteur**: Shock3udt
  - **Date**: 2024
  - **Contenu**: Ajoute PUT `/credentials/{id}` pour update
  - **URL**: https://github.com/n8n-io/n8n/pull/18082

---

## ✅ Endpoints Disponibles (API Actuelle v1.1.1)

### 1. POST `/credentials` - Create Credential

**Source**: `packages/cli/src/public-api/v1/handlers/credentials/spec/paths/credentials.yml`

**Request**:
```http
POST /api/v1/credentials
Content-Type: application/json
X-N8N-API-KEY: your-api-key

{
  "name": "My Credential",
  "type": "httpHeaderAuth",
  "data": {
    "name": "Authorization",
    "value": "Bearer token123"
  }
}
```

**Response 200**:
```json
{
  "id": "credential-uuid",
  "name": "My Credential",
  "type": "httpHeaderAuth",
  "createdAt": "2024-01-01T00:00:00.000Z",
  "updatedAt": "2024-01-01T00:00:00.000Z"
}
```

**Permissions**: Nécessite scope `credential:create`

**Code Handler**:
```typescript
// packages/cli/src/public-api/v1/handlers/credentials/credentials.handler.ts
export const createCredential = [
  validCredentialType,
  validCredentialsProperties,
  apiKeyHasScope('credential:create'),
  async (req: CredentialRequest.Create, res: express.Response): Promise<express.Response> => {
    // ... implementation
  }
]
```

---

### 2. DELETE `/credentials/{id}` - Delete Credential

**Source**: `packages/cli/src/public-api/v1/handlers/credentials/spec/paths/credentials.id.yml`

**Request**:
```http
DELETE /api/v1/credentials/credential-uuid
X-N8N-API-KEY: your-api-key
```

**Response 200**:
```json
{
  "id": "credential-uuid",
  "name": "My Credential",
  "type": "httpHeaderAuth"
}
```

**Permissions**:
- Scope `credential:delete`
- Doit être owner ou admin du credential

**Code Handler**:
```typescript
export const deleteCredential = [
  apiKeyHasScope('credential:delete'),
  projectScope('credential:delete', 'credential'),
  async (req: CredentialRequest.Delete, res: express.Response): Promise<express.Response> => {
    const { id: credentialId } = req.params;
    // ... implementation
  }
]
```

---

### 3. POST `/credentials/{id}/transfer` - Transfer Credential

**Source**: `packages/cli/src/public-api/v1/handlers/credentials/spec/paths/credentials.id.transfer.yml`

**Request**:
```http
POST /api/v1/credentials/credential-uuid/transfer
Content-Type: application/json
X-N8N-API-KEY: your-api-key

{
  "destinationProjectId": "project-uuid"
}
```

**Response**: `204 No Content`

**Permissions**: Scope `credential:move`

---

### 4. GET `/credential-types/{type}` - Get Credential Schema

**Source**: `packages/cli/src/public-api/v1/handlers/credentials/spec/paths/credentials.schema.id.yml`

**Request**:
```http
GET /api/v1/credential-types/httpHeaderAuth
X-N8N-API-KEY: your-api-key
```

**Response 200**:
```json
{
  "type": "object",
  "properties": {
    "name": { "type": "string" },
    "value": { "type": "string" }
  },
  "required": ["name", "value"]
}
```

---

## ❌ Endpoints Manquants

### 1. GET `/credentials` - List Credentials

**Status**: ❌ **N'EXISTE PAS**

**Raison**: Sécurité - Les credentials contiennent des données sensibles

**Impact**: Impossible de lister tous les credentials d'un utilisateur

**Workaround**: Aucun (pas d'alternative)

---

### 2. GET `/credentials/{id}` - Get Credential by ID

**Status**: ❌ **N'EXISTE PAS**

**Raison**: Sécurité - Impossible de récupérer les données d'un credential

**Impact**:
- Impossible de faire Read() dans Terraform
- Impossible de vérifier l'existence d'un credential
- Impossible de rafraîchir le state

**Workaround**: Aucun (pas d'alternative)

**Code Service Existant** (interne seulement):
```typescript
// packages/cli/src/public-api/v1/handlers/credentials/credentials.service.ts
async getCredentials(credentialId: string): Promise<ICredentialsDb | null> {
  return await this.credentialsRepository.findOneBy({ id: credentialId });
}
```
☝️ Cette méthode existe mais n'est **pas exposée** via l'API publique

---

### 3. PUT `/credentials/{id}` - Update Credential

**Status**: ⏳ **EN ATTENTE** (PR #18082 ouvert)

**PR**: https://github.com/n8n-io/n8n/pull/18082

**Proposition dans le PR**:
```http
PUT /api/v1/credentials/credential-uuid
Content-Type: application/json
X-N8N-API-KEY: your-api-key

{
  "name": "Updated Credential Name",
  "data": {
    "name": "Authorization",
    "value": "Bearer new-token"
  }
}
```

**Response Proposée**:
```json
{
  "id": "credential-uuid",
  "name": "Updated Credential Name",
  "type": "httpHeaderAuth",
  "updatedAt": "2024-01-02T00:00:00.000Z"
}
```

**Permissions Proposées**: Scope `credential:update`

**Code Service Proposé**:
```typescript
async updateCredential(
  credentialId: string,
  properties: CredentialRequest.CredentialProperties,
): Promise<CredentialsEntity> {
  const credential = await this.getCredentials(credentialId);
  if (!credential) {
    throw new NotFoundError('Credential not found');
  }

  credential.name = properties.name ?? credential.name;
  if (properties.data) {
    credential.data = await this.encryptCredential(properties.data);
  }

  return await this.credentialsRepository.save(credential);
}
```

**Problèmes Identifiés dans le PR**:
- ❌ Validation d'entrée manquante dans le middleware
- ❌ Type `as any` pour contourner la sécurité des types
- ❌ Tests mocké incorrectement

**Statut**: Le PR est suivi en interne (ticket GHC-3571) mais **pas de timeline**

---

## 🔒 Raisons de Sécurité

### Pourquoi pas de GET?

n8n a **volontairement** limité l'API publique pour les credentials:

1. **Données sensibles**: Les credentials contiennent des tokens, API keys, passwords
2. **Principe du moindre privilège**: Pas besoin de lire pour automatiser
3. **Rotation des credentials**: L'UI web utilise des endpoints **internes** (pas publics)
4. **Architecture**:
   - UI web → Endpoints internes (avec plus de permissions)
   - API publique → Endpoints limités (automation seulement)

### Endpoints Internes vs Publics

```
ENDPOINTS INTERNES (UI web):
/rest/credentials           GET   ✅ Liste
/rest/credentials/:id       GET   ✅ Lecture
/rest/credentials/:id       PATCH ✅ Update

ENDPOINTS PUBLICS (API):
/api/v1/credentials         POST  ✅ Création
/api/v1/credentials/:id     DELETE ✅ Suppression
/api/v1/credentials/:id     GET   ❌ Pas disponible
/api/v1/credentials/:id     PUT   ⏳ PR ouvert
```

---

## 🎯 Impact sur le Provider Terraform

### Limitations Actuelles

Avec l'API actuelle, **impossible** d'implémenter une resource `n8n_credential` complète:

```hcl
resource "n8n_credential" "api_key" {
  name = "My API Key"
  type = "httpHeaderAuth"
  data = {
    name  = "Authorization"
    value = var.api_token
  }
}
```

**Problèmes**:

1. ❌ **Create()** - ✅ Fonctionne (POST disponible)
2. ❌ **Read()** - ❌ **IMPOSSIBLE** (pas de GET)
3. ❌ **Update()** - ❌ **IMPOSSIBLE** (PUT pas mergé)
4. ❌ **Delete()** - ✅ Fonctionne (DELETE disponible)
5. ❌ **Import** - ❌ **IMPOSSIBLE** (pas de GET pour vérifier)

**Conclusion**: Resource `n8n_credential` **NON VIABLE** avec l'API actuelle

---

## 💡 Solutions Possibles

### Option 1: Attendre le Merge du PR #18082 ⏳

**Avantages**:
- ✅ Solution officielle et supportée
- ✅ PUT disponible pour Update()

**Inconvénients**:
- ❌ Timeline inconnue (aucune promesse de n8n)
- ❌ Ne résout pas le problème de Read() (pas de GET)
- ❌ Toujours pas de LIST

**Recommandation**: **NON** - Trop incertain

---

### Option 2: Patch OpenAPI pour Ajouter GET + PUT ⭐

**Approche**: Créer un patch qui ajoute les endpoints manquants

**Patch Proposé**:
```yaml
# credentials-api.patch
# Ajoute GET /credentials (LIST)
/credentials:
  get:
    operationId: credentialsGet
    summary: List all credentials
    responses:
      '200':
        content:
          application/json:
            schema:
              type: array
              items:
                $ref: '#/components/schemas/credential'

# Ajoute GET /credentials/{id}
/credentials/{id}:
  get:
    operationId: credentialsIdGet
    summary: Get credential by ID
    parameters:
      - name: id
        in: path
        required: true
        schema:
          type: string
    responses:
      '200':
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/credential'

  # Ajoute PUT /credentials/{id}
  put:
    operationId: credentialsIdPut
    summary: Update credential
    parameters:
      - name: id
        in: path
        required: true
        schema:
          type: string
    requestBody:
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/credential'
    responses:
      '200':
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/credential'
```

**Avantages**:
- ✅ SDK généré avec méthodes complètes
- ✅ Provider Terraform fonctionnel
- ✅ Réutilisable après chaque téléchargement d'OpenAPI

**Inconvénients**:
- ⚠️ **FAUX** - Les endpoints n'existent PAS réellement dans l'API
- ❌ Le provider **NE FONCTIONNERA PAS** en prod
- ❌ Mentir sur ce que l'API supporte

**Recommandation**: **NON** - C'est mentir

---

### Option 3: Implémenter SANS Credentials ⭐⭐⭐

**Approche**: Ne pas implémenter `n8n_credential`, documenter pourquoi

**Provider Final**:
```
✅ n8n_workflow    (CRUD complet)
✅ n8n_tag         (CRUD complet)
✅ n8n_variable    (CRUD avec workarounds LIST)
✅ n8n_project     (CRUD avec workarounds LIST)
❌ n8n_credential  (API incomplète - non implémentable)
```

**Documentation**:
```markdown
# Limitations

La resource `n8n_credential` n'est pas disponible car l'API publique n8n
ne supporte pas les opérations de lecture (GET) nécessaires pour Terraform.

Endpoints manquants:
- GET /api/v1/credentials (list)
- GET /api/v1/credentials/{id} (read)

Ces endpoints existent dans l'UI web mais ne sont pas exposés dans l'API publique
pour des raisons de sécurité.

## Alternatives

1. Créer les credentials manuellement dans l'UI n8n
2. Utiliser l'API interne (non documentée, peut changer)
3. Attendre que n8n expose ces endpoints dans l'API publique
```

**Avantages**:
- ✅ Honnête sur les capacités réelles
- ✅ Provider fonctionne pour 4 resources sur 5
- ✅ Pas de fausses promesses

**Inconvénients**:
- ⚠️ Credentials non gérés par Terraform

**Recommandation**: **OUI** ⭐⭐⭐ - La seule approche honnête

---

### Option 4: Utiliser l'API Interne (Non Documentée) ⚠️

**Approche**: Utiliser `/rest/credentials` au lieu de `/api/v1/credentials`

**Endpoints Internes**:
```http
GET    /rest/credentials          # Liste
GET    /rest/credentials/:id      # Read
PATCH  /rest/credentials/:id      # Update
DELETE /rest/credentials/:id      # Delete
POST   /rest/credentials          # Create
```

**Avantages**:
- ✅ CRUD complet disponible
- ✅ C'est ce que l'UI utilise

**Inconvénients**:
- ❌ **Non documenté** - Peut changer sans préavis
- ❌ **Non supporté** officiellement
- ❌ Authentification différente (sessions vs API key)
- ❌ Risque de breaking changes
- ❌ Pas éthique (API interne = privée)

**Recommandation**: **NON** - Trop risqué et non éthique

---

## 📊 Comparaison des Options

| Option | Honnêteté | Faisabilité | Maintenabilité | Recommandation |
|--------|-----------|-------------|----------------|----------------|
| **1. Attendre PR** | ✅ | ⏳ Timeline inconnue | ✅ | ⚠️ Trop incertain |
| **2. Patch OpenAPI** | ❌ Mensonge | ✅ Facile | ❌ Provider cassé | ❌ NON |
| **3. Sans Credentials** | ✅ | ✅ | ✅ | ✅ **RECOMMANDÉ** |
| **4. API Interne** | ⚠️ Grey area | ✅ | ❌ Risque | ❌ NON |

---

## 🎯 Recommandation Finale

### ⭐ **Option 3: Provider SANS `n8n_credential`**

**Justification**:

1. **Honnêteté**: Ne pas promettre ce qui ne fonctionne pas
2. **Utilisabilité**: 4/5 resources fonctionnelles, c'est déjà excellent
3. **Documentation**: Expliquer clairement pourquoi Credential manque
4. **Évolutif**: Quand n8n ajoute les endpoints, on pourra ajouter la resource

**Provider Final**:
```hcl
provider "n8n" {
  api_key  = var.n8n_api_key
  base_url = var.n8n_base_url
}

# ✅ Resources Disponibles
resource "n8n_workflow" "example" { }
resource "n8n_tag" "example" { }
resource "n8n_variable" "example" { }
resource "n8n_project" "example" { }

# ❌ Credentials - gérer manuellement dans l'UI
# Raison: API publique n8n ne supporte pas GET/UPDATE
```

**Impact Utilisateur**:

- ✅ Workflows automatisés avec Terraform
- ✅ Tags gérés comme IaC
- ✅ Variables d'environnement versionnées
- ✅ Projects organisés
- ⚠️ Credentials créés manuellement (limitation API, pas provider)

**Message aux utilisateurs**:
```
Le provider n8n Terraform gère 4 des 5 resources principales.

La resource `n8n_credential` n'est pas disponible car l'API publique
n8n ne fournit pas les endpoints nécessaires pour la lecture et mise à jour.

Cette limitation vient de n8n, pas du provider. Les credentials doivent
être gérés manuellement via l'interface web n8n.

Référence: https://github.com/n8n-io/n8n/pull/18082
```

---

## 📚 Références

- **API Documentation**: https://docs.n8n.io/api/api-reference/
- **GitHub Repository**: https://github.com/n8n-io/n8n
- **PR Update Credentials**: https://github.com/n8n-io/n8n/pull/18082
- **Feature Request**: https://community.n8n.io/t/get-update-credentials-via-api/46437
- **OpenAPI Spec**: https://github.com/n8n-io/n8n/blob/master/packages/cli/src/public-api/v1/openApiSpec.ts

---

## ✅ Prochaines Étapes

1. ✅ **Documenter** cette analyse (fait)
2. **Décider** quelle option implémenter
3. Si Option 3:
   - Créer README.md avec limitations
   - Documenter pourquoi Credential manque
   - Fournir workaround (création manuelle)
4. **Compiler** le provider final
5. **Tester** avec une instance n8n réelle
6. **Publier** avec documentation claire

**Décision**: Attendre validation utilisateur pour choisir l'option
