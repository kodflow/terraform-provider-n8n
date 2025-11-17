# BambooHr Node Test

**Category**: Integration
**Type**: `n8n-nodes-base.bamboohr`
**Latest Version**: 1

## Description

Test workflow for the BambooHr node.

## Node Information

- **Inputs**: main
- **Outputs**: main
- **File**: `packages/nodes-base/nodes/BambooHr/BambooHr.node.ts`

## Workflow Structure

```
Manual Trigger → BambooHr → Display Result
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

Edit the `parameters` in `main.tf` to customize the BambooHr node behavior.

See [n8n BambooHr documentation](https://docs.n8n.io/integrations/builtin/core-nodes/n8n-nodes-base.bamboohr/) for available parameters.

## Notes

- This workflow uses a manual trigger for testing
- In production, replace with appropriate trigger node
