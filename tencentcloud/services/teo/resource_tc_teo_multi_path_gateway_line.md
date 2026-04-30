Provides a resource to create a teo multi path gateway line for EdgeOne(TEO).

Example Usage

Custom line type

```hcl
resource "tencentcloud_teo_multi_path_gateway_line" "custom" {
  zone_id      = "zone-359h792djt7h"
  gateway_id   = "mpgw-g3176ppeye"
  line_type    = "custom"
  line_address = "1.2.3.4:81"
}
```

Proxy line type

```hcl
resource "tencentcloud_teo_multi_path_gateway_line" "proxy" {
    gateway_id   = "mpgw-g3176ppeye"
    line_address = "tf-test.359h792djt7h.eo.dnse0.com:82"
    line_type    = "proxy"
    proxy_id     = "sid-3phb7c7we1ns"
    zone_id      = "zone-359h792djt7h"
    rule_id      = "rule-3picr1x9wa0u"
}
```

Import

teo multi path gateway line can be imported using the id, e.g.

```
terraform import tencentcloud_teo_multi_path_gateway_line.example zone-279qso5a4cw9#gw-2qwk1t3g3jx9#line-1
```
