package tag

// Values repeated across this package's data source, plural data source and
// resource schemas. They are the wire-level attribute names and the strings
// Terraform shows to users, so the three files must never drift apart.
const (
	attrName      = "name"
	attrCreatedAt = "created_at"
	attrUpdatedAt = "updated_at"
)

const (
	descCreatedAt = "Timestamp when the tag was created"
	descUpdatedAt = "Timestamp when the tag was last updated"
)
