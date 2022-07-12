---
subcategory: "Anti-DDoS(Dayu)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dayu_l7_rule"
sidebar_current: "docs-tencentcloud-resource-dayu_l7_rule"
description: |-
  Use this resource to create dayu layer 7 rule
---

# tencentcloud_dayu_l7_rule

Use this resource to create dayu layer 7 rule

~> **NOTE:** This resource only support resource Anti-DDoS of type `bgpip`

## Example Usage

```hcl
resource "tencentcloud_dayu_l7_rule" "test_rule" {
  resource_type             = "bgpip"
  resource_id               = "bgpip-00000294"
  name                      = "rule_test"
  domain                    = "zhaoshaona.com"
  protocol                  = "https"
  switch                    = true
  source_type               = 2
  source_list               = ["1.1.1.1:80", "2.2.2.2"]
  ssl_id                    = "%s"
  health_check_switch       = true
  health_check_code         = 31
  health_check_interval     = 30
  health_check_method       = "GET"
  health_check_path         = "/"
  health_check_health_num   = 5
  health_check_unhealth_num = 10
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) Domain that the layer 7 rule works for. Valid string length ranges from 0 to 80.
* `name` - (Required, String, ForceNew) Name of the rule.
* `protocol` - (Required, String) Protocol of the rule. Valid values: `http`, `https`.
* `resource_id` - (Required, String, ForceNew) ID of the resource that the layer 7 rule works for.
* `resource_type` - (Required, String, ForceNew) Type of the resource that the layer 7 rule works for, valid value is `bgpip`.
* `source_list` - (Required, Set: [`String`]) Source list of the rule, it can be a set of ip sources or a set of domain sources. The number of items ranges from 1 to 16.
* `source_type` - (Required, Int) Source type, `1` for source of host, `2` for source of IP.
* `switch` - (Required, Bool) Indicate the rule will take effect or not.
* `health_check_code` - (Optional, Int) HTTP Status Code. The default is `26`. Valid value ranges: [1~31]. `1` means the return value '1xx' is health. `2` means the return value '2xx' is health. `4` means the return value '3xx' is health. `8` means the return value '4xx' is health. `16` means the return value '5xx' is health. If you want multiple return codes to indicate health, need to add the corresponding values.
* `health_check_health_num` - (Optional, Int) Health threshold of health check, and the default is `3`. If a success result is returned for the health check 3 consecutive times, indicates that the forwarding is normal. The value range is [2-10].
* `health_check_interval` - (Optional, Int) Interval time of health check. Valid value ranges: [10~60]sec. The default is 15 sec.
* `health_check_method` - (Optional, String) Methods of health check. The default is 'HEAD', the available value are 'HEAD' and 'GET'.
* `health_check_path` - (Optional, String) Path of health check. The default is `/`.
* `health_check_switch` - (Optional, Bool) Indicates whether health check is enabled. The default is `false`.
* `health_check_unhealth_num` - (Optional, Int) Unhealthy threshold of health check, and the default is `3`. If the unhealthy result is returned 3 consecutive times, indicates that the forwarding is abnormal. The value range is [2-10].
* `ssl_id` - (Optional, String) SSL ID, when the `protocol` is `https`, the field should be set with valid SSL id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `rule_id` - ID of the layer 7 rule.
* `status` - Status of the rule. `0` for create/modify success, `2` for create/modify fail, `3` for delete success, `5` for delete failed, `6` for waiting to be created/modified, `7` for waiting to be deleted and 8 for waiting to get SSL ID.


