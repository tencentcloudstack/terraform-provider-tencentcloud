Provides a resource to create a dnspod domain_alias

Example Usage

```hcl
resource "tencentcloud_dnspod_domain_alias" "domain_alias" {
  domain_alias = "dnspod.com"
  domain = "dnspod.cn"
}
```

Import

dnspod domain_alias can be imported using the id, e.g.

```
terraform import tencentcloud_dnspod_domain_alias.domain_alias domain#domain_alias_id
```