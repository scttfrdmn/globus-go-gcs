# ⚠️ Breaking Changes in v2.0

## Critical Security Change: Secret Input Method

### What Changed

**v2.0 removes the ability to pass secrets via command-line arguments.**

This is a **security-critical breaking change** required for HIPAA/PHI compliance.

### Why This Change Was Necessary

**Security Risk**: Passing secrets as command-line arguments has serious security implications:

1. **Process Listing Exposure**: Secrets are visible to any user who can run `ps` or view process lists
   ```bash
   # Any user on the system can see this:
   $ ps aux | grep globus
   user  12345  globus-connect-server ... --secret-access-key wJalrXUtnFEMI/K7MDENG...
   ```

2. **Shell History**: Secrets are saved in shell history files (~/.bash_history, ~/.zsh_history)
   ```bash
   $ history | grep secret
   42  globus-connect-server user-credential s3-keys add --secret-access-key wJalrXUtnFEMI/K7MDENG...
   ```

3. **Log Files**: Command-line arguments may be logged by system monitoring tools, audit systems, or debugging logs

4. **Non-Compliance**:
   - **NIST 800-53 IA-5(7)**: Prohibits embedded unprotected passwords
   - **HIPAA Security Rule § 164.312(a)(2)(iv)**: Requires encryption of authentication credentials
   - **PCI DSS 8.2.1**: Prohibits passwords in clear text

**This vulnerability makes v1.x unsuitable for healthcare, financial, or other regulated environments.**

### Affected Commands

Five commands that previously accepted secrets via CLI arguments:

1. `user-credential s3-keys add`
   - Flag removed: `--secret-access-key`
2. `user-credential s3-keys update`
   - Flag removed: `--secret-access-key`
3. `user-credential activescale-create`
   - Flag removed: `--password`
4. `oidc create`
   - Flag removed: `--client-secret`
5. `oidc update`
   - Flag removed: `--client-secret`

### Migration Guide

#### Option 1: Interactive Prompt (Recommended)

The command will prompt for the secret with hidden input:

```bash
# v1.x (insecure - DO NOT USE)
globus-connect-server user-credential s3-keys add \
  --access-key-id AKIA... \
  --secret-access-key wJalrXUtnFEMI/K7MDENG...

# v2.0 (secure)
globus-connect-server user-credential s3-keys add \
  --access-key-id AKIA...
# Prompts: Enter secret access key: ********
#          Confirm secret access key: ********
```

**Benefits**:
- No secrets in process list
- No secrets in shell history
- Input is hidden (not displayed on screen)
- Confirmation prevents typos

#### Option 2: Standard Input (For Scripts)

Pipe the secret to the command using `--secret-stdin`:

```bash
# From file (secure file with 0600 permissions)
cat /secure/secret.txt | globus-connect-server user-credential s3-keys add \
  --access-key-id AKIA... \
  --secret-stdin

# From environment variable
echo "$SECRET_VALUE" | globus-connect-server user-credential s3-keys add \
  --access-key-id AKIA... \
  --secret-stdin

# From password manager
pass show aws/secret-key | globus-connect-server user-credential s3-keys add \
  --access-key-id AKIA... \
  --secret-stdin

# From heredoc in script
globus-connect-server user-credential s3-keys add \
  --access-key-id AKIA... \
  --secret-stdin <<EOF
wJalrXUtnFEMI/K7MDENG...
EOF
```

**Benefits**:
- No secrets in command line
- Works in automated scripts
- Can integrate with password managers

#### Option 3: Environment Variable

Set the secret in an environment variable and use `--secret-env`:

```bash
# Set environment variable
export GLOBUS_SECRET_VALUE="wJalrXUtnFEMI/K7MDENG..."

# Run command
globus-connect-server user-credential s3-keys add \
  --access-key-id AKIA... \
  --secret-env

# Clear environment variable
unset GLOBUS_SECRET_VALUE
```

**Benefits**:
- Secrets not in command line
- Works in CI/CD with secret injection
- Simple for automation

**⚠️ Warning**: Environment variables are still visible in `/proc/<pid>/environ` on Linux. Use stdin method for highest security.

