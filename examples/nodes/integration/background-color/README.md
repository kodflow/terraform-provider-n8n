# Background Color Node Test

**Category**: Integration
**Type**: `n8n-nodes-base.Blur`
**Latest Version**: 1

## Description

Adds a blur to the image and so makes it less sharp

## Node Information

- **Inputs**: nodeconnectiontypes.main
- **Outputs**: nodeconnectiontypes.main
- **File**: `packages/nodes-base/nodes/EditImage/EditImage.node.ts`

## Workflow Structure

```
Manual Trigger → Background Color → Display Result
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

Edit the `parameters` in `main.tf` to customize the Background Color node behavior.

See [n8n Background Color documentation](https://docs.n8n.io/integrations/builtin/core-nodes/n8n-nodes-base.background-color/) for available parameters.

## Notes

- This workflow uses a manual trigger for testing
- In production, replace with appropriate trigger node
