---
subcategory: "Anti-DDoS(DayuV2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dayu_l4_rule_v2"
sidebar_current: "docs-tencentcloud-resource-dayu_l4_rule_v2"
description: |-
  Use this resource to create dayu layer 4 rule
---

# tencentcloud_dayu_l4_rule_v2

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

* `business` - (Required, String, ForceNew) Business of the resource that the layer 4 rule works for. Valid values: `bgpip` and `net`.
* `resource_id` - (Required, String, ForceNew) Resource id.
* `rules` - (Required, List, ForceNew) A list of layer 4 rules. Each element contains the following attributes:
* `virtual_port` - (Required, Int, ForceNew) The virtual port of the layer 4 rule.
* `vpn` - (Required, String, ForceNew) Resource vpn.

The `rules` object supports the following:

* `keep_enable` - (Required, Bool) session hold switch.
* `keeptime` - (Required, Int) The keeptime of the layer 4 rule.
* `lb_type` - (Required, Int) LB type of the rule, `1` for weight cycling and `2` for IP hash.
* `protocol` - (Required, String) Protocol of the rule.
* `region` - (Required, Int) Corresponding regional information.
* `remove_switch` - (Required, Bool) Remove the watermark state.
* `rule_name` - (Required, String) Name of the rule.
* `source_list` - (Required, List) Source list of the rule.
* `source_port` - (Required, Int) The source port of the layer 4 rule.
* `source_type` - (Required, Int) Source type, `1` for source of host, `2` for source of IP.
* `virtual_port` - (Required, Int) The virtual port of the layer 4 rule.

The `source_list` object of `rules` supports the following:

* `source` - (Required, String) Source IP or domain.
* `weight` - (Required, Int) Weight of the source.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateNewL4Rules](https://www.tencentcloud.com/search/?wd=CreateNewL4Rules&type=document&lang=en)
* Read: [DescribeNewL4Rules](https://www.tencentcloud.com/search/?wd=DescribeNewL4Rules&type=document&lang=en)
* Delete:
  - [DeleteNewL4Rules](https://www.tencentcloud.com/search/?wd=DeleteNewL4Rules&type=document&lang=en)
  - [DescribeNewL4Rules](https://www.tencentcloud.com/search/?wd=DescribeNewL4Rules&type=document&lang=en)

