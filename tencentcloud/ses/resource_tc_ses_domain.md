Provides a resource to create a ses domain

Example Usage

```hcl
resource "tencentcloud_ses_domain" "domain" {
    email_identity = "iac.cloud"
}

```
Import

ses domain can be imported using the id, e.g.
```
$ terraform import tencentcloud_ses_domain.domain iac.cloud
```