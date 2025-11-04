# Exemple d'utilisation du provider n8n

Ce dossier contient un exemple d'utilisation du provider Terraform pour n8n en développement local.

## Prérequis

1. Compiler et installer le provider localement :
   ```bash
   cd ..
   make build
   ```

2. Vérifier l'installation :
   ```bash
   ls -la ~/.terraform.d/plugins/registry.terraform.io/kodflow/n8n/0.0.1/
   ```

## Utilisation

1. Initialiser Terraform :
   ```bash
   terraform init
   ```

2. Voir le plan d'exécution :
   ```bash
   terraform plan
   ```

3. Appliquer la configuration :
   ```bash
   terraform apply
   ```

## Configuration

Le provider n8n nécessite la configuration suivante :

- `api_url` : URL de votre instance n8n (ex: `https://your-n8n.com`)
- `api_key` : Clé API n8n pour l'authentification

### Variables d'environnement

Vous pouvez également utiliser des variables d'environnement :

```bash
export N8N_API_URL="https://your-n8n.com"
export N8N_API_KEY="your-api-key"
```

### Fichier de variables

Créez un fichier `terraform.tfvars` (ignoré par Git) :

```hcl
n8n_api_url = "https://your-n8n.com"
n8n_api_key = "your-api-key"
```

## Nettoyage

Pour détruire les ressources créées :

```bash
terraform destroy
```
