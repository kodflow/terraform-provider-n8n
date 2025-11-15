<!-- Copyright (c) 2024 Florent (Kodflow). All rights reserved. -->
<!-- Licensed under the Sustainable Use License 1.0 -->
<!-- See LICENSE.md in the project root for license information. -->

# Security Policy

## Supported Versions

We release patches for security vulnerabilities in the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

**Note:** Only the latest minor version of each major version receives security updates.

## Reporting a Vulnerability

We take the security of the n8n Terraform Provider seriously. If you believe you have found a security vulnerability, please report it to us responsibly.

### How to Report

**DO NOT** report security vulnerabilities through public GitHub issues.

Instead, please report them via one of the following methods:

1. **Email** (Preferred): Send details to `133899878+kodflow@users.noreply.github.com`
2. **GitHub Security Advisories**: Use the [Security tab](https://github.com/kodflow/terraform-provider-n8n/security/advisories) to report privately

### What to Include

Please include the following information in your report:

- **Type of vulnerability** (e.g., injection, authentication bypass, etc.)
- **Affected version(s)** of the provider
- **Step-by-step instructions** to reproduce the issue
- **Proof of concept** or exploit code (if possible)
- **Impact assessment** - What can an attacker do with this vulnerability?
- **Suggested fix** (if you have one)

### Example Report Format

```markdown
## Vulnerability Type

[e.g., SQL Injection, Authentication Bypass, etc.]

## Affected Versions

[e.g., v1.0.0 - v1.2.3]

## Description

[Clear description of the vulnerability]

## Steps to Reproduce

1. Configure provider with...
2. Create resource...
3. Execute action...
4. Observe vulnerability...

## Impact

[What can an attacker do? What data is at risk?]

## Proof of Concept

[Code or steps to demonstrate the vulnerability]

## Suggested Fix

[If you have suggestions for fixing the issue]
```

## Response Timeline

We are committed to responding to security reports promptly:

- **Initial Response**: Within 48 hours of receiving your report
- **Status Update**: Weekly updates on the progress of fixing the issue
- **Fix Timeline**: We aim to release a fix within 30 days for critical vulnerabilities
- **Disclosure**: Coordinated disclosure after a fix is available

## Disclosure Policy

We follow **Coordinated Vulnerability Disclosure**:

1. **You report** the vulnerability privately
2. **We acknowledge** and investigate
3. **We develop** a fix
4. **We release** a security patch
5. **We publicly disclose** the vulnerability (crediting you, if desired)

### Timeline

- **Day 0**: Vulnerability reported
- **Day 1-2**: Initial acknowledgment
- **Day 3-30**: Investigation and fix development
- **Day 30**: Security patch released
- **Day 37**: Public disclosure (7 days after patch release)

We may adjust this timeline based on the complexity of the issue.

## Security Best Practices

### For Users

When using the n8n Terraform Provider:

1. **Keep Updated**: Always use the latest version
2. **Secure Credentials**: Never commit API keys or credentials to version control
3. **Use Environment Variables**: Store sensitive data in environment variables or secure vaults
4. **Enable Logging**: Use `TF_LOG=DEBUG` only in secure environments
5. **Review Plans**: Always review `terraform plan` output before applying
6. **Limit Access**: Restrict provider permissions to the minimum required
7. **Audit Regularly**: Review your Terraform state files and configurations

### For Contributors

When developing the provider:

1. **Input Validation**: Always validate and sanitize user input
2. **Secure Defaults**: Use secure defaults for all configurations
3. **Error Handling**: Don't expose sensitive information in error messages
4. **Dependencies**: Keep dependencies updated and audited
5. **Code Review**: All code must be reviewed before merging
6. **Testing**: Include security-focused tests
7. **Secrets**: Never commit secrets, API keys, or credentials

### Sensitive Data Handling

The provider handles sensitive data including:

- n8n API keys
- Workflow credentials
- User authentication tokens
- Configuration data

**Best practices:**

```hcl
# ✅ Good - Use environment variables
provider "n8n" {
  api_url = var.n8n_api_url
  api_key = var.n8n_api_key
}

# ❌ Bad - Hardcoded credentials
provider "n8n" {
  api_url = "https://n8n.example.com"
  api_key = "n8n_api_key_1234567890"  # NEVER DO THIS
}
```

## Known Security Considerations

### API Key Management

- **API keys are transmitted** to the n8n instance via HTTPS
- **Keys are not logged** unless debug logging is explicitly enabled
- **Keys are stored** in Terraform state files (encrypt state files!)

### State File Security

Terraform state files may contain sensitive information:

- Use [remote state backends](https://developer.hashicorp.com/terraform/language/state/remote) with encryption
- Enable [state encryption](https://developer.hashicorp.com/terraform/language/state/encryption)
- Restrict access to state files
- Never commit state files to version control

### Network Security

- All API communication uses HTTPS
- Certificate validation is enforced
- No sensitive data is sent via query parameters

## Security Updates

Security updates will be:

1. Released as patch versions (e.g., 1.2.3 → 1.2.4)
2. Documented in the changelog with a `[SECURITY]` tag
3. Announced via:
   - GitHub Security Advisories
   - GitHub Releases
   - Project README

## Hall of Fame

We appreciate security researchers who responsibly disclose vulnerabilities. Contributors will be credited here (with permission):

<!-- Contributors will be listed here after disclosure -->

---

## Questions?

If you have questions about this security policy, please open a discussion in the
[GitHub Discussions](https://github.com/kodflow/terraform-provider-n8n/discussions) or contact us at `133899878+kodflow@users.noreply.github.com`.

Thank you for helping keep the n8n Terraform Provider secure! 🔒
