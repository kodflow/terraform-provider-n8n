output "admin_user_id" {
  description = "ID of the admin user"
  value       = n8n_user.admin.id
}

output "admin_user_email" {
  description = "Email of the admin user"
  value       = n8n_user.admin.email
}

output "admin_user_role" {
  description = "Role of the admin user"
  value       = n8n_user.admin.role
}

output "member_user_id" {
  description = "ID of the member user"
  value       = n8n_user.member.id
}

output "member_user_email" {
  description = "Email of the member user"
  value       = n8n_user.member.email
}

output "member_user_role" {
  description = "Role of the member user"
  value       = n8n_user.member.role
}

output "default_user_id" {
  description = "ID of the user with default role"
  value       = n8n_user.default_role.id
}

output "default_user_role" {
  description = "Role of the user with default role"
  value       = n8n_user.default_role.role
}

output "all_users_count" {
  description = "Total number of users in the instance"
  value       = length(data.n8n_users.all.users)
}

output "queried_admin_status" {
  description = "Admin user queried by ID - is_pending status"
  value       = data.n8n_user.admin_by_id.is_pending
}

output "queried_member_status" {
  description = "Member user queried by email - is_pending status"
  value       = data.n8n_user.member_by_email.is_pending
}
