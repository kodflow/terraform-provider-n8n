# Postmark Trigger Node Test

**Category**: Trigger
**Type**: `n8n-nodes-base.postmarkTrigger`
**Latest Version**: 1

## Description

Starts the workflow when Postmark events occur

## Node Information

- **Inputs**: main
- **Outputs**: nodeconnectiontypes.main
- **File**: `packages/nodes-base/nodes/Postmark/PostmarkTrigger.node.ts`

## Workflow Structure

```
Postmark Trigger → Display Result
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

Edit the `parameters` in `main.tf` to customize the Postmark Trigger node behavior.

See [n8n Postmark Trigger documentation](https://docs.n8n.io/integrations/builtin/core-nodes/n8n-nodes-base.postmark-trigger/) for available parameters.

## Notes

- This is a trigger node, so it doesn't need an input
- The workflow will be triggered by postmark trigger events
