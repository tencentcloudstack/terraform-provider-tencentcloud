Provides a resource to create a dnspod modify_domain_owner

Example Usage

```hcl
resource "tencentcloud_dnspod_modify_domain_owner_operation" "modify_domain_owner" {
  domain = "dnspod.cn"
  account = "xxxxxxxxx"
  domain_id = 123
}
```