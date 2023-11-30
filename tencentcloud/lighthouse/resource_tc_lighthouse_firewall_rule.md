Provides a resource to create a lighthouse firewall rule

~> **NOTE:**  Use an empty template to clean up the default rules before using this resource manage firewall rules.

Example Usage

```hcl
resource "tencentcloud_lighthouse_firewall_rule" "firewall_rule" {
  instance_id = "lhins-xxxxxxx"
  firewall_rules {
	protocol = "TCP"
	port = "80"
	cidr_block = "10.0.0.1"
	action = "ACCEPT"
	firewall_rule_description = "description 1"
  }
  firewall_rules {
	protocol = "TCP"
	port = "80"
	cidr_block = "10.0.0.2"
	action = "ACCEPT"
	firewall_rule_description = "description 2"
  }
}
```

Import

lighthouse firewall_rule can be imported using the id, e.g.

```
terraform import tencentcloud_lighthouse_firewall_rule.firewall_rule lighthouse_instance_id
```