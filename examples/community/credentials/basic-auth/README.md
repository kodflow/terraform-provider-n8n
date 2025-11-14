# HTTP Basic Auth Credential Example

Creates an HTTP Basic Auth credential that can be used in workflows for API authentication.

## What This Does

1. Creates an HTTP Basic Auth credential with username and password
2. Reads the credential back using a data source
3. Outputs credential information (without sensitive data)

## Usage

```bash
# Set your credentials
export TF_VAR_n8n_api_key="your-api-key-here"
export TF_VAR_basic_auth_user="your-username"
export TF_VAR_basic_auth_password="your-password"

# Initialize and apply
terraform init
terraform apply
```

## Using the Credential

Once created, this credential can be referenced in workflows by its ID:

```hcl
resource "n8n_workflow" "api_workflow" {
  nodes = jsonencode([
    {
      id   = "http-node"
      type = "n8n-nodes-base.httpRequest"
      credentials = {
        httpBasicAuth = {
          id   = n8n_credential.http_basic_auth.id
          name = n8n_credential.http_basic_auth.name
        }
      }
      parameters = {
        url    = "https://api.example.com/data"
        method = "GET"
      }
    }
  ])
}
```

## Cleanup

```bash
terraform destroy
```
