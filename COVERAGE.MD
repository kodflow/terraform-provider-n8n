# Coverage Report

Rapport de couverture généré automatiquement.

**Légende:**
- 🟢 ≥90% - Excellente couverture
- 🟡 70-89% - Bonne couverture
- 🔴 <70% - Couverture insuffisante

---

## Coverage Global

| Metric | Value |
|--------|-------|
| **Total Coverage** | **68.4%** |
| **Threshold** | **70.0%** |
| **Status** | ❌ FAILED |

---

## Coverage par Package

| Icon | Package | Coverage |
|:----:|---------|----------|
| 🟢 | `github.com/kodflow/n8n/src/internal/provider` | 100.0% |
| 🔴 | `github.com/kodflow/n8n/src/internal/provider/credential` | 59.9% |
| 🟡 | `github.com/kodflow/n8n/src/internal/provider/execution` | 79.2% |
| 🔴 | `github.com/kodflow/n8n/src/internal/provider/project` | 55.9% |
| 🟢 | `github.com/kodflow/n8n/src/internal/provider/shared` | 100.0% |
| 🟢 | `github.com/kodflow/n8n/src/internal/provider/shared/client` | 100.0% |
| 🟢 | `github.com/kodflow/n8n/src/internal/provider/sourcecontrol` | 90.8% |
| 🔴 | `github.com/kodflow/n8n/src/internal/provider/tag` | 42.1% |
| 🟡 | `github.com/kodflow/n8n/src/internal/provider/user` | 73.5% |
| 🟡 | `github.com/kodflow/n8n/src/internal/provider/variable` | 70.6% |
| 🟡 | `github.com/kodflow/n8n/src/internal/provider/workflow` | 83.9% |

---

## Coverage Détaillée par Fichier

Cette section liste uniquement les **fonctions publiques** (exportées) pour identifier rapidement les fonctions non testées.

### 📦 Package: `github.com/kodflow/n8n/src/internal/provider`

<details>
<summary><b>📄 resource.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🔴 `Create` | 0.0% |
| 🔴 `Read` | 0.0% |
| 🔴 `Update` | 0.0% |
| 🔴 `Delete` | 0.0% |
| 🟢 `ImportState` | 100.0% |

</details>

<details>
<summary><b>📄 transfer_resource.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🔴 `Metadata` | 0.0% |
| 🔴 `Schema` | 0.0% |
| 🟢 `Configure` | 100.0% |
| 🔴 `Create` | 0.0% |
| 🔴 `Read` | 0.0% |
| 🔴 `Update` | 0.0% |
| 🔴 `Delete` | 0.0% |
| 🔴 `ImportState` | 0.0% |

</details>

<details>
<summary><b>📄 datasource.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🔴 `Read` | 0.0% |

</details>

<details>
<summary><b>📄 datasources.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🔴 `Read` | 0.0% |

</details>

<details>
<summary><b>📄 retry_resource.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🟢 `Create` | 100.0% |
| 🟡 `Read` | 80.0% |
| 🔴 `Update` | 0.0% |
| 🔴 `Delete` | 0.0% |
| 🔴 `ImportState` | 0.0% |

</details>

<details>
<summary><b>📄 datasource.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🔴 `Read` | 0.0% |

</details>

<details>
<summary><b>📄 datasources.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🔴 `Read` | 0.0% |

</details>

<details>
<summary><b>📄 resource.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🔴 `Create` | 0.0% |
| 🔴 `Read` | 0.0% |
| 🔴 `Update` | 0.0% |
| 🔴 `Delete` | 0.0% |
| 🔴 `ImportState` | 0.0% |

</details>

<details>
<summary><b>📄 user_resource.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🔴 `Create` | 0.0% |
| 🔴 `Read` | 0.0% |
| 🔴 `Update` | 0.0% |
| 🔴 `Delete` | 0.0% |
| 🔴 `ImportState` | 0.0% |

</details>

<details>
<summary><b>📄 provider.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🟢 `Resources` | 100.0% |
| 🟢 `DataSources` | 100.0% |
| 🟢 `ValidateProvider` | 100.0% |

</details>

