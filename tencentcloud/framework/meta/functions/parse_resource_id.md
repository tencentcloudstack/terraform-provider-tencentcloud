Provides a provider-defined function that parses a TencentCloud composite
resource ID into its constituent fields. The function is pure (no cloud
API call) and is safe to use in any Terraform expression context that
allows function calls.

Example Usage

```hcl
locals {
  parts = provider::tencentcloud::parse_resource_id("ins-abcd1234#vpc-xyz0987")
}

output "instance_id" {
  value = local.parts.instance_id
}
```
