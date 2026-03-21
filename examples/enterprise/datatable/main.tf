# Enterprise Data Table Example
# Requires n8n Enterprise license with data tables feature enabled
#
# Note: Data tables are only available with an Enterprise license.
# Columns cannot be changed after creation; to change columns, delete and recreate.

terraform {
  required_providers {
    n8n = {
      source  = "kodflow/n8n"
      version = ">= 2.0"
    }
  }
}

provider "n8n" {
  api_key  = var.n8n_api_key
  base_url = var.n8n_base_url
}

# ============================================================================
# Data Table Resource
# ============================================================================

resource "n8n_data_table" "contacts" {
  name = "ci_${var.run_id}_contacts"

  columns {
    name = "email"
    type = "String"
  }

  columns {
    name = "first_name"
    type = "String"
  }

  columns {
    name = "last_name"
    type = "String"
  }

  columns {
    name = "active"
    type = "Boolean"
  }
}

# ============================================================================
# Data Table Datasource (single)
# ============================================================================

data "n8n_data_table" "contacts_lookup" {
  id = n8n_data_table.contacts.id
}

# ============================================================================
# Data Tables Datasource (list)
# ============================================================================

data "n8n_data_tables" "all" {
  depends_on = [n8n_data_table.contacts]
}

# ============================================================================
# Outputs
# ============================================================================

output "data_table_id" {
  description = "The ID of the created data table"
  value       = n8n_data_table.contacts.id
}

output "data_table_name" {
  description = "The name of the created data table"
  value       = n8n_data_table.contacts.name
}

output "data_tables_count" {
  description = "Total number of data tables"
  value       = length(data.n8n_data_tables.all.data_tables)
}
