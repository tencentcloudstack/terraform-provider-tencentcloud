Provides a resource to create a DNSPod package order

Example Usage

```hcl
resource "tencentcloud_dnspod_package_order" "example" {
  domain = "demo.com"
  grade  = "DPG_ULTIMATE"
}
```

Import

DNSPod package order can be imported using the id, e.g.

```
terraform import tencentcloud_dnspod_package_order.example demo.com
```
