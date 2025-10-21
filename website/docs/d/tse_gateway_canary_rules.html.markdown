---
subcategory: "Tencent Cloud Service Engine(TSE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tse_gateway_canary_rules"
sidebar_current: "docs-tencentcloud-datasource-tse_gateway_canary_rules"
description: |-
  Use this data source to query detailed information of tse gateway_canary_rules
---

# tencentcloud_tse_gateway_canary_rules

Use this data source to query detailed information of tse gateway_canary_rules

## Example Usage

```hcl
data "tencentcloud_tse_gateway_canary_rules" "gateway_canary_rules" {
  gateway_id = "gateway-xxxxxx"
  service_id = "451a9920-e67a-4519-af41-fccac0e72005"
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, String) gateway ID.
* `service_id` - (Required, String) service ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - canary rule configuration.
  * `canary_rule_list` - canary rule list.
    * `balanced_service_list` - service weight configuration.
      * `percent` - percent, 10 is 10%, valid values: 0 to 100.
      * `service_id` - service ID.
      * `service_name` - service name.
      * `upstream_name` - upstream name.
    * `condition_list` - parameter matching condition list.
      * `delimiter` - delimiter. valid when operator is in or not in, reference value:`,`, `;`,`\n`.
      * `global_config_id` - global configuration ID.
      * `global_config_name` - global configuration name.
      * `key` - parameter name.
      * `operator` - operator.Reference value:`le`, `eq`, `lt`, `ne`, `ge`, `gt`, `regex`, `exists`, `in`, `not in`,  `prefix`, `exact`, `regex`.
      * `type` - type.Reference value:- path- method- query- header- cookie- body- system.
      * `value` - parameter value.
    * `enabled` - the status of canary rule.
    * `priority` - priority. The value ranges from 0 to 100; the larger the value, the higher the priority; the priority cannot be repeated between different rules.
    * `service_id` - service ID.
    * `service_name` - service name.
  * `total_count` - total count.


