Provides a resource to create a teo l4proxy

Example Usage

```hcl
resource "tencentcloud_teo_l4proxy" "l4proxy" {
  zone_id = "zone-21xfqlh4qjee"
  proxy_name = "test-proxy"
  area = "mainland"
  ipv6 = "on"
  static_ip = "on"
  accelerate_mainland = "off"
  d_dos_protection_config {
		level_mainland = "BASE30_MAX300"
		max_bandwidth_mainland = 
		level_overseas = ""

  }
}
```

Import

teo l4proxy can be imported using the id, e.g.

```
terraform import tencentcloud_teo_l4proxy.l4proxy l4proxy_id
```
