# NocoDB Node Test

**Category**: Integration
**Type**: `n8n-nodes-base.nocoDb`
**Latest Version**: 1

## Description

Read, update, write and delete data from NocoDB

## Node Information

- **Inputs**: nodeconnectiontypes.main
- **Outputs**: nodeconnectiontypes.main
- **File**: `packages/nodes-base/nodes/NocoDB/NocoDB.node.ts`

## Workflow Structure

```
Manual Trigger → NocoDB → Display Result
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

Edit the `parameters` in `main.tf` to customize the NocoDB node behavior.

See [n8n NocoDB documentation](https://docs.n8n.io/integrations/builtin/core-nodes/n8n-nodes-base.nocodb/) for available parameters.

## Notes

- This workflow uses a manual trigger for testing
- In production, replace with appropriate trigger node
