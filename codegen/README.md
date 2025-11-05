# N8N SDK Code Generation

Simple tooling for generating the n8n Go SDK from the official OpenAPI specification.

## 🚀 Quick Start

```bash
make openapi
```

That's it! This downloads the OpenAPI spec from GitHub, applies patches, and generates the SDK.

## 📁 Files

| File | Purpose |
|------|---------|
| `build-sdk.py` | Main orchestration script (~100 lines) |
| `generate-sdk.sh` | Runs openapi-generator and fixes imports |
| `openapi-generator-config.yaml` | OpenAPI generator settings |

## 🔧 How It Works

```bash
make openapi
  └─ build-sdk.py
      ├─ 1. Download from GitHub (sparse checkout)
      ├─ 2. Bundle YAML files (redocly)
      ├─ 3. Fix schema aliases (Workflow → workflow)
      ├─ 4. Apply patches (Python YAML operations)
      ├─ 5. Generate SDK (generate-sdk.sh → openapi-generator)
      └─ 6. Generate BUILD files (gazelle)
```

## 📝 Modifying Patches

**Patches are defined in `build-sdk.py` (lines 72-96):**

```python
# Patch workflow
workflow['properties']['versionId'] = {'type': 'string', 'readOnly': True}
workflow['properties']['isArchived'] = {'type': 'boolean', 'readOnly': True}
# ...
```

To modify patches:
1. Edit `build-sdk.py` directly
2. Run `make openapi`
3. Verify changes in generated `sdk/n8nsdk/*.go` files

## 🎯 Current Patches

**Workflow:** `versionId`, `isArchived`, `triggerCount`, `meta`, `pinData`
**Credential:** Remove `writeOnly` from `data`, add `isManaged`
**Create-Credential-Response:** Add `isManaged`
**Project:** Add `createdAt`, `updatedAt`, `icon`, `description`, `projectRelations`

## 🔄 Updating n8n Version

Edit `build-sdk.py`:

```python
N8N_COMMIT = "0ccf47044a2ba5b94140bfdd2ba36b868091288d"  # Change this
```

## 📊 File Flow

```
GitHub n8n repo (commit 0ccf4704)
    ↓
sdk/n8nsdk/api/openapi-source/ (committable source files)
    ↓ redocly bundle
sdk/n8nsdk/api/openapi.yaml (bundled)
    ↓ fix aliases + apply patches
sdk/n8nsdk/api/openapi.yaml (patched)
    ↓ openapi-generator
sdk/n8nsdk/*.go (generated SDK)
```

## 🐛 Troubleshooting

**Regenerate from scratch:**
```bash
rm -rf sdk/n8nsdk
mkdir -p sdk/n8nsdk/api
make openapi
```

**Java not found:**
```bash
sudo apt install openjdk-17-jre-headless  # Ubuntu/Debian
brew install openjdk@17                    # macOS
```

That's all you need to know!