<details>
<summary><b>📄 pointers.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Ptr` | 100.0% |
| 🟢 `String` | 100.0% |
| 🟢 `Bool` | 100.0% |
| 🟢 `Int` | 100.0% |
| 🟢 `Int32` | 100.0% |
| 🟢 `Float32` | 100.0% |

</details>

<details>
<summary><b>📄 pull_resource.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🟡 `Create` | 78.6% |
| 🟢 `Read` | 100.0% |
| 🟢 `Update` | 100.0% |
| 🔴 `Delete` | 0.0% |
| 🟢 `ImportState` | 100.0% |

</details>

<details>
<summary><b>📄 datasource.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🔴 `Read` | 0.0% |

</details>

<details>
<summary><b>📄 datasources.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🔴 `Read` | 0.0% |

</details>

<details>
<summary><b>📄 resource.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🔴 `Create` | 0.0% |
| 🔴 `Read` | 0.0% |
| 🔴 `Update` | 0.0% |
| 🔴 `Delete` | 0.0% |
| 🔴 `ImportState` | 0.0% |

</details>

<details>
<summary><b>📄 datasource.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🔴 `Read` | 0.0% |

</details>

<details>
<summary><b>📄 datasources.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🔴 `Read` | 0.0% |

</details>

<details>
<summary><b>📄 resource.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🟢 `Create` | 100.0% |
| 🔴 `Read` | 0.0% |
| 🟢 `Update` | 100.0% |
| 🔴 `Delete` | 0.0% |
| 🔴 `ImportState` | 0.0% |

</details>

<details>
<summary><b>📄 datasource.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🔴 `Read` | 0.0% |

</details>

<details>
<summary><b>📄 datasources.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🔴 `Read` | 0.0% |

</details>

<details>
<summary><b>📄 resource.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🔴 `Create` | 0.0% |
| 🔴 `Read` | 0.0% |
| 🔴 `Update` | 0.0% |
| 🔴 `Delete` | 0.0% |
| 🔴 `ImportState` | 0.0% |

</details>

<details>
<summary><b>📄 datasource.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🔴 `Metadata` | 0.0% |
| 🔴 `Schema` | 0.0% |
| 🟢 `Configure` | 100.0% |
| 🔴 `Read` | 0.0% |

</details>

<details>
<summary><b>📄 datasources.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🔴 `Metadata` | 0.0% |
| 🔴 `Schema` | 0.0% |
| 🟢 `Configure` | 100.0% |
| 🔴 `Read` | 0.0% |

</details>

<details>
<summary><b>📄 resource.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🟢 `Create` | 100.0% |
| 🟢 `Read` | 100.0% |
| 🟢 `Update` | 100.0% |
| 🟠 `Delete` | 40.0% |
| 🔴 `ImportState` | 0.0% |

</details>

<details>
<summary><b>📄 transfer_resource.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🟢 `Create` | 100.0% |
| 🟢 `Read` | 100.0% |
| 🟢 `Update` | 100.0% |
| 🔴 `Delete` | 0.0% |
| 🟢 `ImportState` | 100.0% |

</details>


### 📦 Package: `credential`

<details>
<summary><b>📄 resource.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🔴 `Create` | 0.0% |
| 🔴 `Read` | 0.0% |
| 🔴 `Update` | 0.0% |
| 🔴 `Delete` | 0.0% |
| 🟢 `ImportState` | 100.0% |

</details>

<details>
<summary><b>📄 transfer_resource.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🔴 `Metadata` | 0.0% |
| 🔴 `Schema` | 0.0% |
| 🟢 `Configure` | 100.0% |
| 🔴 `Create` | 0.0% |
| 🔴 `Read` | 0.0% |
| 🔴 `Update` | 0.0% |
| 🔴 `Delete` | 0.0% |
| 🔴 `ImportState` | 0.0% |

</details>


### 📦 Package: `execution`

<details>
<summary><b>📄 datasource.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🔴 `Read` | 0.0% |

</details>

<details>
<summary><b>📄 datasources.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🔴 `Read` | 0.0% |

</details>

<details>
<summary><b>📄 retry_resource.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🟢 `Create` | 100.0% |
| 🟡 `Read` | 80.0% |
| 🔴 `Update` | 0.0% |
| 🔴 `Delete` | 0.0% |
| 🔴 `ImportState` | 0.0% |

</details>


### 📦 Package: `project`

<details>
<summary><b>📄 datasource.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🔴 `Read` | 0.0% |

</details>

<details>
<summary><b>📄 datasources.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🔴 `Read` | 0.0% |

</details>

<details>
<summary><b>📄 resource.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🔴 `Create` | 0.0% |
| 🔴 `Read` | 0.0% |
| 🔴 `Update` | 0.0% |
| 🔴 `Delete` | 0.0% |
| 🔴 `ImportState` | 0.0% |

</details>

<details>
<summary><b>📄 user_resource.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🔴 `Create` | 0.0% |
| 🔴 `Read` | 0.0% |
| 🔴 `Update` | 0.0% |
| 🔴 `Delete` | 0.0% |
| 🔴 `ImportState` | 0.0% |

</details>


### 📦 Package: `shared`

<details>
<summary><b>📄 pointers.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Ptr` | 100.0% |
| 🟢 `String` | 100.0% |
| 🟢 `Bool` | 100.0% |
| 🟢 `Int` | 100.0% |
| 🟢 `Int32` | 100.0% |
| 🟢 `Float32` | 100.0% |

