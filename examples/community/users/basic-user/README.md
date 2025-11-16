# User Management Example

This example demonstrates how to create and manage n8n users with different roles using Terraform.

## Features Demonstrated

- Creating users with different global roles (`global:admin`, `global:member`)
- Creating users with default instance role
- Querying users by ID
- Querying users by email
- Listing all users in the instance

## Important Notes

### API Limitations

The n8n user management API has several important limitations:

1. **Instance Owner Only**: Only the instance owner can create, update, and delete users
2. **Email Cannot Be Changed**: Once a user is created, their email address cannot be modified. You must delete and recreate the user to change the email
3. **Limited Updates**: Only the user's `role` can be updated. Fields like `first_name` and `last_name` are set by the user during account activation
4. **Pending Status**: Newly created users will have `is_pending = true` until they complete their account setup

### User Lifecycle

1. **Creation**: User is created with an email and optional role
2. **Invitation**: n8n sends an invitation email to the user
3. **Activation**: User completes setup (sets password, name, etc.)
4. **Active**: User can log in and use n8n

## Resources Created

- `n8n_user.admin` - Admin user with `global:admin` role
- `n8n_user.member` - Regular user with `global:member` role
- `n8n_user.default_role` - User with instance default role

## Data Sources Used

- `n8n_user` - Query single user by ID or email
- `n8n_users` - List all users in the instance

## Usage

```bash
# Initialize Terraform
terraform init

# Plan the changes
terraform plan \
  -var="n8n_api_key=YOUR_API_KEY" \
  -var="n8n_base_url=https://your-n8n-instance.com"

# Apply the configuration
terraform apply \
  -var="n8n_api_key=YOUR_API_KEY" \
  -var="n8n_base_url=https://your-n8n-instance.com"

# View outputs
terraform output

# Destroy resources
terraform destroy \
  -var="n8n_api_key=YOUR_API_KEY" \
  -var="n8n_base_url=https://your-n8n-instance.com"
```

## Expected Outputs

```
admin_user_id = "user-uuid-1"
admin_user_email = "admin-ci-timestamp@example.com"
admin_user_role = "global:admin"
member_user_id = "user-uuid-2"
member_user_email = "member-ci-timestamp@example.com"
member_user_role = "global:member"
default_user_id = "user-uuid-3"
default_user_role = "global:member"  # or whatever your instance default is
all_users_count = 4  # Including the instance owner
queried_admin_status = true  # Will be true until user activates their account
queried_member_status = true
```

## Available Roles

- `global:admin` - Full administrative access
- `global:member` - Regular user access (default)

## Troubleshooting

### "User already exists" Error

If you see a conflict error about a user already existing, it means a user with that email already exists in the instance. The example uses unique timestamps to
avoid this in CI/CD, but you may need to:

1. Use a different email address
2. Delete the existing user first
3. Import the existing user into Terraform state

### "Only instance owner can manage users" Error

This API endpoint is restricted to the instance owner. Make sure your API key belongs to the instance owner account.

### Users Not Appearing in Outputs

Newly created users will appear in the outputs immediately, but they won't be able to log in until they:

1. Receive and click the invitation email
2. Complete the account setup process
3. Set their password and personal information

## See Also

- [n8n User API Documentation](https://docs.n8n.io/api/users/)
- [Terraform n8n Provider Documentation](https://registry.terraform.io/providers/kodflow/n8n/latest/docs)
