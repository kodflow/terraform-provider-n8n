# 🔧 Analyse des Workarounds pour n8n_credential

## 🎯 Objectif

Implémenter `resource "n8n_credential"` malgré les limitations de l'API n8n:
- ❌ Pas de GET /credentials/{id}
- ❌ Pas de GET /credentials (list)
- ❌ Pas de PUT /credentials/{id}

## 🔧 Workaround Proposés

### Workaround 1: Update = DELETE + POST

**Principe**: Au lieu de PUT, faire DELETE puis POST

**Implémentation**:
```go
func (r *CredentialResource) Update(ctx, req, resp) {
    // 1. DELETE l'ancien credential
    _, err := r.client.APIClient.CredentialAPI.
        DeleteCredential(ctx, state.ID.ValueString()).
        Execute()

    // 2. POST un nouveau credential
    newCred, _, err := r.client.APIClient.CredentialAPI.
        CredentialsPost(ctx).
        Credential(request).
        Execute()

    // 3. ⚠️ NOUVEL ID!
    plan.ID = types.StringPointerValue(newCred.Id)
}
```

### Workaround 2: Read = State Only (pas d'API)

**Principe**: Ne pas appeler l'API pour Read, utiliser le tfstate local

**Implémentation**:
```go
func (r *CredentialResource) Read(ctx, req, resp) {
    var state CredentialResourceModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

    // ⚠️ Pas d'API call - on assume que rien n'a changé

    resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
```

---

## ⚖️ Analyse des Risques

### Workaround 1: DELETE + POST

#### ✅ Avantages

1. **Fonctionne techniquement** - API supporte DELETE et POST
2. **Pas de dépendance à des endpoints manquants**
3. **Terraform gère avec `create_before_destroy`**

#### ❌ Inconvénients Majeurs

1. **L'ID Change** 🚨
   ```hcl
   # Avant update
   credential_id = "cred-abc123"

   # Après update
   credential_id = "cred-xyz789"  # ⚠️ DIFFÉRENT!
   ```

   **Impact**: Toutes les références doivent être mises à jour:
   ```hcl
   resource "n8n_workflow" "example" {
     # Ces références sont CASSÉES après update du credential
     nodes = [{
       credentials = {
         id = n8n_credential.api.id  # ❌ Ancien ID, plus valide
       }
     }]
   }
   ```

2. **Downtime** ⏱️
   - Entre DELETE et POST: credential n'existe pas
   - Workflows qui l'utilisent échouent pendant ce temps
   - Même avec `create_before_destroy`, il y a un moment de transition

3. **Rollback Impossible** 💥
   ```go
   // 1. DELETE réussit ✅
   DeleteCredential(ctx, oldID).Execute()

   // 2. POST échoue ❌ (erreur réseau, validation, etc.)
   newCred, err := CredentialsPost(ctx).Execute()
   // ❌ On a perdu le credential! Pas de rollback possible
   ```

4. **Perte de l'Historique** 📜
   - Nouvelle création = nouveau timestamp
   - Audit trail cassé
   - Plus moyen de tracer l'historique du credential

5. **Problème avec les Workflows Actifs** 🔴
   ```
   Workflow actif avec credential X (id: old-123)
   → Update credential
   → Credential recréé (id: new-456)
   → Workflow référence toujours old-123
   → Workflow CASSE!
   ```

#### 🔧 Mitigation Partielle

**Lifecycle Policy**:
```hcl
resource "n8n_credential" "api" {
  name = "API Key"
  type = "httpHeaderAuth"

  lifecycle {
    # Crée le nouveau AVANT de supprimer l'ancien
    create_before_destroy = true

    # Avertir l'utilisateur
    # (commentaire, pas une vraie option Terraform)
  }
}
```

**Limites de la mitigation**:
- ✅ Réduit le downtime (nouveau existe avant suppression ancien)
- ❌ Ne résout PAS le problème de l'ID changé
- ❌ Les références doivent TOUJOURS être mises à jour manuellement

---

### Workaround 2: State-Only Read

