Provides a resource to create a teo application_proxy_rule

Example Usage

```hcl
resource "tencentcloud_teo_application_proxy_rule" "application_proxy_rule" {
  forward_client_ip = "TOA"
  origin_type       = "custom"
  origin_port       = "8083"
  origin_value      = [
    "127.0.0.1",
  ]
  port              = [
    "8083",
  ]
  proto             = "TCP"
  proxy_id          = "proxy-6972528a-373a-11ed-afca-52540044a456"
  session_persist   = false
  status            = "online"
  zone_id           = "zone-2983wizgxqvm"
}

```
Import

teo application_proxy_rule can be imported using the zoneId#proxyId#ruleId, e.g.
```
terraform import tencentcloud_teo_application_proxy_rule.application_proxy_rule zone-2983wizgxqvm#proxy-6972528a-373a-11ed-afca-52540044a456#rule-90b13bb4-373a-11ed-8794-525400eddfed
```