# Facebook Lead Ads Trigger Node Test

**Category**: Trigger **Type**: `n8n-nodes-base.facebookLeadAdsTrigger` **Latest Version**: 1

## Description

Handle Facebook Lead Ads events via webhooks

## Node Information

- **Inputs**: main
- **Outputs**: nodeconnectiontypes.main
- **File**: `packages/nodes-base/nodes/FacebookLeadAds/FacebookLeadAdsTrigger.node.ts`

## Workflow Structure

```
Facebook Lead Ads Trigger → Display Result
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

Edit the `parameters` in `main.tf` to customize the Facebook Lead Ads Trigger node behavior.

See [n8n Facebook Lead Ads Trigger documentation](https://docs.n8n.io/integrations/builtin/core-nodes/n8n-nodes-base.facebook-lead-ads-trigger/) for available
parameters.

## Notes

- This is a trigger node, so it doesn't need an input
- The workflow will be triggered by facebook lead ads trigger events
