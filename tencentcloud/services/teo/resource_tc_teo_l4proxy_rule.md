Provides a resource to create a teo l4proxy_rule

Example Usage

```hcl
resource "tencentcloud_teo_l4proxy_rule" "l4proxy_rule" {
  zone_id = "zone-21xfqlh4qjee"
  proxy_id = "proxy-00dde483-9a2b-11ec-bcb0"
  l4_proxy_rules {
		rule_id = "rule-2qzkbvx3hxl7"
		protocol = "TCP"
		port_range = &lt;nil&gt;
		origin_type = "IP_DOMAIN"
		origin_value = &lt;nil&gt;
		origin_port_range = "80-90"
		client_ip_pass_through_mode = "TOA"
		session_persist = "off"
		session_persist_time = 300
		rule_tag = "rule-service1	"
		status = "offline"

  }
}
```

Import

teo l4proxy_rule can be imported using the id, e.g.

```
terraform import tencentcloud_teo_l4proxy_rule.l4proxy_rule l4proxy_rule_id
```
