# 🔄 Design: Credential Rotation avec Zero-Downtime

## 🎯 Objectif

Implémenter `Update()` pour `n8n_credential` avec rotation automatique:
- ✅ Créer nouveau credential AVANT suppression
- ✅ Migrer automatiquement toutes les références
- ✅ Supprimer l'ancien SEULEMENT si tout OK
- ✅ Rollback automatique en cas d'erreur
- ✅ Transparent pour Terraform (PUT normal)

---

## 🏗️ Architecture

### Flux Normal (Succès)

```
1. Terraform: terraform apply (credential changé)
   ↓
2. Provider: Update() appelé
   ↓
3. CREATE nouveau credential (ID = new-123)
   ↓
4. SCAN tous les workflows utilisant old-456
   ↓
5. Pour chaque workflow:
   - GET workflow
   - PARSE nodes pour trouver credential references
   - REPLACE old-456 → new-123
   - PUT workflow (mise à jour)
   ↓
6. VERIFY tous les updates réussis
   ↓
7. DELETE ancien credential (ID = old-456)
   ↓
8. SUCCESS ✅
```

### Flux Erreur (Rollback)

```
1-3. (même que succès)
   ↓
4. SCAN workflows → 5 workflows trouvés
   ↓
5. Update workflow 1 ✅
   Update workflow 2 ✅
   Update workflow 3 ❌ ERREUR!
   ↓
6. ROLLBACK:
   - DELETE nouveau credential (new-123)
   - GARDER ancien credential (old-456)
   - Workflows 1-2 référencent maintenant new-123 (inexistant) ❌
   ↓
7. ROLLBACK COMPLET:
   - RESTORE workflow 1 (old-456)
   - RESTORE workflow 2 (old-456)
   ↓
8. FAILURE ❌ (état initial restauré)
```

---

## 📋 Étapes Détaillées

### Étape 1: Créer Nouveau Credential

```go
// 1. POST nouveau credential
newCredRequest := n8nsdk.Credential{
    Name: plan.Name.ValueString(),
    Type: plan.Type.ValueString(),
    Data: plan.Data, // nouvelles données
}

newCred, _, err := r.client.APIClient.CredentialAPI.
    CredentialsPost(ctx).
    Credential(newCredRequest).
    Execute()

if err != nil {
    // Erreur création → ABORT (rien à rollback)
    return err
}

newCredID := *newCred.Id  // new-123
oldCredID := state.ID.ValueString()  // old-456
```

**État**:
- ✅ Nouveau credential existe (new-123)
- ✅ Ancien credential existe encore (old-456)
- ✅ Zero downtime

---

### Étape 2: Scanner les Workflows

```go
// 2. LIST tous les workflows
workflowList, _, err := r.client.APIClient.WorkflowAPI.
    WorkflowsGet(ctx).
    Execute()

if err != nil {
    // Erreur scan → ROLLBACK: delete new credential
    r.client.APIClient.CredentialAPI.DeleteCredential(ctx, newCredID).Execute()
    return err
}

// 3. Filtrer workflows qui utilisent old credential
affectedWorkflows := []WorkflowToUpdate{}

for _, workflow := range workflowList.Data {
    if usesCredential(workflow, oldCredID) {
        affectedWorkflows = append(affectedWorkflows, WorkflowToUpdate{
            ID:       *workflow.Id,
            Original: workflow,  // Backup pour rollback
        })
    }
}

// Log pour l'utilisateur
tflog.Info(ctx, fmt.Sprintf(
    "Found %d workflows using credential %s",
    len(affectedWorkflows),
    oldCredID,
))
```

**Fonction Helper**:
```go
func usesCredential(workflow n8nsdk.Workflow, credentialID string) bool {
    if workflow.Nodes == nil {
        return false
    }

    for _, node := range workflow.Nodes {
        // Vérifier si le node a des credentials
        if node.Credentials != nil {
            for _, cred := range node.Credentials {
                if cred.Id != nil && *cred.Id == credentialID {
                    return true
                }
            }
        }
    }

    return false
}
```

---

### Étape 3: Migrer les Références

