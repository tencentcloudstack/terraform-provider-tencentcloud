Provides a resource to create a teo teo_l4_proxy

Example Usage

```hcl
resource "tencentcloud_teo_l4_proxy" "proxy" {
  accelerate_mainland = "off"
  area                = "overseas"
  ipv6                = "off"
  proxy_name          = "proxy-test"
  static_ip           = "off"
  zone_id             = "zone-2qtuhspy6cr7"
}
```

Import

teo teo_l4_proxy can be imported using the id, e.g.

```
terraform import tencentcloud_teo_l4_proxy.teo_l4_proxy teo_l4_proxy_id
```
