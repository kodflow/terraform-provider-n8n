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
| **Total Coverage** | **98.6%** |
| **Threshold** | **70.0%** |
| **Status** | ✅ PASSED |

---

## Acceptance Tests (E2E)

Acceptance tests validate the real behavior of the provider against an n8n instance.

| Resource | Status | Tests |
|----------|:------:|-------|
| `credential` | ✅ | `TestAccCredentialResource` |
| `tag` | ✅ | `TestAccTagResource` |
| `variable` | ✅ | `TestAccVariableResource` |
| `workflow` | ✅ | `TestAccWorkflowResource` |

**Legend:**
- ✅ Acceptance tests present
- Acceptance tests verify real operations via the n8n API

---

## Coverage by Package

| Icon | Package | Coverage |
|:----:|---------|----------|
| 🟢 | `github.com/kodflow/n8n/src/internal/provider` | 100.0% |
| 🟢 | `github.com/kodflow/n8n/src/internal/provider/credential` | 94.8% |
| 🟢 | `github.com/kodflow/n8n/src/internal/provider/execution` | 100.0% |
| 🟢 | `github.com/kodflow/n8n/src/internal/provider/project` | 99.1% |
| 🟢 | `github.com/kodflow/n8n/src/internal/provider/shared` | 100.0% |
| 🟢 | `github.com/kodflow/n8n/src/internal/provider/shared/client` | 100.0% |
| 🟢 | `github.com/kodflow/n8n/src/internal/provider/sourcecontrol` | 97.0% |
| 🟢 | `github.com/kodflow/n8n/src/internal/provider/tag` | 100.0% |
| 🟢 | `github.com/kodflow/n8n/src/internal/provider/user` | 100.0% |
| 🟢 | `github.com/kodflow/n8n/src/internal/provider/variable` | 98.1% |
| 🟢 | `github.com/kodflow/n8n/src/internal/provider/workflow` | 99.4% |

---

## Detailed Coverage by Category

This section lists only **public functions** (exported) to quickly identify untested functions.
Tables are organized by resource category to facilitate understanding of the provider architecture.

## 📦 Primary Resources (CRUD Entities)

Complete lifecycle management of n8n resources (Create, Read, Update, Delete).

| Function | credential | project | tag | user | variable | workflow |
|----------|:--------:|:--------:|:--------:|:--------:|:--------:|:--------:|
| `Configure` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |
| `Create` | 🟡 85.7% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |
| `Delete` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |
| `ImportState` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |
| `Metadata` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |
| `Read` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |
| `Schema` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |
| `Update` | 🟡 87.5% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |

---

## 🔧 Secondary Resources (Operations/Relations)

Special operations and resource relationship management.

### Source Control Pull

| Function | sourcecontrol |
|----------|:--------:|
| `Configure` | 🟢 100.0% |
| `Create` | 🟢 100.0% |
| `Delete` | 🟢 100.0% |
| `ImportState` | 🟢 100.0% |
| `Metadata` | 🟢 100.0% |
| `Read` | 🟢 100.0% |
| `Schema` | 🟢 100.0% |
| `Update` | 🟢 100.0% |

### Execution Retry

| Function | execution |
|----------|:--------:|
| `Configure` | 🟢 100.0% |
| `Create` | 🟢 100.0% |
| `Delete` | 🔵 N/A |
| `ImportState` | 🟢 100.0% |
| `Metadata` | 🟢 100.0% |
| `Read` | 🟢 100.0% |
| `Schema` | 🟢 100.0% |
| `Update` | 🟢 100.0% |

### Transfer

| Function | credential | workflow |
|----------|:--------:|:--------:|
| `Configure` | 🟢 100.0% | 🟢 100.0% |
| `Create` | 🟢 100.0% | 🟢 100.0% |
| `Delete` | 🔵 N/A | 🔵 N/A |
| `ImportState` | 🟢 100.0% | 🟢 100.0% |
| `Metadata` | 🟢 100.0% | 🟢 100.0% |
| `Read` | 🟢 100.0% | 🟢 100.0% |
| `Schema` | 🟢 100.0% | 🟢 100.0% |
| `Update` | 🟢 100.0% | 🟢 100.0% |

### Project User

| Function | project |
|----------|:--------:|
| `Configure` | 🟢 100.0% |
| `Create` | 🟢 100.0% |
| `Delete` | 🟢 100.0% |
| `ImportState` | 🟢 100.0% |
| `Metadata` | 🟢 100.0% |
| `Read` | 🟢 100.0% |
| `Schema` | 🟢 100.0% |
| `Update` | 🟢 100.0% |

---

## 📊 Data Sources

Reading n8n resources without managing their lifecycle.

### datasource (singular)

| Function | execution | project | tag | user | variable | workflow |
|----------|:--------:|:--------:|:--------:|:--------:|:--------:|:--------:|
| `Configure` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |
| `Metadata` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |
| `Read` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |
| `Schema` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |

### datasources (plural)

| Function | execution | project | tag | user | variable | workflow |
|----------|:--------:|:--------:|:--------:|:--------:|:--------:|:--------:|
| `Configure` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |
| `Metadata` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |
| `Read` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 95.2% |
| `Schema` | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% | 🟢 100.0% |


---

## Icon Legend

- 🟢 **≥90%** - Excellent coverage
- 🟡 **70-89%** - Good coverage
- 🟠 **1-69%** - Partial coverage (needs improvement)
- 🔴 **0%** - Untested function (implemented but no tests)
- 🔵 **N/A** - Not applicable (not in API or intentionally empty)

**Note:** Only public functions (exported) are listed. Private functions and constructors (`New*`) are excluded.

---

*Report generated on: 2025-11-14*
*Threshold: 70.0%*