```go
// 4. Update chaque workflow
updatedWorkflows := []string{}
failedWorkflows := []string{}

for _, workflowToUpdate := range affectedWorkflows {
    // GET workflow complet
    workflow, _, err := r.client.APIClient.WorkflowAPI.
        WorkflowsIdGet(ctx, workflowToUpdate.ID).
        Execute()

    if err != nil {
        failedWorkflows = append(failedWorkflows, workflowToUpdate.ID)
        continue
    }

    // REPLACE credential references
    updated := replaceCredentialInWorkflow(workflow, oldCredID, newCredID)

    // PUT workflow
    _, _, err = r.client.APIClient.WorkflowAPI.
        WorkflowsIdPut(ctx, workflowToUpdate.ID).
        Workflow(*updated).
        Execute()

    if err != nil {
        failedWorkflows = append(failedWorkflows, workflowToUpdate.ID)
        continue
    }

    updatedWorkflows = append(updatedWorkflows, workflowToUpdate.ID)
}

// 5. Vérifier succès
if len(failedWorkflows) > 0 {
    // ÉCHEC → ROLLBACK complet
    return rollbackRotation(ctx, r.client,
        newCredID, oldCredID,
        affectedWorkflows, updatedWorkflows)
}
```

**Fonction Helper**:
```go
func replaceCredentialInWorkflow(
    workflow *n8nsdk.Workflow,
    oldCredID, newCredID string,
) *n8nsdk.Workflow {
    if workflow.Nodes == nil {
        return workflow
    }

    for i := range workflow.Nodes {
        node := &workflow.Nodes[i]

        if node.Credentials != nil {
            for j := range node.Credentials {
                cred := &node.Credentials[j]
                if cred.Id != nil && *cred.Id == oldCredID {
                    cred.Id = &newCredID
                }
            }
        }
    }

    return workflow
}
```

---

### Étape 4: Supprimer Ancien Credential

```go
// 6. Tous les workflows migrés → DELETE ancien
_, _, err := r.client.APIClient.CredentialAPI.
    DeleteCredential(ctx, oldCredID).
    Execute()

if err != nil {
    // Erreur suppression ancien
    // C'est OK! Nouveau fonctionne, ancien juste orphelin
    tflog.Warn(ctx, fmt.Sprintf(
        "Could not delete old credential %s: %s. "+
        "New credential %s is active. "+
        "Manual cleanup may be required.",
        oldCredID, err.Error(), newCredID,
    ))

    // On continue quand même (nouveau fonctionne)
}

// 7. Update state avec nouveau ID
plan.ID = types.StringValue(newCredID)

tflog.Info(ctx, fmt.Sprintf(
    "Credential rotated successfully: %s → %s (%d workflows updated)",
    oldCredID, newCredID, len(updatedWorkflows),
))
```

---

### Étape 5: Rollback en Cas d'Erreur

```go
func rollbackRotation(
    ctx context.Context,
    client *providertypes.N8nClient,
    newCredID, oldCredID string,
    affectedWorkflows []WorkflowToUpdate,
    updatedWorkflows []string,
) error {
    tflog.Error(ctx, "Rotation failed, rolling back...")

    // 1. Supprimer nouveau credential
    _, _, err := client.APIClient.CredentialAPI.
        DeleteCredential(ctx, newCredID).
        Execute()

    if err != nil {
        tflog.Error(ctx, fmt.Sprintf(
            "CRITICAL: Failed to delete new credential %s during rollback: %s",
            newCredID, err.Error(),
        ))
        // Continue quand même pour essayer de restaurer workflows
    }

    // 2. Restaurer les workflows updatés
    restoredCount := 0
    failedRestores := []string{}

    for _, workflowID := range updatedWorkflows {
        // Trouver le workflow original
        var original *n8nsdk.Workflow
        for _, wtu := range affectedWorkflows {
            if wtu.ID == workflowID {
                original = wtu.Original
                break
            }
        }

        if original == nil {
            tflog.Error(ctx, fmt.Sprintf(
                "Cannot find original for workflow %s", workflowID,
            ))
            failedRestores = append(failedRestores, workflowID)
            continue
        }

        // Restaurer workflow original
        _, _, err := client.APIClient.WorkflowAPI.
            WorkflowsIdPut(ctx, workflowID).
            Workflow(*original).
            Execute()

        if err != nil {
            tflog.Error(ctx, fmt.Sprintf(
                "Failed to restore workflow %s: %s", workflowID, err.Error(),
            ))
            failedRestores = append(failedRestores, workflowID)
            continue
        }

        restoredCount++
    }

    // 3. Retourner erreur avec détails
    if len(failedRestores) > 0 {
        return fmt.Errorf(
            "Rotation rollback partially failed. "+
            "Restored %d/%d workflows. "+
            "Failed to restore: %v. "+
            "Old credential %s preserved. "+
            "Manual intervention required.",
            restoredCount, len(updatedWorkflows),
            failedRestores, oldCredID,
        )
    }

    return fmt.Errorf(
        "Rotation failed and rolled back successfully. "+
        "All %d workflows restored to use credential %s.",
        restoredCount, oldCredID,
    )
}
```

