---
subcategory: "TencentCloud Lighthouse(Lighthouse)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lighthouse_firewall_rule"
sidebar_current: "docs-tencentcloud-resource-lighthouse_firewall_rule"
description: |-
  Provides a resource to create a lighthouse firewall rule
---

# tencentcloud_lighthouse_firewall_rule

Provides a resource to create a lighthouse firewall rule

~> **NOTE:**  Use an empty template to clean up the default rules before using this resource manage firewall rules.

## Example Usage

```hcl
resource "tencentcloud_lighthouse_firewall_rule" "firewall_rule" {
  instance_id = "lhins-xxxxxxx"
  firewall_rules {
    protocol                  = "TCP"
    port                      = "80"
    cidr_block                = "10.0.0.1"
    action                    = "ACCEPT"
    firewall_rule_description = "description 1"
  }
  firewall_rules {
    protocol                  = "TCP"
    port                      = "80"
    cidr_block                = "10.0.0.2"
    action                    = "ACCEPT"
    firewall_rule_description = "description 2"
  }
}
```

## Argument Reference

The following arguments are supported:

* `firewall_rules` - (Required, List) Firewall rule list.
* `instance_id` - (Required, String, ForceNew) Instance ID.

The `firewall_rules` object supports the following:

* `protocol` - (Required, String) Protocol. Valid values are TCP, UDP, ICMP, ALL.
* `action` - (Optional, String) Valid values are ACCEPT, DROP. Default value is ACCEPT.
* `cidr_block` - (Optional, String) IP range or IP (mutually exclusive). Default value is 0.0.0.0/0, which indicates all sources.
* `firewall_rule_description` - (Optional, String) Firewall rule description.
* `port` - (Optional, String) Port. Valid values are ALL, one single port, multiple ports separated by commas, or port range indicated by a minus sign.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

lighthouse firewall_rule can be imported using the id, e.g.

```
terraform import tencentcloud_lighthouse_firewall_rule.firewall_rule lighthouse_instance_id
```

