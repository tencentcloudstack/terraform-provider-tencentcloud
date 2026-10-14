---
subcategory: "Anti-DDoS(DayuV2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dayu_l4_rule"
sidebar_current: "docs-tencentcloud-resource-dayu_l4_rule"
description: |-
  Use this resource to create dayu layer 4 rule
---

# tencentcloud_dayu_l4_rule

Use this resource to create dayu layer 4 rule

~> **NOTE:** This resource only support resource Anti-DDoS of type `bgpip` and `net`

## Example Usage

```hcl
resource "tencentcloud_dayu_l4_rule" "test_rule" {
  resource_type             = "bgpip"
  resource_id               = "bgpip-00000294"
  name                      = "rule_test"
  protocol                  = "TCP"
  s_port                    = 80
  d_port                    = 60
  source_type               = 2
  health_check_switch       = true
  health_check_timeout      = 30
  health_check_interval     = 35
  health_check_health_num   = 5
  health_check_unhealth_num = 10
  session_switch            = false
  session_time              = 300

  source_list {
    source = "1.1.1.1"
    weight = 100
  }
  source_list {
    source = "2.2.2.2"
    weight = 50
  }
}
```

## Argument Reference

The following arguments are supported:

* `d_port` - (Required, Int) The destination port of the L4 rule.
* `name` - (Required, String, ForceNew) Name of the rule. When the `resource_type` is `net`, this field should be set with valid domain.
* `protocol` - (Required, String) Protocol of the rule. Valid values: `http`, `https`. When `source_type` is 1(host source), the value of this field can only set with `tcp`.
* `resource_id` - (Required, String, ForceNew) ID of the resource that the layer 4 rule works for.
* `resource_type` - (Required, String, ForceNew) Type of the resource that the layer 4 rule works for. Valid values: `bgpip` and `net`.
* `s_port` - (Required, Int) The source port of the L4 rule.
* `source_list` - (Required, Set) Source list of the rule, it can be a set of ip sources or a set of domain sources. The number of items ranges from 1 to 20.
* `source_type` - (Required, Int, ForceNew) Source type, `1` for source of host, `2` for source of IP.
* `health_check_health_num` - (Optional, Int) Health threshold of health check, and the default is 3. If a success result is returned for the health check 3 consecutive times, indicates that the forwarding is normal. The value range is 2-10.
* `health_check_interval` - (Optional, Int) Interval time of health check. The value range is 10-60 sec, and the default is 15 sec.
* `health_check_switch` - (Optional, Bool) Indicates whether health check is enabled. The default is `false`. Only valid when source list has more than one source item.
* `health_check_timeout` - (Optional, Int) HTTP Status Code. The default is 26 and value range is 2-60.
* `health_check_unhealth_num` - (Optional, Int) Unhealthy threshold of health check, and the default is 3. If the unhealthy result is returned 3 consecutive times, indicates that the forwarding is abnormal. The value range is 2-10.
* `session_switch` - (Optional, Bool) Indicate that the session will keep or not, and default value is `false`.
* `session_time` - (Optional, Int) Session keep time, only valid when `session_switch` is true, the available value ranges from 1 to 300 and unit is second.

The `source_list` object supports the following:

* `source` - (Required, String) Source IP or domain, valid format of ip is like `1.1.1.1` and valid format of host source is like `abc.com`.
* `weight` - (Required, Int) Weight of the source, the valid value ranges from 0 to 100.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `lb_type` - LB type of the rule. Valid values: `1`, `2`. `1` for weight cycling and `2` for IP hash.
* `rule_id` - ID of the layer 4 rule.


