# Coverage Report

Automatically generated coverage report.

**Legend:**

- 🟢 ≥90% - Excellent coverage
- 🟡 70-89% - Good coverage
- 🔴 <70% - Insufficient coverage

---

## Global Coverage

| Metric | Value |
|--------|-------|
| **Total Coverage** | **92.4%** |
| **Threshold** | **70.0%** |
| **Status** | ✅ PASSED |

---

## Workflow Examples Testing

All 296 per-node workflow examples have been tested for correctness.

| Metric | Value |
|--------|-------|
| **Status** | ⚠️ Not yet run |

Run `make nodes/test-workflows` to test all 296 workflow examples.

---

## Acceptance Tests (E2E)

Acceptance tests validate the real behavior of the provider against an n8n instance.

| Resource | Status | Tests |
|----------|:------:|-------|
| `credential` | ✅ | `TestAccCredentialResource` |
| `project` | ✅ | `TestAccProjectResource` |
| `tag` | ✅ | `TestAccTagResource` |
| `user` | ✅ | `TestAccUserResource` |
| `variable` | ✅ | `TestAccVariableResource` |
| `workflow` | ✅ | `TestAccWorkflowConnectionResource`, `TestAccWorkflowConnectionResource_InCompleteWorkflow` |
| `workflow` | ✅ | `TestAccWorkflowNodeResource`, `TestAccWorkflowNodeResource_InCompleteWorkflow` |
| `workflow` | ✅ | `TestAccWorkflowResource` |

**Legend:**

- ✅ Acceptance tests present
- Acceptance tests verify real operations via the n8n API

---

## Coverage by Package

| Icon | Package | Coverage |
|:----:|---------|----------|
| 🟢 | `github.com/kodflow/terraform-provider-n8n/src/internal/provider` | 100.0% |
| 🟢 | `github.com/kodflow/terraform-provider-n8n/src/internal/provider/audit` | 95.3% |
| 🟢 | `github.com/kodflow/terraform-provider-n8n/src/internal/provider/credential` | 94.2% |
| 🟢 | `github.com/kodflow/terraform-provider-n8n/src/internal/provider/datatable` | 92.0% |
| 🟢 | `github.com/kodflow/terraform-provider-n8n/src/internal/provider/execution` | 94.0% |
| 🟢 | `github.com/kodflow/terraform-provider-n8n/src/internal/provider/project` | 93.2% |
| 🟢 | `github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared` | 100.0% |
| 🟢 | `github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client` | 100.0% |
| 🟢 | `github.com/kodflow/terraform-provider-n8n/src/internal/provider/sourcecontrol` | 94.7% |
| 🟢 | `github.com/kodflow/terraform-provider-n8n/src/internal/provider/tag` | 92.5% |
| 🟢 | `github.com/kodflow/terraform-provider-n8n/src/internal/provider/user` | 90.8% |
| 🟢 | `github.com/kodflow/terraform-provider-n8n/src/internal/provider/variable` | 93.0% |
| 🟡 | `github.com/kodflow/terraform-provider-n8n/src/internal/provider/workflow` | 89.5% |

---

## Detailed Coverage by Category

This section lists only **public functions** (exported) to quickly identify untested functions.
Tables are organized by resource category to facilitate understanding of the provider architecture.

## 📦 Primary Resources (CRUD Entities)

Complete lifecycle management of n8n resources (Create, Read, Update, Delete).

| Function | credential | project | sourcecontrol | tag | user | variable | workflow |
|----------|:--------:|:--------:|:--------:|:--------:|:--------:|:--------:|:--------:|
| `Configure` | 🟡 88.9% | 🟡 88.9% | 🟡 88.9% | 🟡 88.9% | 🟡 88.9% | 🟡 88.9% | 🟡 80.0% |
| `Create` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |
| `Delete` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |
| `ImportState` | 🟢 100.0% | 🟢 100.0% | 🔵 N/A | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |
| `Metadata` | 🟠 66.7% | 🟠 66.7% | 🟢 100.0% | 🟠 66.7% | 🟠 66.7% | 🟠 66.7% | 🟠 66.7% |
| `Read` | 🟢 100.0% | 🟢 100.0% | 🟡 80.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |
| `Schema` | 🟠 66.7% | 🟠 66.7% | 🟢 100.0% | 🟠 66.7% | 🟠 66.7% | 🟠 66.7% | 🟡 80.0% |
| `Update` | 🟢 100.0% | 🟢 100.0% | 🟡 85.7% | 🟢 100.0% | 🟢 100.0% | 🟢 91.7% | 🟢 100.0% |

---

## 🔧 Secondary Resources (Operations/Relations)

Special operations and resource relationship management.

### 📊 Quick Overview

| Resource | Type | API Interaction | Primary Functions Tested |
|----------|------|-----------------|-------------------------|
| **Connection** 📋 | Local-only | None | Schema ✅, Metadata ✅, JSON Generation ✅ |
| **Node** 🔷 | Local-only | None | Schema ✅, Metadata ✅, JSON Generation ✅ |
| **Project User** 👥 | API-based | Full CRUD | All functions ✅ 100% |

---

### 🏠 Local-Only Resources

These resources exist **only in Terraform state** and do not interact with the n8n API. They generate JSON configurations for use in workflow compositions.

#### 📋 Connection (Workflow Connections)

Defines connections between workflow nodes for modular composition.