#### ✅ Avantages

1. **Simple à implémenter**
2. **Pas d'API call** - pas de limitations
3. **Terraform fonctionne "normalement"**

#### ❌ Inconvénients Majeurs

1. **State Drift Non Détecté** 🚨
   ```bash
   # Scénario:
   # 1. Créer credential via Terraform
   $ terraform apply

   # 2. Quelqu'un supprime via UI n8n
   # (Le credential n'existe plus dans n8n)

   # 3. Terraform ne le détecte JAMAIS
   $ terraform plan
   No changes. Infrastructure is up-to-date.
   # ❌ FAUX! Le credential n'existe plus!

   $ terraform refresh
   # Ne fait rien car Read() ne vérifie pas l'API

   # 4. Essayer d'utiliser le credential
   # ❌ ÉCHEC - il n'existe pas!
   ```

2. **Faux Sentiment de Sécurité** 🎭
   - Le tfstate dit "credential existe"
   - La réalité: peut-être supprimé, modifié, ou corrompu
   - Utilisateur pense que son infra est synchro

3. **Import Impossible** 🔗
   ```bash
   $ terraform import n8n_credential.api cred-123
   # Comment vérifier que cred-123 existe réellement?
   # Read() ne vérifie pas l'API
   # ❌ On importe peut-être un credential qui n'existe pas!
   ```

4. **Refresh Inutile** 🔄
   ```bash
   $ terraform refresh
   # Censé synchroniser avec l'infra réelle
   # Mais Read() ne fait rien
   # = Commande inutile
   ```

5. **Violations Best Practices** 📋
   - **Principle of Truth**: Terraform doit refléter la réalité
   - **Idempotence**: Plusieurs apply = même résultat
   - **Declarative**: State = source de vérité

   Avec state-only:
   - ❌ State peut être faux
   - ❌ Apply peut échouer silencieusement
   - ❌ Source de vérité = mensonge

6. **Debugging Impossible** 🐛
   ```bash
   # Utilisateur: "Mon workflow ne marche plus!"
   # Support: "Vérifions le credential"

   $ terraform state show n8n_credential.api
   # Montre les données du state
   # ❌ Mais ça ne prouve PAS que le credential existe dans n8n!

   # Impossible de dire si le problème vient de:
   # - Credential supprimé dans n8n
   # - Credential modifié dans n8n
   # - State désynchronisé
   # - Autre chose
   ```

#### 🔧 Mitigation Partielle

**Documentation Claire**:
```hcl
resource "n8n_credential" "api" {
  name = "API Key"

  # ⚠️ WARNING: This resource cannot detect drift!
  # If the credential is modified or deleted in n8n UI,
  # Terraform will not detect it.
  #
  # Manual verification required:
  # 1. Check n8n UI regularly
  # 2. Test workflows using this credential
  # 3. Consider using n8n audit logs
}
```

**Workaround pour Import**:
```go
func (r *CredentialResource) ImportState(ctx, req, resp) {
    // ⚠️ On ne peut PAS vérifier que le credential existe
    // Accepter l'ID et espérer que c'est bon

    resp.Diagnostics.AddWarning(
        "Cannot verify credential existence",
        "The n8n API does not support reading credentials. "+
        "This import assumes the credential ID is valid. "+
        "If the credential does not exist, operations will fail.",
    )

    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
```

**Limites de la mitigation**:
- ✅ Utilisateur est averti
- ❌ Problème reste présent
- ❌ Pas de vraie solution

---

## 🎯 Recommandation Finale

### Option A: Implémenter avec Workarounds ⚠️

**Approche**: DELETE+POST pour Update, State-only pour Read

**Utilisation**:
```hcl
resource "n8n_credential" "api" {
  name = "API Key"
  type = "httpHeaderAuth"
  data = {
    name  = "Authorization"
    value = var.api_token
  }

  lifecycle {
    create_before_destroy = true

    # ⚠️ WARNINGS OBLIGATOIRES
    # 1. Update = destroy + recreate (ID changes)
    # 2. Drift detection not supported
    # 3. Manual verification required
  }
}

# ⚠️ Après un update, mettre à jour les références:
resource "n8n_workflow" "example" {
  # Doit être mis à jour manuellement si credential change
  nodes = [{
    credentials = {
      id = n8n_credential.api.id
    }
  }]
}
```

