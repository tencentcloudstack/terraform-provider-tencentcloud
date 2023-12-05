Provides a resource to create a dnspod domain_lock

Example Usage

```hcl
resource "tencentcloud_dnspod_domain_lock" "domain_lock" {
  domain = "dnspod.cn"
  lock_days = 30
}
```