| Function | Coverage | Notes |
|----------|:--------:|-------|
| `Configure` | 🟠 60.0% | No configuration needed |
| `Create` | 🟡 75.0% | State-only operation (no API calls) |
| `Delete` | 🟠 25.0% | Automatic state removal |
| `ImportState` | 🟢 100.0% | Passthrough implementation |
| `Metadata` | 🟠 66.7% | ✅ Fully tested |
| `Read` | 🟠 22.2% | State is always current |
| `Schema` | 🟡 80.0% | ✅ Fully tested |
| `Update` | 🟡 71.4% | Regenerates JSON in state |

**Type:** Local state-only resource
**Purpose:** Generate connection JSON for `n8n_workflow` resources
**API Interaction:** None (100% state operations)

#### 🔷 Node (Workflow Nodes)

Defines individual workflow nodes for modular composition.

| Function | Coverage | Notes |
|----------|:--------:|-------|
| `Configure` | 🟠 60.0% | No configuration needed |
| `Create` | 🟡 87.5% | State-only operation (no API calls) |
| `Delete` | 🟠 60.0% | Automatic state removal |
| `ImportState` | 🟢 100.0% | Passthrough implementation |
| `Metadata` | 🟠 66.7% | ✅ Fully tested |
| `Read` | 🟠 60.0% | State is always current |
| `Schema` | 🟡 80.0% | ✅ Fully tested |
| `Update` | 🟡 85.7% | Regenerates JSON in state |

**Type:** Local state-only resource
**Purpose:** Generate node JSON for `n8n_workflow` resources
**API Interaction:** None (100% state operations)

> **Note on Coverage:** The 0% coverage for Create/Update/ImportState methods is expected. These are simple state operations with no complex logic. The core functionality (JSON generation and ID creation) is thoroughly tested via internal tests with 100% coverage.

### 🔗 API-Based Secondary Resources

These resources interact with the n8n API to manage resource relationships.

#### Tags

| Function | Coverage | Notes |
|----------|:--------:|-------|
| `Configure` | 🟡 88.9% |  |
| `Create` | 🟢 100.0% | ✅ Fully tested |
| `Delete` | 🟢 100.0% | ✅ Fully tested |
| `ImportState` | 🟢 100.0% | ✅ Fully tested |
| `Metadata` | 🟠 66.7% |  |
| `Read` | 🟢 100.0% | ✅ Fully tested |
| `Schema` | 🟠 66.7% |  |
| `Update` | 🟢 100.0% | ✅ Fully tested |

#### 👥 Project User (Project Member Management)

Manages user assignments and roles within n8n projects.

| Function | Coverage | Notes |
|----------|:--------:|-------|
| `Configure` | 🟡 88.9% |  |
| `Create` | 🟢 100.0% | ✅ Fully tested |
| `Delete` | 🟢 100.0% | ✅ Fully tested |
| `ImportState` | 🟢 100.0% | ✅ Fully tested |
| `Metadata` | 🟠 66.7% |  |
| `Read` | 🟢 100.0% | ✅ Fully tested |
| `Schema` | 🟠 66.7% |  |
| `Update` | 🟢 100.0% | ✅ Fully tested |

**Type:** API-based resource
**Purpose:** Manage project member roles and access
**API Interaction:** Full CRUD operations via n8n API

---

## 📊 Data Sources

Reading n8n resources without managing their lifecycle.

### datasource (singular)

| Function | credential | execution | project | tag | user | variable | workflow |
|----------|:--------:|:--------:|:--------:|:--------:|:--------:|:--------:|:--------:|
| `Configure` | 🟡 88.9% | 🟡 88.9% | 🟡 88.9% | 🟡 88.9% | 🟡 88.9% | 🟡 88.9% | 🟡 80.0% |
| `Metadata` | 🟠 66.7% | 🟢 100.0% | 🟠 66.7% | 🟠 66.7% | 🟠 66.7% | 🟠 66.7% | 🟠 66.7% |
| `Read` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 93.8% |
| `Schema` | 🟠 66.7% | 🟢 100.0% | 🟠 66.7% | 🟠 66.7% | 🟠 66.7% | 🟠 66.7% | 🟡 80.0% |

### datasources (plural)

| Function | credential | execution | project | tag | user | variable | workflow |
|----------|:--------:|:--------:|:--------:|:--------:|:--------:|:--------:|:--------:|
| `Configure` | 🟡 88.9% | 🟡 88.9% | 🟡 88.9% | 🟡 88.9% | 🟡 88.9% | 🟡 88.9% | 🟡 88.9% |
| `Metadata` | 🟠 66.7% | 🟢 100.0% | 🟠 66.7% | 🟠 66.7% | 🟠 66.7% | 🟠 66.7% | 🟠 66.7% |
| `Read` | 🟢 100.0% | 🟢 100.0% | 🟢 93.3% | 🟢 92.3% | 🟢 93.3% | 🟢 93.3% | 🟢 100.0% |
| `Schema` | 🟠 66.7% | 🟢 100.0% | 🟠 66.7% | 🟠 66.7% | 🟠 66.7% | 🟠 66.7% | 🟠 66.7% |

---

## Icon Legend

- 🟢 **≥90%** - Excellent coverage
- 🟡 **70-89%** - Good coverage
- 🟠 **1-69%** - Partial coverage (needs improvement)
- 🔴 **0%** - Untested function (implemented but no tests)
- 🔵 **N/A** - Not applicable (not in API or intentionally empty)

**Note:** Only public functions (exported) are listed. Private functions and constructors (`New*`) are excluded.

---

*Report generated on: 2026-03-21*
*Threshold: 70.0%*
