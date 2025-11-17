# Form.io Trigger Node Test

**Category**: Trigger
**Type**: `n8n-nodes-base.formIoTrigger`
**Latest Version**: 1

## Description

Handle form.io events via webhooks

## Node Information

- **Inputs**: main
- **Outputs**: nodeconnectiontypes.main
- **File**: `packages/nodes-base/nodes/FormIo/FormIoTrigger.node.ts`

## Workflow Structure

```
Form.io Trigger → Display Result
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

Edit the `parameters` in `main.tf` to customize the Form.io Trigger node behavior.

See [n8n Form.io Trigger documentation](https://docs.n8n.io/integrations/builtin/core-nodes/n8n-nodes-base.form-io-trigger/) for available parameters.

## Notes

- This is a trigger node, so it doesn't need an input
- The workflow will be triggered by form.io trigger events
