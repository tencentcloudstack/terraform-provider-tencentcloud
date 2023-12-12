Provides a resource to create a ses verify_domain

~> **NOTE:** Please add the `attributes` information returned by `tencentcloud_ses_domain` to the domain name resolution record through `tencentcloud_dnspod_record`, and then verify it.

Example Usage

```hcl
resource "tencentcloud_ses_verify_domain" "verify_domain" {
  email_identity = "example.com"
}
```