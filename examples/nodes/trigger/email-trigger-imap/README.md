# Email Trigger (IMAP) Node Test

**Category**: Trigger **Type**: `n8n-nodes-base.emailReadImap` **Latest Version**: 1

## Description

Triggers the workflow when a new email is received

## Node Information

- **Inputs**: main
- **Outputs**: main
- **File**: `packages/nodes-base/nodes/EmailReadImap/EmailReadImap.node.ts`

## Workflow Structure

```
Email Trigger (IMAP) → Display Result
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

Edit the `parameters` in `main.tf` to customize the Email Trigger (IMAP) node behavior.

See [n8n Email Trigger (IMAP) documentation](https://docs.n8n.io/integrations/builtin/core-nodes/n8n-nodes-base.email-trigger-imap/) for available parameters.

## Notes

- This is a trigger node, so it doesn't need an input
- The workflow will be triggered by email trigger (imap) events