**Documentation Requise**:
```markdown
# ⚠️ Limitations Importantes

## n8n_credential

Cette resource a des limitations dues à l'API n8n:

### Update = Destroy + Recreate
- Modifier un credential le SUPPRIME puis RECRÉE
- L'ID change à chaque update
- Toutes les références doivent être mises à jour manuellement

### Drift Detection Non Supportée
- Terraform ne peut pas détecter si le credential est modifié/supprimé dans n8n
- `terraform refresh` ne synchronise pas
- Vérification manuelle requise

### Recommandation
- Éviter d'updater les credentials (utiliser versionning)
- Vérifier régulièrement dans l'UI n8n
- Préférer créer de nouveaux credentials plutôt qu'updater

### Alternative
Créer les credentials manuellement dans n8n UI et utiliser
data source pour les référencer (quand disponible).
```

**Risques**:
- 🔴 **ÉLEVÉ**: Références cassées après update
- 🔴 **ÉLEVÉ**: State drift non détecté
- 🟡 **MOYEN**: Confusion utilisateur
- 🟡 **MOYEN**: Support complexe

**Verdict**: ⚠️ **Faisable mais RISQUÉ**

---

### Option B: NE PAS Implémenter ✅ RECOMMANDÉ

**Approche**: Provider sans `n8n_credential`, documentation claire

**Utilisation**:
```hcl
# ✅ Resources disponibles
resource "n8n_workflow" "example" { }
resource "n8n_tag" "example" { }
resource "n8n_variable" "example" { }
resource "n8n_project" "example" { }

# ❌ Credentials - créer manuellement dans l'UI n8n
# Puis référencer par ID si nécessaire
```

**Documentation**:
```markdown
# Limitations

## Credentials Non Supportés

La resource `n8n_credential` n'est pas disponible.

### Raison
L'API publique n8n ne fournit pas les endpoints nécessaires:
- ❌ GET /api/v1/credentials/{id} (lecture)
- ❌ GET /api/v1/credentials (liste)
- ⏳ PUT /api/v1/credentials/{id} (update - PR ouvert)

### Workarounds Possibles Mais Non Recommandés

Nous avons étudié des workarounds:

1. **Update via DELETE+POST**
   - ❌ Change l'ID du credential
   - ❌ Casse toutes les références
   - ❌ Downtime entre suppression et création

2. **Read via tfstate uniquement**
   - ❌ Drift detection impossible
   - ❌ Faux sentiment de sécurité
   - ❌ Violations best practices Terraform

Ces workarounds créent plus de problèmes qu'ils n'en résolvent.

### Solution Recommandée

Créer les credentials manuellement dans l'UI n8n:
1. Ouvrir l'interface web n8n
2. Aller dans Credentials
3. Créer le credential
4. Noter l'ID pour référence dans les workflows Terraform

### Évolution Future

Si n8n ajoute les endpoints manquants à l'API publique,
nous ajouterons la resource `n8n_credential` au provider.

Suivre: https://github.com/n8n-io/n8n/pull/18082
```

**Avantages**:
- ✅ **Honnête** sur les capacités
- ✅ **Pas de fausses promesses**
- ✅ **Pas de comportements surprenants**
- ✅ **Support simplifié**
- ✅ **4/5 resources fonctionnelles**

**Verdict**: ✅ **RECOMMANDÉ**

---

## 📊 Comparaison

| Critère | Option A (Workarounds) | Option B (Sans Credential) |
|---------|------------------------|----------------------------|
| **Fonctionnalité** | ⚠️ Partielle (CRUD cassé) | ✅ 4/5 resources complètes |
| **Honnêteté** | ⚠️ Cache les problèmes | ✅ Documentation claire |
| **UX** | 🔴 Surprises négatives | ✅ Prévisible |
| **Maintenance** | 🔴 Complexe | ✅ Simple |
| **Support** | 🔴 Difficile | ✅ Facile |
| **Risques** | 🔴 Élevés | ✅ Aucun |
| **Best Practices** | ❌ Violations | ✅ Respectées |