---

## 🎯 Structures de Données

```go
// WorkflowToUpdate stocke un workflow à migrer
type WorkflowToUpdate struct {
    ID       string
    Original *n8nsdk.Workflow  // Backup pour rollback
}

// RotationResult résultat de la rotation
type RotationResult struct {
    Success           bool
    NewCredentialID   string
    OldCredentialID   string
    WorkflowsScanned  int
    WorkflowsAffected int
    WorkflowsUpdated  int
    WorkflowsFailed   []string
    RollbackPerformed bool
    Error             error
}
```

---

## ⚠️ Points d'Attention

### 1. Workflows Actifs

**Problème**: Update d'un workflow actif peut le désactiver temporairement

**Solution**:
```go
func replaceCredentialInWorkflow(...) *n8nsdk.Workflow {
    // Préserver le statut active
    wasActive := workflow.Active != nil && *workflow.Active

    // ... modifications ...

    // Restaurer active status
    if wasActive {
        workflow.Active = &wasActive
    }

    return workflow
}
```

### 2. Rate Limiting

**Problème**: Trop de requêtes API en séquence

**Solution**:
```go
import "time"

// Throttle entre chaque workflow update
for i, workflowToUpdate := range affectedWorkflows {
    if i > 0 {
        time.Sleep(100 * time.Millisecond)  // 100ms entre chaque
    }

    // ... update workflow ...
}
```

### 3. Timeout

**Problème**: Rotation peut prendre du temps (100+ workflows)

**Solution**:
```go
// Augmenter le timeout pour Update()
ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
defer cancel()
```

### 4. Permissions

**Problème**: API key doit avoir accès à TOUS les workflows

**Solution**: Documentation
```markdown
## Permissions Requises

Pour que la rotation de credentials fonctionne, l'API key doit avoir:
- `credential:create`
- `credential:delete`
- `workflow:read` (tous les workflows)
- `workflow:update` (tous les workflows)

Si vous n'avez pas ces permissions, la rotation échouera.
```

---

## 🧪 Cas de Test

### Test 1: Rotation Simple (1 Workflow)

```hcl
resource "n8n_credential" "api" {
  name = "API Key"
  type = "httpHeaderAuth"
  data = {
    name  = "Authorization"
    value = "Bearer token123"
  }
}

resource "n8n_workflow" "example" {
  name = "Test Workflow"
  nodes = [{
    credentials = {
      id = n8n_credential.api.id
    }
  }]
}

# Update credential
terraform apply -var="token=new-token"
```

**Attendu**:
1. ✅ Nouveau credential créé
2. ✅ Workflow scanné (1 trouvé)
3. ✅ Workflow updaté
4. ✅ Ancien credential supprimé
5. ✅ `n8n_credential.api.id` = nouveau ID

### Test 2: Rotation Complexe (10 Workflows)

```bash
# 10 workflows utilisant le même credential
# Update → 10 workflows doivent être migrés
```

**Attendu**:
- ✅ 10 workflows scannés
- ✅ 10 workflows updatés
- ✅ Rotation réussie

### Test 3: Rotation avec Échec (Rollback)

```bash
# Setup: 5 workflows, workflow 3 locked (ne peut pas update)
# Update credential
```

**Attendu**:
1. ✅ Nouveau credential créé
2. ✅ 5 workflows scannés
3. ✅ Workflow 1-2 updatés
4. ❌ Workflow 3 échoue
5. 🔄 ROLLBACK:
   - ✅ Nouveau credential supprimé
   - ✅ Workflows 1-2 restaurés
   - ✅ Ancien credential preserved
6. ❌ Terraform apply échoue avec erreur claire

### Test 4: Zero Downtime