### Migration Checklist

- [ ] Identify all scripts using affected commands
- [ ] Update scripts to use one of the three secure methods
- [ ] Test updated scripts in non-production environment
- [ ] Update CI/CD pipelines to inject secrets securely
- [ ] Remove old command syntax from documentation
- [ ] Train team members on new secret input methods
- [ ] Audit logs for any leaked secrets in history files
- [ ] Consider rotating secrets after migration

### Migration Examples

#### Example 1: Shell Script

**Before (v1.x - INSECURE)**:
```bash
#!/bin/bash
# deploy.sh

ACCESS_KEY="AKIA..."
SECRET_KEY="wJalrXUtnFEMI/K7MDENG..."

globus-connect-server user-credential s3-keys add \
  --access-key-id "$ACCESS_KEY" \
  --secret-access-key "$SECRET_KEY"  # ❌ VISIBLE IN ps OUTPUT
```

**After (v2.0 - SECURE)**:
```bash
#!/bin/bash
# deploy.sh

ACCESS_KEY="AKIA..."
SECRET_KEY="wJalrXUtnFEMI/K7MDENG..."

echo "$SECRET_KEY" | globus-connect-server user-credential s3-keys add \
  --access-key-id "$ACCESS_KEY" \
  --secret-stdin  # ✅ NOT VISIBLE IN ps OUTPUT
```

#### Example 2: CI/CD Pipeline (GitHub Actions)

**Before (v1.x - INSECURE)**:
```yaml
- name: Add S3 credentials
  run: |
    globus-connect-server user-credential s3-keys add \
      --access-key-id ${{ secrets.AWS_ACCESS_KEY_ID }} \
      --secret-access-key ${{ secrets.AWS_SECRET_ACCESS_KEY }}  # ❌ LOGGED
```

**After (v2.0 - SECURE)**:
```yaml
- name: Add S3 credentials
  run: |
    echo "${{ secrets.AWS_SECRET_ACCESS_KEY }}" | \
      globus-connect-server user-credential s3-keys add \
        --access-key-id ${{ secrets.AWS_ACCESS_KEY_ID }} \
        --secret-stdin  # ✅ NOT LOGGED
```

#### Example 3: Ansible Playbook

**Before (v1.x - INSECURE)**:
```yaml
- name: Configure S3 credentials
  command: >
    globus-connect-server user-credential s3-keys add
    --access-key-id {{ aws_access_key }}
    --secret-access-key {{ aws_secret_key }}  # ❌ VISIBLE
  no_log: true  # Doesn't help - still visible in ps output
```

**After (v2.0 - SECURE)**:
```yaml
- name: Configure S3 credentials
  shell: |
    echo "{{ aws_secret_key }}" | \
      globus-connect-server user-credential s3-keys add \
        --access-key-id {{ aws_access_key }} \
        --secret-stdin  # ✅ SECURE
  no_log: true
```

#### Example 4: Terraform

**Before (v1.x - INSECURE)**:
```hcl
resource "null_resource" "s3_creds" {
  provisioner "local-exec" {
    command = <<-EOT
      globus-connect-server user-credential s3-keys add \
        --access-key-id ${var.aws_access_key_id} \
        --secret-access-key ${var.aws_secret_access_key}  # ❌ INSECURE
    EOT
  }
}
```

**After (v2.0 - SECURE)**:
```hcl
resource "null_resource" "s3_creds" {
  provisioner "local-exec" {
    command = <<-EOT
      echo "${var.aws_secret_access_key}" | \
        globus-connect-server user-credential s3-keys add \
          --access-key-id ${var.aws_access_key_id} \
          --secret-stdin  # ✅ SECURE
    EOT
  }
}
```

### Detection Script

Use this script to find instances of the old syntax in your codebase:

```bash
#!/bin/bash
# find-insecure-secret-usage.sh

echo "Scanning for insecure secret usage patterns..."
echo ""

patterns=(
  "--secret-access-key"
  "--client-secret [^-]"
  "--password [^-]"
)

for pattern in "${patterns[@]}"; do
  echo "Checking for: $pattern"
  grep -r "$pattern" . \
    --exclude-dir=.git \
    --exclude-dir=vendor \
    --exclude="*.md" \
    --color=always || echo "  ✓ Not found"
  echo ""
done

echo "Review any matches above and update to use --secret-stdin or --secret-env"
```

