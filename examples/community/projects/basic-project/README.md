# Project Management Example

This example demonstrates how to create and manage n8n projects using Terraform.

## Features Demonstrated

- Creating multiple projects with unique names
- Querying projects by ID
- Listing all projects in the instance

## What Are n8n Projects?

Projects in n8n are organizational containers that help you:

- **Organize workflows** by environment (dev, staging, prod)
- **Group related workflows** by team, client, or purpose
- **Manage access** with project-specific permissions (Enterprise feature)
- **Isolate resources** for better organization

**Note**: Projects are NOT the same as Folders. Folders are a separate feature accessed via the internal n8n REST API (`/rest/*`) and are not available in the
public API.

## Important Notes

### API Limitations

The n8n project API has specific limitations:

1. **No GET by ID**: The API doesn't provide a direct endpoint to fetch a single project by ID
2. **POST Returns 201 with No Body**: Creating a project returns success but no project details
3. **Workaround Required**: The provider works around this by using the LIST endpoint and filtering by name

### Project Attributes

- `name` - **Required** - Project name (must be unique)
- `id` - **Computed** - Project identifier (UUID)
- `type` - **Computed** - Project type (e.g., "team")

## Resources Created

- `n8n_project.development` - Project for development workflows
- `n8n_project.production` - Project for production workflows
- `n8n_project.testing` - Project for testing workflows

## Data Sources Used

- `n8n_project` - Query single project by ID
- `n8n_projects` - List all projects in the instance

## Usage

```bash
# Initialize Terraform
terraform init

# Plan the changes
terraform plan \
  -var="n8n_api_key=YOUR_API_KEY" \
  -var="n8n_base_url=https://your-n8n-instance.com"

# Apply the configuration
terraform apply \
  -var="n8n_api_key=YOUR_API_KEY" \
  -var="n8n_base_url=https://your-n8n-instance.com"

# View outputs
terraform output

# Destroy resources
terraform destroy \
  -var="n8n_api_key=YOUR_API_KEY" \
  -var="n8n_base_url=https://your-n8n-instance.com"
```

## Expected Outputs

```
development_project_id = "project-uuid-1"
development_project_name = "ci-timestamp-dev-project"
development_project_type = "team"
production_project_id = "project-uuid-2"
production_project_name = "ci-timestamp-prod-project"
testing_project_id = "project-uuid-3"
all_projects_count = 4  # Including any pre-existing projects
queried_dev_project_name = "ci-timestamp-dev-project"
queried_dev_project_type = "team"
```

## Use Cases

### 1. Multi-Environment Setup

```hcl
resource "n8n_project" "dev" {
  name = "Development"
}

resource "n8n_project" "staging" {
  name = "Staging"
}

resource "n8n_project" "prod" {
  name = "Production"
}
```

### 2. Team-Based Organization

```hcl
resource "n8n_project" "engineering" {
  name = "Engineering Team"
}

resource "n8n_project" "marketing" {
  name = "Marketing Team"
}

resource "n8n_project" "sales" {
  name = "Sales Team"
}
```

### 3. Client Segregation

```hcl
resource "n8n_project" "client_a" {
  name = "Client A - Automation"
}

resource "n8n_project" "client_b" {
  name = "Client B - Integration"
}
```

## Combining with Other Resources

Projects can be combined with workflows, tags, and other resources:

```hcl
# Create project
resource "n8n_project" "my_project" {
  name = "My Project"
}

# Create workflows in the project
# Note: Project assignment for workflows requires Enterprise license
resource "n8n_workflow" "project_workflow" {
  name   = "Project Workflow"
  active = false
  # project_id = n8n_project.my_project.id  # Enterprise only
}

# Create project-specific tags
resource "n8n_tag" "project_tag" {
  name = "project:${n8n_project.my_project.name}"
}
```

## Troubleshooting

### "Project already exists" Error

If you see a conflict error, a project with that name already exists. The example uses unique timestamps to avoid this in CI/CD, but you may need to:

1. Use a different project name
2. Delete the existing project first
3. Import the existing project into Terraform state

### Project Not Found After Creation

Due to API limitations, there may be a brief delay between project creation and availability. The provider handles this automatically by polling the LIST
endpoint.

## Project vs. Folders

**Projects** (Available in Public API):

- ✅ Organizational containers
- ✅ Can be managed via Terraform
- ✅ Support team collaboration (Enterprise)

**Folders** (NOT in Public API):

- ❌ Sub-organization within projects
- ❌ Only available via internal REST API (`/rest/*`)
- ❌ Cannot be managed via this Terraform provider

## See Also

- [n8n Projects API Documentation](https://docs.n8n.io/api/projects/)
- [Terraform n8n Provider Documentation](https://registry.terraform.io/providers/kodflow/n8n/latest/docs)