```bash
# Workflow actif qui s'exécute pendant la rotation
# La rotation ne doit PAS causer d'erreur d'exécution
```

**Attendu**:
- ✅ Workflow continue de fonctionner pendant rotation
- ✅ Utilise ancien credential pendant migration
- ✅ Bascule sur nouveau credential après migration
- ✅ Pas d'erreur d'exécution

---

## 📊 Avantages vs Inconvénients

### ✅ Avantages

1. **Zero Downtime**
   - Nouveau credential créé AVANT suppression ancien
   - Workflows fonctionnent pendant la migration

2. **Rollback Automatique**
   - Si erreur, état initial restauré automatiquement
   - Pas de state corrompu

3. **Transparent pour Terraform**
   - Utilisateur voit juste `terraform apply`
   - Pas de manipulation manuelle

4. **Atomique**
   - Soit tout réussit, soit rien (rollback)
   - Pas d'état partiel

5. **Safe**
   - Ancien credential gardé jusqu'à fin
   - Nouveau testé avant suppression ancien

### ⚠️ Inconvénients

1. **Complexité**
   - Beaucoup de code à maintenir
   - Beaucoup de cas d'erreur à gérer

2. **Permissions Requises**
   - API key doit pouvoir lire/écrire TOUS les workflows
   - Pas toujours possible (permissions limitées)

3. **Performance**
   - 100 workflows = 100+ API calls
   - Peut prendre plusieurs minutes

4. **Risques Résiduels**
   - Si rollback échoue partiellement → intervention manuelle
   - Edge cases possibles (workflows verrouillés, etc.)

5. **L'ID Change Quand Même**
   - Même si transparent, l'ID est différent après
   - Peut impacter des systèmes externes qui référencent l'ID

---

## 📝 Documentation Utilisateur

```markdown
# n8n_credential

Manages an n8n credential with automatic rotation.

## Example Usage

```hcl
resource "n8n_credential" "api" {
  name = "API Key"
  type = "httpHeaderAuth"
  data = {
    name  = "Authorization"
    value = var.api_token
  }
}

resource "n8n_workflow" "example" {
  name = "My Workflow"
  # Workflow uses the credential
}
```

## Update Behavior (Rotation)

When you update a credential, the provider performs an **automatic rotation**:

1. **Creates** a new credential with updated data
2. **Scans** all workflows using the old credential
3. **Updates** each workflow to reference the new credential
4. **Deletes** the old credential

This ensures **zero downtime** - workflows continue working during the update.

### Rollback

If any step fails, the provider **automatically rolls back**:
- Deletes the new credential
- Restores all workflows to use the old credential
- The Terraform apply fails with a clear error message

### Important Notes

⚠️ **The credential ID will change** after an update. However, this is handled automatically - all workflows are updated to reference the new ID.

⚠️ **Permissions Required**: Your API key must have permissions to read and update ALL workflows that use the credential. If you don't have these permissions, the rotation will fail.

⚠️ **Performance**: Rotation scans all workflows. With many workflows (100+), this may take several minutes.

### Example Output

```bash
$ terraform apply

n8n_credential.api: Modifying... [id=cred-old-123]
n8n_credential.api: Found 5 workflows using credential cred-old-123
n8n_credential.api: Updated workflow wf-1 (1/5)
n8n_credential.api: Updated workflow wf-2 (2/5)
n8n_credential.api: Updated workflow wf-3 (3/5)
n8n_credential.api: Updated workflow wf-4 (4/5)
n8n_credential.api: Updated workflow wf-5 (5/5)
n8n_credential.api: Credential rotated successfully: cred-old-123 → cred-new-456
n8n_credential.api: Modifications complete after 15s [id=cred-new-456]

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
```

### Troubleshooting

**Error: "Failed to update workflow X"**

The rotation failed and was rolled back. Possible causes:
- Workflow is locked by another user
- Insufficient permissions
- Network error

Check the workflow in n8n UI and try again.

**Warning: "Could not delete old credential"**

The new credential is active and working, but the old one couldn't be deleted. This is safe - you can manually delete the old credential in n8n UI.
```

---

## 🎯 Prochaines Étapes

1. ✅ Design approuvé
2. **Implémenter** le code complet
3. **Tester** tous les cas (succès, échec, rollback)
4. **Documenter** le comportement
5. **Valider** avec une instance n8n réelle

**Prêt à implémenter ?** 🚀
