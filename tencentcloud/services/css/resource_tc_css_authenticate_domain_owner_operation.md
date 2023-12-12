Provides a resource to verify the domain ownership by specified way when DomainNeedVerifyOwner failed in domain creation.

Example Usage
dnsCheck way:
```hcl
resource "tencentcloud_css_authenticate_domain_owner_operation" "dnsCheck" {
  domain_name = "your_domain_name"
  verify_type = "dnsCheck"
}
```

fileCheck way:
```hcl
resource "tencentcloud_css_authenticate_domain_owner_operation" "fileCheck" {
  domain_name = "your_domain_name"
  verify_type = "fileCheck"
}
```