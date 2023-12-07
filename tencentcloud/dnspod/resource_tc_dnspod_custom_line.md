Provides a resource to create a dnspod custom_line

~> **NOTE:** Terraform uses the combined id of doamin and name when importing. When the name changes, the combined id will also change.

Example Usage

```hcl
resource "tencentcloud_dnspod_custom_line" "custom_line" {
    domain = "dnspod.com"
    name   = "testline8"
	area   = "6.6.6.1-6.6.6.2"
}
```

Import

dnspod custom_line can be imported using the id, e.g.

```
terraform import tencentcloud_dnspod_custom_line.custom_line domain#name
```