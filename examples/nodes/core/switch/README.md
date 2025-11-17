# Switch Node Test

**Category**: Core
**Type**: `n8n-nodes-base.switch`
**Latest Version**: 1

## Description

Route items depending on defined expression or rules

## Node Information

- **Inputs**: main
- **Outputs**: main
- **File**: `packages/nodes-base/nodes/Switch/Switch.node.ts`

## Workflow Structure

```
Manual Trigger → Switch → Display Result
```

## Usage

### 1. Initialize Terraform

```bash
terraform init
```

### 2. Plan

```bash
terraform plan -var="n8n_api_key=YOUR_API_KEY"
```

### 3. Apply

```bash
terraform apply -var="n8n_api_key=YOUR_API_KEY"
```

### 4. Test in n8n

1. Open the workflow in n8n UI
2. Click "Execute Workflow"
3. Check the result in the Display Result node

## Customization

Edit the `parameters` in `main.tf` to customize the Switch node behavior.

See [n8n Switch documentation](https://docs.n8n.io/integrations/builtin/core-nodes/n8n-nodes-base.switch/) for available parameters.

## Notes

- This workflow uses a manual trigger for testing
- In production, replace with appropriate trigger node
