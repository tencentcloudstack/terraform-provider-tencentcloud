---
subcategory: "Anti-DDoS(Dayu)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dayu_l4_rules"
sidebar_current: "docs-tencentcloud-datasource-dayu_l4_rules"
description: |-
  Use this data source to query dayu layer 4 rules
---

# tencentcloud_dayu_l4_rules

Use this data source to query dayu layer 4 rules

## Example Usage

```hcl
data "tencentcloud_dayu_l4_rules" "name_test" {
  resource_type = tencentcloud_dayu_l4_rule.test_rule.resource_type
  resource_id   = tencentcloud_dayu_l4_rule.test_rule.resource_id
  name          = tencentcloud_dayu_l4_rule.test_rule.name
}
data "tencentcloud_dayu_l4_rules" "id_test" {
  resource_type = tencentcloud_dayu_l4_rule.test_rule.resource_type
  resource_id   = tencentcloud_dayu_l4_rule.test_rule.resource_id
  rule_id       = tencentcloud_dayu_l4_rule.test_rule.rule_id
}
```

## Argument Reference

The following arguments are supported:

* `resource_id` - (Required) Id of the resource that the layer 4 rule works for.
* `resource_type` - (Required) Type of the resource that the layer 4 rule works for, valid values are `bgpip`, `bgp`, `bgp-multip` and `net`.
* `name` - (Optional) Name of the layer 4 rule to be queried.
* `result_output_file` - (Optional) Used to save results.
* `rule_id` - (Optional) Id of the layer 4 rule to be queried.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of layer 4 rules. Each element contains the following attributes:
  * `d_port` - The destination port of the layer 4 rule.
  * `health_check_health_num` - Health threshold of health check.
  * `health_check_interval` - Interval time of health check.
  * `health_check_switch` - Indicates whether health check is enabled.
  * `health_check_timeout` - HTTP Status Code. `1` means the return value `1xx` is health. `2` means the return value `2xx` is health. `4` means the return value `3xx` is health. `8` means the return value `4xx` is health. `16` means the return value `5xx` is health. If you want multiple return codes to indicate health, need to add the corresponding values.
  * `health_check_unhealth_num` - Unhealthy threshold of health check.
  * `lb_type` - LB type of the rule, `1` for weight cycling and `2` for IP hash.
  * `name` - Name of the rule.
  * `protocol` - Protocol of the rule.
  * `rule_id` - ID of the 4 layer rule.
  * `s_port` - The source port of the layer 4 rule.
  * `session_switch` - Indicate that the session will keep or not.
  * `session_time` - Session keep time, only valid when `session_switch` is true, the available value ranges from 1 to 300 and unit is second.
  * `source_type` - Source type, `1` for source of host, `2` for source of IP.


