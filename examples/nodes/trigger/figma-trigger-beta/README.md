# Figma Trigger (Beta) Node Test

**Category**: Trigger **Type**: `n8n-nodes-base.figmaTrigger` **Latest Version**: 1

## Description

Starts the workflow when Figma events occur

## Node Information

- **Inputs**: main
- **Outputs**: nodeconnectiontypes.main
- **File**: `packages/nodes-base/nodes/Figma/FigmaTrigger.node.ts`

## Workflow Structure

```
Figma Trigger (Beta) → Display Result
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
2. Trigger the workflow (webhook, schedule, etc.)
3. Check the result in the Display Result node

## Customization

Edit the `parameters` in `main.tf` to customize the Figma Trigger (Beta) node behavior.

See [n8n Figma Trigger (Beta) documentation](https://docs.n8n.io/integrations/builtin/core-nodes/n8n-nodes-base.figma-trigger-beta/) for available parameters.

## Notes

- This is a trigger node, so it doesn't need an input
- The workflow will be triggered by figma trigger (beta) events
