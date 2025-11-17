# crowd.dev Node Test

**Category**: Integration
**Type**: `n8n-nodes-base.crowdDev`
**Latest Version**: 1

## Description

crowd.dev is an open-source suite of community and data tools built to unlock community-led growth for your organization.

## Node Information

- **Inputs**: nodeconnectiontypes.main
- **Outputs**: nodeconnectiontypes.main
- **File**: `packages/nodes-base/nodes/CrowdDev/CrowdDev.node.ts`

## Workflow Structure

```
Manual Trigger → crowd.dev → Display Result
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

Edit the `parameters` in `main.tf` to customize the crowd.dev node behavior.

See [n8n crowd.dev documentation](https://docs.n8n.io/integrations/builtin/core-nodes/n8n-nodes-base.crowd-dev/) for available parameters.

## Notes

- This workflow uses a manual trigger for testing
- In production, replace with appropriate trigger node
