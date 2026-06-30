Provides a resource to manage DNSPod package domain binding

Example Usage

```hcl
resource "tencentcloud_dnspod_package_domain" "package_domain" {
  resource_id = "res-xxxxx"
  domain_id   = 12345
}
```

Import

dnspod package_domain can be imported using the id, e.g.

```
terraform import tencentcloud_dnspod_package_domain.package_domain resource_id#domain_id
```