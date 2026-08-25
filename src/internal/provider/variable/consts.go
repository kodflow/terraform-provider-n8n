package variable

// Values repeated across this package's data source, plural data source and
// resource schemas. They are the wire-level attribute names and the strings
// Terraform shows to users, so the three files must never drift apart.
const (
	attrKey       = "key"
	attrValue     = "value"
	attrType      = "type"
	attrProjectID = "project_id"
)

const (
	descValue = "Variable value"
)
