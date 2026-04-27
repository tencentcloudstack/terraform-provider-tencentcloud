Provides a resource to create a teo multi path gateway line for EdgeOne(TEO).

Example Usage

Custom line type

```hcl
resource "tencentcloud_teo_multi_path_gateway_line" "example" {
  zone_id      = "zone-279qso5a4cw9"
  gateway_id   = "gw-2qwk1t3g3jx9"
  line_type    = "custom"
  line_address = "1.2.3.4:80"
}
```

Proxy line type

```hcl
resource "tencentcloud_teo_multi_path_gateway_line" "example" {
  zone_id      = "zone-279qso5a4cw9"
  gateway_id   = "gw-2qwk1t3g3jx9"
  line_type    = "proxy"
  line_address = "5.6.7.8:443"
  proxy_id     = "sid-38hbn26osico"
  rule_id      = "rule-abcdef"
}
```

Import

teo multi path gateway line can be imported using the id, e.g.

```
terraform import tencentcloud_teo_multi_path_gateway_line.example zone-279qso5a4cw9#gw-2qwk1t3g3jx9#line-1
```
