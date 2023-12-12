Provides a resource to create a lighthouse firewall template

Example Usage

```hcl
resource "tencentcloud_lighthouse_firewall_template" "firewall_template" {
  template_name = "firewall-template-test"
  template_rules {
    protocol = "TCP"
    port = "8080"
    cidr_block = "127.0.0.1"
    action = "ACCEPT"
    firewall_rule_description = "test description"
  }
  template_rules {
    protocol = "TCP"
    port = "8090"
    cidr_block = "127.0.0.0/24"
    action = "DROP"
    firewall_rule_description = "test description"
  }
}
```

Import

lighthouse firewall_template can be imported using the id, e.g.

```
terraform import tencentcloud_lighthouse_firewall_template.firewall_template firewall_template_id
```