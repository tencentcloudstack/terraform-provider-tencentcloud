---
subcategory: "TencentCloud Lighthouse(Lighthouse)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lighthouse_firewall_template"
sidebar_current: "docs-tencentcloud-resource-lighthouse_firewall_template"
description: |-
  Provides a resource to create a lighthouse firewall template
---

# tencentcloud_lighthouse_firewall_template

Provides a resource to create a lighthouse firewall template

## Example Usage

```hcl
resource "tencentcloud_lighthouse_firewall_template" "firewall_template" {
  template_name = "firewall-template-test"
  template_rules {
    protocol                  = "TCP"
    port                      = "8080"
    cidr_block                = "127.0.0.1"
    action                    = "ACCEPT"
    firewall_rule_description = "test description"
  }
  template_rules {
    protocol                  = "TCP"
    port                      = "8090"
    cidr_block                = "127.0.0.0/24"
    action                    = "DROP"
    firewall_rule_description = "test description"
  }
}
```

## Argument Reference

The following arguments are supported:

* `template_name` - (Required, String) Template name.
* `template_rules` - (Optional, List) List of firewall rules.

The `template_rules` object supports the following:

* `protocol` - (Required, String) Protocol. Values: TCP, UDP, ICMP, ALL.
* `action` - (Optional, String) Action. Values: ACCEPT, DROP. The default is `ACCEPT`.
* `cidr_block` - (Optional, String) Network segment or IP (mutually exclusive). The default is `0.0.0.0`, indicating all sources.
* `firewall_rule_description` - (Optional, String) Firewall rule description.
* `port` - (Optional, String) Port. Values: ALL, Separate ports, comma-separated discrete ports, minus sign-separated port ranges.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateFirewallTemplate](https://www.tencentcloud.com/search/?wd=CreateFirewallTemplate&type=document&lang=en)
* Read: [DescribeFirewallTemplateRules](https://www.tencentcloud.com/search/?wd=DescribeFirewallTemplateRules&type=document&lang=en)
* Read: [DescribeFirewallTemplates](https://www.tencentcloud.com/search/?wd=DescribeFirewallTemplates&type=document&lang=en)
* Update: [ModifyFirewallTemplate](https://www.tencentcloud.com/search/?wd=ModifyFirewallTemplate&type=document&lang=en)
* Delete: [DeleteFirewallTemplate](https://www.tencentcloud.com/search/?wd=DeleteFirewallTemplate&type=document&lang=en)

## Import

lighthouse firewall_template can be imported using the id, e.g.

```
terraform import tencentcloud_lighthouse_firewall_template.firewall_template firewall_template_id
```

