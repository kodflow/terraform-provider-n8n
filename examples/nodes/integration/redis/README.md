# Redis Node Test

**Category**: Integration **Type**: `n8n-nodes-base.redis` **Latest Version**: 1

## Description

Get, send and update data in Redis

## Node Information

- **Inputs**: nodeconnectiontypes.main
- **Outputs**: nodeconnectiontypes.main
- **File**: `packages/nodes-base/nodes/Redis/Redis.node.ts`

## Workflow Structure

```
Manual Trigger → Redis → Display Result
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

Edit the `parameters` in `main.tf` to customize the Redis node behavior.

See [n8n Redis documentation](https://docs.n8n.io/integrations/builtin/core-nodes/n8n-nodes-base.redis/) for available parameters.

## Notes

- This workflow uses a manual trigger for testing
- In production, replace with appropriate trigger node
