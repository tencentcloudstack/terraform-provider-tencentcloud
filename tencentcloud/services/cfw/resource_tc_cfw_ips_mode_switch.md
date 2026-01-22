Provides a resource to create a CFW ips mode switch

Example Usage

```hcl
resource "tencentcloud_cfw_ips_mode_switch" "example" {
  mode = 1
}
```

Import

CFW ips mode switch can be imported using the customId(like uuid or base64 string), e.g.

```
terraform import tencentcloud_cfw_ips_mode_switch.example FTNxVFqU1BeA5JKfQlmkPg==
```