---

## 🎯 Décision

### ⭐ Recommandation: **Option B**

**Pourquoi**:

1. **Intégrité**: Ne pas promettre ce qui ne fonctionne pas correctement
2. **Prévisibilité**: Pas de comportements surprenants
3. **Documentation**: Expliquer clairement les limitations (API n8n, pas provider)
4. **Évolutivité**: Quand n8n ajoute les endpoints, on ajoutera la resource
5. **Utilisateur**: 4 resources fonctionnelles > 5 resources dont 1 cassée

**Message aux utilisateurs**:
```
Le provider n8n Terraform gère 4 des 5 resources principales.

La resource n8n_credential n'est pas disponible car l'API publique n8n
ne fournit pas les endpoints nécessaires (GET, PUT).

Cette limitation vient de n8n, pas du provider. Les credentials doivent
être créés manuellement via l'interface web n8n.

Nous avons exploré des workarounds (DELETE+POST pour update, state-only
pour read) mais ils créent plus de problèmes qu'ils n'en résolvent:
- Update change l'ID (casse les références)
- Drift detection impossible
- Violations des best practices Terraform

Nous préférons être honnêtes sur les limitations plutôt que de livrer
une resource cassée.
```

---

## 📚 Annexes

### Code Workaround 1: DELETE+POST

```go
func (r *CredentialResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var plan, state CredentialResourceModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // WORKAROUND: DELETE + POST
    resp.Diagnostics.AddWarning(
        "Update will change credential ID",
        "Updating a credential deletes and recreates it with a new ID. "+
        "All references must be updated manually.",
    )

    // 1. DELETE old credential
    _, httpResp, err := r.client.APIClient.CredentialAPI.
        DeleteCredential(ctx, state.ID.ValueString()).
        Execute()

    if err != nil {
        resp.Diagnostics.AddError(
            "Failed to delete old credential",
            fmt.Sprintf("Could not delete credential %s: %s", state.ID.ValueString(), err.Error()),
        )
        return
    }

    // 2. POST new credential
    credRequest := n8nsdk.Credential{
        Name: plan.Name.ValueString(),
        Type: plan.Type.ValueString(),
        Data: plan.Data.ValueMap(), // Assuming proper conversion
    }

    newCred, httpResp, err := r.client.APIClient.CredentialAPI.
        CredentialsPost(ctx).
        Credential(credRequest).
        Execute()

    if err != nil {
        // ❌ PROBLÈME: Old credential deleted, new creation failed
        // Credential is LOST!
        resp.Diagnostics.AddError(
            "Failed to create new credential after delete",
            fmt.Sprintf("Old credential was deleted but new creation failed: %s", err.Error()),
        )
        // State is now inconsistent
        return
    }

    // 3. ⚠️ NEW ID
    plan.ID = types.StringPointerValue(newCred.Id)

    resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
```

### Code Workaround 2: State-Only Read

```go
func (r *CredentialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state CredentialResourceModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // WORKAROUND: No API call, use state as-is
    // ⚠️ Cannot detect if credential was modified or deleted in n8n

    resp.Diagnostics.AddWarning(
        "Drift detection not supported",
        "The n8n API does not support reading credentials. "+
        "Terraform cannot detect if the credential was modified or deleted in n8n. "+
        "Manual verification required.",
    )

    // Just keep the state as-is
    resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
```

### Import avec Warning

```go
func (r *CredentialResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    // Cannot verify credential exists
    resp.Diagnostics.AddWarning(
        "Cannot verify credential existence",
        "The n8n API does not support reading credentials. "+
        "This import assumes the credential ID '"+req.ID+"' is valid. "+
        "If the credential does not exist or has wrong type, operations will fail.",
    )

    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
```
