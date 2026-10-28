Provides a resource to manage DNSPod package domain binding

Example Usage

```hcl
resource "tencentcloud_dnspod_package_domain" "example" {
  resource_id = "91d8006a"
  domain_id   = 92435817
}
```

Import

dnspod package_domain can be imported using the id, e.g.

```
terraform import tencentcloud_dnspod_package_domain.example 91d8006a
```
