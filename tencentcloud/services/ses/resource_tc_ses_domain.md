Provides a resource to create a ses domain

Example Usage

```hcl
resource "tencentcloud_ses_domain" "domain" {
    email_identity = "iac.cloud"
    dkim_option    = 1
    tag_key        = "env"
    tag_value      = "prod"
}

```
Import

ses domain can be imported using the id, e.g.
```
$ terraform import tencentcloud_ses_domain.domain iac.cloud
```