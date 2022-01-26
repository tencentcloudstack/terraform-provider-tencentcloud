---
subcategory: "Anti-DDoS(DayuV2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dayu_l4_rules_v2"
sidebar_current: "docs-tencentcloud-datasource-dayu_l4_rules_v2"
description: |-
  Use this data source to query dayu new layer 4 rules
---

# tencentcloud_dayu_l4_rules_v2

Use this data source to query dayu new layer 4 rules

## Example Usage

```hcl
data "tencentcloud_dayu_l4_rules_v2" "tencentcloud_dayu_l4_rules_v2" {
  business = "bgpip"
}
```

## Argument Reference

The following arguments are supported:

* `business` - (Required) Type of the resource that the layer 4 rule works for, valid values are `bgpip`, `bgp`, `bgp-multip` and `net`.
* `ip` - (Optional) Ip of the resource.
* `result_output_file` - (Optional) Used to save results.
* `virtual_port` - (Optional) Virtual port of resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of layer 4 rules. Each element contains the following attributes:
  * `id` - Bind the resource ID information.
  * `ip` - Bind the resource IP information.
  * `keep_enable` - session hold switch.
  * `keeptime` - The keeptime of the layer 4 rule.
  * `lb_type` - LB type of the rule, `1` for weight cycling and `2` for IP hash.
  * `modify_time` - Rule modification time.
  * `protocol` - Protocol of the rule.
  * `region` - Corresponding regional information.
  * `remove_switch` - Remove the watermark state.
  * `rule_id` - ID of the 4 layer rule.
  * `rule_name` - Name of the rule.
  * `source_port` - The source port of the layer 4 rule.
  * `source_type` - Source type, `1` for source of host, `2` for source of IP.
  * `virtual_port` - The virtual port of the layer 4 rule.


