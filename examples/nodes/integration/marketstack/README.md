# Marketstack Node Test

**Category**: Integration **Type**: `n8n-nodes-base.marketstack` **Latest Version**: 1

## Description

Consume Marketstack API

## Node Information

- **Inputs**: nodeconnectiontypes.main
- **Outputs**: nodeconnectiontypes.main
- **File**: `packages/nodes-base/nodes/Marketstack/Marketstack.node.ts`

## Workflow Structure

```
Manual Trigger → Marketstack → Display Result
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

Edit the `parameters` in `main.tf` to customize the Marketstack node behavior.

See [n8n Marketstack documentation](https://docs.n8n.io/integrations/builtin/core-nodes/n8n-nodes-base.marketstack/) for available parameters.

## Notes

- This workflow uses a manual trigger for testing
- In production, replace with appropriate trigger node
