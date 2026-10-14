---
subcategory: "Anti-DDoS(Dayu)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dayu_l7_rules"
sidebar_current: "docs-tencentcloud-datasource-dayu_l7_rules"
description: |-
  Use this data source to query dayu layer 7 rules
---

# tencentcloud_dayu_l7_rules

Use this data source to query dayu layer 7 rules

## Example Usage

```hcl
data "tencentcloud_dayu_l7_rules" "domain_test" {
  resource_type = tencentcloud_dayu_l7_rule.test_rule.resource_type
  resource_id   = tencentcloud_dayu_l7_rule.test_rule.resource_id
  domain        = tencentcloud_dayu_l7_rule.test_rule.domain
}
data "tencentcloud_dayu_l7_rules" "id_test" {
  resource_type = tencentcloud_dayu_l7_rule.test_rule.resource_type
  resource_id   = tencentcloud_dayu_l7_rule.test_rule.resource_id
  rule_id       = tencentcloud_dayu_l7_rule.test_rule.rule_id
}
```

## Argument Reference

The following arguments are supported:

* `resource_id` - (Required, String) Id of the resource that the layer 7 rule works for.
* `resource_type` - (Required, String) Type of the resource that the layer 7 rule works for, valid value is `bgpip`.
* `domain` - (Optional, String) Domain of the layer 7 rule to be queried.
* `result_output_file` - (Optional, String) Used to save results.
* `rule_id` - (Optional, String) Id of the layer 7 rule to be queried.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of layer 7 rules. Each element contains the following attributes:
  * `domain` - Domain that the 7 layer rule works for.
  * `health_check_code` - HTTP Status Code. `1` means the return value `1xx` is health. `2` means the return value `2xx` is health. `4` means the return value `3xx` is health. `8` means the return value `4xx` is health. `16` means the return value `5xx` is health. If you want multiple return codes to indicate health, need to add the corresponding values.
  * `health_check_health_num` - Health threshold of health check.
  * `health_check_interval` - Interval time of health check.
  * `health_check_method` - Methods of health check.
  * `health_check_path` - Path of health check.
  * `health_check_switch` - Indicates whether health check is enabled.
  * `health_check_unhealth_num` - Unhealthy threshold of health check.
  * `name` - Name of the rule.
  * `protocol` - Protocol of the rule.
  * `rule_id` - Id of the 7 layer rule.
  * `source_list` - Source list of the rule.
  * `source_type` - Source type, 1 for source of host, 2 for source of ip.
  * `ssl_id` - SSL id.
  * `status` - Status of the rule. `0` for create/modify success, `2` for create/modify fail, `3` for delete success, `5` for waiting to be created/modified, `7` for waiting to be deleted and `8` for waiting to get SSL id.
  * `switch` - Indicate the rule will take effect or not.
  * `threshold` - Threshold of the rule.


