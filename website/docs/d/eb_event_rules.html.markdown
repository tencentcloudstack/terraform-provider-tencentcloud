---
subcategory: "EventBridge(EB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eb_event_rules"
sidebar_current: "docs-tencentcloud-datasource-eb_event_rules"
description: |-
  Use this data source to query detailed information of eb event_rules
---

# tencentcloud_eb_event_rules

Use this data source to query detailed information of eb event_rules

## Example Usage

```hcl
resource "tencentcloud_eb_event_bus" "foo" {
  event_bus_name = "tf-event_bus_rule"
  description    = "event bus desc"
  enable_store   = false
  save_days      = 1
  tags = {
    "createdBy" = "terraform"
  }
}
resource "tencentcloud_eb_event_rule" "event_rule" {
  event_bus_id = tencentcloud_eb_event_bus.foo.id
  rule_name    = "tf-event_rule"
  description  = "event rule desc"
  enable       = true
  event_pattern = jsonencode(
    {
      source = "apigw.cloud.tencent"
      type = [
        "connector:apigw",
      ]
    }
  )
  tags = {
    "createdBy" = "terraform"
  }
}
data "tencentcloud_eb_event_rules" "event_rules" {
  event_bus_id = tencentcloud_eb_event_bus.foo.id
  order_by     = "AddTime"
  order        = "DESC"
  depends_on   = [tencentcloud_eb_event_rule.event_rule]
}
```

## Argument Reference

The following arguments are supported:

* `event_bus_id` - (Required, String) event bus Id.
* `order_by` - (Optional, String) According to which field to sort the returned results, the following fields are supported: AddTime (creation time), ModTime (modification time).
* `order` - (Optional, String) Return results in ascending or descending order, optional values ASC (ascending) and DESC (descending).
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `rules` - Event rule information.
  * `add_time` - create time.
  * `dead_letter_config` - The dlq rule set by rule. It may be null. Note: this field may return null, indicating that no valid value can be obtained.
    * `ckafka_delivery_params` - After setting the DLQ mode, this option is required. The error message will be delivered to the corresponding kafka topic Note: This field may return null, indicating that no valid value can be obtained.
      * `resource_description` - ckafka resource qcs six-segment.
      * `topic_name` - ckafka topic name.
    * `dispose_method` - Support three modes of dlq, discarding, ignoring errors and continuing to pass, corresponding to: DLQ, DROP, IGNORE_ERROR.
  * `description` - description.
  * `enable` - enable switch.
  * `event_bus_id` - event bus Id.
  * `mod_time` - modify time.
  * `rule_id` - rule Id.
  * `rule_name` - rule name.
  * `status` - Status.
  * `targets` - Target brief information, note: this field may return null, indicating that no valid value can be obtained.
    * `target_id` - target Id.
    * `type` - target type.


