package user

// Values repeated across this package's data source, plural data source and
// resource schemas. They are the wire-level attribute names and the strings
// Terraform shows to users, so the three files must never drift apart.
const (
	attrEmail     = "email"
	attrFirstName = "first_name"
	attrLastName  = "last_name"
	attrIsPending = "is_pending"
	attrRole      = "role"
	attrCreatedAt = "created_at"
	attrUpdatedAt = "updated_at"
)

const (
	descFirstName = "User's first name"
	descLastName  = "User's last name"
	descCreatedAt = "Timestamp when the user was created"
	descUpdatedAt = "Timestamp when the user was last updated"
)