### Compatibility

| Version | Old Syntax | New Syntax | Notes |
|---------|-----------|------------|-------|
| v1.x | ✅ Supported (insecure) | ❌ Not available | Not HIPAA compliant |
| v2.0-beta | ⚠️ Deprecated (warning) | ✅ Supported | Transition period |
| v2.0+ | ❌ Removed | ✅ Required | HIPAA compliant |

### Frequently Asked Questions

#### Q: Why wasn't this change made in v1.x?

**A**: We prioritized feature parity with the Python version for v1.0. Security hardening for production/HIPAA compliance is the focus of v2.0.

#### Q: Can I temporarily enable the old behavior?

**A**: No. The old behavior is a security vulnerability and has been completely removed. There is no flag to re-enable it.

#### Q: What if I'm not in a regulated environment?

**A**: The security risk exists regardless of regulatory requirements. Even in non-regulated environments, secrets in command-line arguments can be compromised.

#### Q: Will v1.x continue to receive updates?

**A**: v1.x will receive critical bug fixes for 6 months after v2.0 release, but will not receive security updates. All users should migrate to v2.0.

#### Q: Is there a tool to automatically migrate my scripts?

**A**: We provide detection scripts (see above) to find instances. Migration must be done manually as it depends on your specific use case (interactive vs. automated).

#### Q: What about other sensitive data (tokens, credentials)?

**A**: Authentication tokens are already handled securely via the `login` command and stored encrypted. User credentials other than secrets (like usernames, IDs) are not considered sensitive and can remain as CLI arguments.

#### Q: Does this affect reading secrets from config files?

**A**: No. Reading secrets from properly secured config files (0600 permissions, encrypted) is acceptable. This change only affects command-line arguments.

### Security Best Practices

After migrating to v2.0, follow these additional security practices:

1. **Rotate Secrets**: If secrets may have been exposed in logs or history, rotate them
2. **Secure Storage**: Store secret files with 0600 permissions
3. **Clean History**: Clear shell history files after migration
   ```bash
   history -c  # Clear current session
   rm ~/.bash_history ~/.zsh_history  # Remove history files
   ```
4. **Audit Logs**: Review system logs for any recorded secrets
5. **Secret Managers**: Consider using password managers (1Password, LastPass) or secret managers (HashiCorp Vault, AWS Secrets Manager)
6. **Documentation**: Update team documentation to reflect new patterns
7. **Code Review**: Review and approve all changes that handle secrets

### Support

If you encounter issues during migration:

1. **Check the migration guide**: Review all three options (interactive, stdin, env)
2. **Review examples**: See the command-specific examples above
3. **Test in non-production**: Always test in a safe environment first
4. **Open an issue**: If you find a use case not covered, open a GitHub issue
5. **Security concerns**: For security-related questions, email security@globus.org

### Timeline

| Date | Event |
|------|-------|
| v2.0-beta.1 (Week 6) | Old syntax shows deprecation warning |
| v2.0-rc.1 (Week 7) | Old syntax removed completely |
| v2.0.0 (Week 8) | Production release with breaking change |
| v1.x EOL (Week 34) | v1.x critical bug fixes end (6 months after v2.0) |

### Acknowledgments

This change aligns with industry best practices:

- **AWS CLI**: Removed credentials from arguments, requires config files or env vars
- **kubectl**: Secrets only via files or stdin
- **docker**: Secrets via stdin or secret management
- **Terraform**: Sensitive values marked and not logged

We apologize for the inconvenience but believe this change is essential for operating in regulated environments and protecting sensitive data.

---

**Need Help?**
- Migration Issues: [Open a GitHub issue](https://github.com/scttfrdmn/globus-go-gcs/issues)
- Security Questions: security@globus.org
- Documentation: See [MIGRATION_V2.md](./MIGRATION_V2.md)