</details>


### 📦 Package: `shared/client`


### 📦 Package: `sourcecontrol`

<details>
<summary><b>📄 pull_resource.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🟡 `Create` | 78.6% |
| 🟢 `Read` | 100.0% |
| 🟢 `Update` | 100.0% |
| 🔴 `Delete` | 0.0% |
| 🟢 `ImportState` | 100.0% |

</details>


### 📦 Package: `tag`

<details>
<summary><b>📄 datasource.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🔴 `Read` | 0.0% |

</details>

<details>
<summary><b>📄 datasources.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🔴 `Read` | 0.0% |

</details>

<details>
<summary><b>📄 resource.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🔴 `Create` | 0.0% |
| 🔴 `Read` | 0.0% |
| 🔴 `Update` | 0.0% |
| 🔴 `Delete` | 0.0% |
| 🔴 `ImportState` | 0.0% |

</details>


### 📦 Package: `user`

<details>
<summary><b>📄 datasource.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🔴 `Read` | 0.0% |

</details>

<details>
<summary><b>📄 datasources.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🔴 `Read` | 0.0% |

</details>

<details>
<summary><b>📄 resource.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🟢 `Create` | 100.0% |
| 🔴 `Read` | 0.0% |
| 🟢 `Update` | 100.0% |
| 🔴 `Delete` | 0.0% |
| 🔴 `ImportState` | 0.0% |

</details>


### 📦 Package: `variable`

<details>
<summary><b>📄 datasource.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🔴 `Read` | 0.0% |

</details>

<details>
<summary><b>📄 datasources.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🔴 `Read` | 0.0% |

</details>

<details>
<summary><b>📄 resource.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🔴 `Create` | 0.0% |
| 🔴 `Read` | 0.0% |
| 🔴 `Update` | 0.0% |
| 🔴 `Delete` | 0.0% |
| 🔴 `ImportState` | 0.0% |

</details>


### 📦 Package: `workflow`

<details>
<summary><b>📄 datasource.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🔴 `Metadata` | 0.0% |
| 🔴 `Schema` | 0.0% |
| 🟢 `Configure` | 100.0% |
| 🔴 `Read` | 0.0% |

</details>

<details>
<summary><b>📄 datasources.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🔴 `Metadata` | 0.0% |
| 🔴 `Schema` | 0.0% |
| 🟢 `Configure` | 100.0% |
| 🔴 `Read` | 0.0% |

</details>

<details>
<summary><b>📄 resource.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🟢 `Create` | 100.0% |
| 🟢 `Read` | 100.0% |
| 🟢 `Update` | 100.0% |
| 🟠 `Delete` | 40.0% |
| 🔴 `ImportState` | 0.0% |

</details>

<details>
<summary><b>📄 transfer_resource.go</b></summary>

| Function | Coverage |
|----------|----------|
| 🟢 `Metadata` | 100.0% |
| 🟢 `Schema` | 100.0% |
| 🟢 `Configure` | 100.0% |
| 🟢 `Create` | 100.0% |
| 🟢 `Read` | 100.0% |
| 🟢 `Update` | 100.0% |
| 🔴 `Delete` | 0.0% |
| 🟢 `ImportState` | 100.0% |

</details>



---

## Légende des Icônes

- 🟢 **≥90%** - Excellente couverture
- 🟡 **70-89%** - Bonne couverture
- 🟠 **1-69%** - Couverture partielle (à améliorer)
- 🔴 **0%** - Fonction non testée

**Note:** Seules les fonctions publiques (exportées) sont listées. Les fonctions privées et constructeurs (`New*`) sont exclus.

---

*Rapport généré le: 2025-11-11*
*Threshold: 70.0%*
