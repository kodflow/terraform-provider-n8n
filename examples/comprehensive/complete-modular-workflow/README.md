# Complete Modular Workflow Example

This is a **REAL, PRODUCTION-READY** workflow that demonstrates the full power of modular n8n workflow composition in Terraform.

## Workflow Architecture

```
┌─────────────┐
│  Webhook    │ ← HTTP POST /webhook/data-processor
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ Validate    │ ← JavaScript validation & enrichment
│  (Code)     │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ Check Valid │ ← Conditional routing
│   (If)      │
└──┬──────┬───┘
   │      │
   │      └─────────────┐
   │                    │
   ▼                    ▼
┌──────────┐      ┌──────────┐
│ Prepare  │      │ Prepare  │
│  Valid   │      │  Error   │
│  (Set)   │      │  (Set)   │
└────┬─────┘      └────┬─────┘
     │                 │
     ▼                 │
┌──────────┐           │
│   Send   │           │
│   to     │           │
│   API    │           │
│  (HTTP)  │           │
└────┬─────┘           │
     │                 │
     └────────┬────────┘
              │
              ▼
        ┌──────────┐
        │  Format  │
        │ Response │
        │  (Code)  │
        └──────────┘
```

## Node Categories Used

- **Trigger (1)**: Webhook - Receives HTTP requests
- **Core (5)**:
  - Code (2x) - Validation & Response formatting
  - If (1x) - Conditional routing
  - Set (2x) - Data transformation
- **Integration (1)**: HTTP Request - External API call

## What This Workflow Does

1. **Receives Data**: Webhook receives POST requests with JSON data
2. **Validates**: JavaScript code validates email format and enriches data
3. **Routes**: If node splits valid/invalid data to different paths
4. **Transforms**: Set nodes prepare data for next steps
5. **Integrates**: HTTP Request sends valid data to external API
6. **Responds**: Code node formats final response

## Usage

### 1. Initialize Terraform

```bash
terraform init
```

### 2. Plan

```bash
terraform plan
```

### 3. Apply

```bash
terraform apply
```

### 4. Test the Workflow

#### Valid Data Test:

```bash
curl -X POST http://localhost:5678/webhook/data-processor \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "score": 75}'
```

Expected response:

```json
{
  "success": true,
  "message": "Data processed successfully",
  "result": {
    "data": {
      "email": "user@example.com",
      "score": 75,
      "category": "high",
      "status": "success"
    }
  }
}
```

#### Invalid Data Test:

```bash
curl -X POST http://localhost:5678/webhook/data-processor \
  -H "Content-Type: application/json" \
  -d '{"email": "invalid-email", "score": 25}'
```

Expected response:

```json
{
  "success": true,
  "message": "Data processed successfully",
  "result": {
    "status": "error",
    "message": "Invalid email format"
  }
}
```

## Key Features Demonstrated

### 1. Modular Design

Each node is defined separately, making it easy to:

- Understand individual components
- Modify specific nodes without affecting others
- Reuse nodes in other workflows
- Test nodes in isolation

### 2. Branching Logic

The If node demonstrates conditional routing:

- Valid data → API integration path
- Invalid data → Error handling path

### 3. Data Transformation

Set nodes show how to:

- Select specific fields
- Add new fields
- Structure data for APIs

### 4. External Integration

HTTP Request node demonstrates:

- POST requests to external APIs
- Body parameters
- Response handling

### 5. Error Handling

The workflow handles both success and error cases gracefully.

## Terraform Benefits

### Before (Monolithic JSON):

```hcl
resource "n8n_workflow" "example" {
  nodes_json = jsonencode([
    { /* 50 lines of webhook config */ },
    { /* 30 lines of code config */ },
    { /* 40 lines of if config */ },
    // ... 300+ lines of unreadable JSON
  ])
}
```

### After (Modular):

```hcl
resource "n8n_workflow_node" "webhook" {
  name = "API Webhook"
  type = "n8n-nodes-base.webhook"
  parameters = jsonencode({ /* clear config */ })
}

resource "n8n_workflow_connection" "webhook_to_validate" {
  source_node = n8n_workflow_node.webhook.name
  target_node = n8n_workflow_node.validate.name
}
```

## Customization

### Change Validation Logic

Edit `validate_data` node parameters:

```terraform
resource "n8n_workflow_node" "validate_data" {
  parameters = jsonencode({
    jsCode = <<-EOT
      // Your custom validation logic here
    EOT
  })
}
```

### Add More Integration Steps

Add new HTTP Request nodes and connections:

```terraform
resource "n8n_workflow_node" "another_api" {
  name = "Send to Another API"
  type = "n8n-nodes-base.httpRequest"
  // ...
}
```

### Change Routing Logic

Modify the If node conditions:

```terraform
resource "n8n_workflow_node" "check_validity" {
  parameters = jsonencode({
    conditions = {
      number = [{
        value1 = "={{ $json.score }}"
        operation = "larger"
        value2 = 80
      }]
    }
  })
}
```

## Files

- `main.tf` - Main workflow definition
- `variables.tf` - Input variables
- `README.md` - This file

## Next Steps

1. Activate the workflow: `terraform apply -var="active=true"`
2. Monitor executions in n8n UI
3. Extend with more nodes (databases, notifications, etc.)
4. Add error handling with Error Trigger node
5. Implement retry logic for failed API calls

## Notes

- This workflow uses httpbin.org for testing (no auth required)
- Replace with your actual API endpoints in production
- Consider adding credentials management for real APIs
- Test thoroughly before activating in production
