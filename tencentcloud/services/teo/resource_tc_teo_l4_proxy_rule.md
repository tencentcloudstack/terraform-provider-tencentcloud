Provides a resource to create a teo teo_l4_proxy_rule

Example Usage

```hcl
resource "tencentcloud_teo_l4_proxy_rule" "teo_l4_proxy_rule" {
    proxy_id = "sid-38hbn26osico"
    zone_id  = "zone-36bjhygh1bxe"

    l4_proxy_rules {
        client_ip_pass_through_mode = "OFF"
        origin_port_range           = "1212"
        origin_type                 = "IP_DOMAIN"
        origin_value                = [
            "www.aaa.com",
        ]
        port_range                  = [
            "1212",
        ]
        protocol                    = "TCP"
        rule_tag                    = "aaa"
        session_persist             = "off"
        session_persist_time        = 3600
    }
}
```

Import

teo teo_l4_proxy can be imported using the id, e.g.

```
terraform import tencentcloud_teo_l4_proxy_rule.teo_l4_proxy_rule zoneId#proxyId#ruleId
```
