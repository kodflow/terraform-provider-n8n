package workflow

// Values repeated across this package's data source, plural data source and
// resource schemas. They are the wire-level attribute names and the strings
// Terraform shows to users, so the three files must never drift apart.
const (
	attrName   = "name"
	attrActive = "active"
)

const (
	descID     = "Workflow identifier"
	descName   = "Workflow name"
	descActive = "Whether the workflow is active"
)
