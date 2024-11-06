Provides a resource to create a ssl Check Certificate Domain Verification

~> **NOTE:** You can customize the maximum timeout time by setting parameter `timeouts`, which defaults to 15 minutes.

Example Usage

Check certificate domain

```hcl
resource "tencentcloud_ssl_check_certificate_domain_verification_operation" "example" {
  certificate_id = "6BE701Jx"
}
```

Check certificate domain and set the maximum timeout period

```hcl
resource "tencentcloud_ssl_check_certificate_domain_verification_operation" "example" {
  certificate_id = "6BE701Jx"

  timeouts {
    create = "30m"
  }
}
```