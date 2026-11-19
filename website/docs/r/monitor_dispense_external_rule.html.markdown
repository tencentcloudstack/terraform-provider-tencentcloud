---
subcategory: "Cloud Monitor(Monitor)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_dispense_external_rule"
sidebar_current: "docs-tencentcloud-resource-monitor_dispense_external_rule"
description: |-
  Provides a resource to create a monitor dispense external rule (dispense rule 3.0) which forwards cloud monitor external metrics to Kafka.
---

# tencentcloud_monitor_dispense_external_rule

Provides a resource to create a monitor dispense external rule (dispense rule 3.0) which forwards cloud monitor external metrics to Kafka.

## Example Usage

```hcl
resource "tencentcloud_monitor_dispense_external_rule" "example" {
  name             = "tf-example"
  ext_namespace    = "QCE/CVM"
  dispense_regions = ["gz"]
  status           = 1
  ext_metrics = [
    "AccOuttraffic",
    "Acs",
    "ActiveMtu",
    "AtsEnabled",
    "BaseCpuUsage",
    "BondingMode",
    "Cnp802pPrio",
    "CbsVolumeFsUsage",
  ]
  period = [60]
  producer = {
    protocol_type = 1
    type          = "Kafka"
    brokers       = "13.132.191.101:8080"
    topic         = "topic-demo"
  }

  dispense_conditions = {
    ext_metric = "AtsEnabled"
    dispense_filters = [
      {
        key        = "InstanceId"
        expression = "="
        values     = ["ins-p97oyxt7"]
      },
      {
        key        = "InstanceId"
        expression = "in"
        values     = ["ins-p97oyxt7,ins-ly4ibb5w"]
      },
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `ext_namespace` - (Required, String) Cloud monitor external namespace.
* `name` - (Required, String) Rule name.
* `producer` - (Required, Object) Dispense target producer information.
* `dispense_conditions` - (Optional, Object) Dispense filter condition information.
* `dispense_regions` - (Optional, Set) Dispense deployment region list.
* `ext_metrics` - (Optional, Set) Cloud monitor external metric name list. The describe API reports metric details through `rule.ext_metric` only; the name-to-detail mapping is not confirmed, so this attribute is not refreshed from the response and drift on it is not detected.
* `period` - (Optional, Set) Metric statistical period list.
* `status` - (Optional, Int) Rule trigger status applied through `ModifyDispenseExternalRuleStatus` after the rule is created or updated. `0` means disabled, `1` means enabled. The effective status is always reported by `rule.status`.

The `dispense_conditions` object supports the following:

* `condition_id` - (Optional, Int) Filter condition id.
* `dispense_filters` - (Optional, List) Filter condition list.
* `ext_metric` - (Optional, String) External metric name.

The `dispense_filters` object of `dispense_conditions` supports the following:

* `expression` - (Optional, String) Expression.
* `key` - (Optional, String) Dimension name.
* `values` - (Optional, Set) Dimension value list.

The `producer` object supports the following:

* `brokers` - (Required, String) Kafka broker address to dispense to.
* `protocol_type` - (Required, Int) Dispense protocol type. Valid values: `0` stormRetPb, `1` tcbDispensePb, `2` stormRetJson, `3` ADPPb (deprecated), `4` zhongtai pb.
* `topic` - (Required, String) Kafka topic to dispense to.
* `type` - (Required, String) Target type.
* `default_tags` - (Optional, Set) Default dimension groups. Only dimensions are provided.
* `global_tags` - (Optional, List) Global dimension groups.
* `merge` - (Optional, Int) Whether to merge metrics. Default value: `1` (merged).
* `password` - (Optional, String) Kafka password.
* `username` - (Optional, String) Kafka user name.

The `global_tags` object of `producer` supports the following:

* `key` - (Optional, String) Dimension key.
* `value` - (Optional, String) Dimension value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

      * `expression` - Expression.
      * `key` - Dimension key.
      * `key` - Dimension name.
      * `value` - Dimension value.
      * `values` - Dimension value list.
    * `brokers` - Kafka broker address to dispense to.
    * `cn_meaning` - Chinese meaning of the metric.
    * `condition_id` - Filter condition id.
    * `default_tags` - Default dimension groups. Only dimensions are provided.
    * `dimension_flag` - Whether external dimensions are configured for the metric.
    * `dispense_filters` - Filter condition list.
    * `en_meaning` - English meaning of the metric.
    * `ext_metric` - External metric name.
    * `global_tags` - Global dimension groups.
    * `merge` - Whether to merge metrics. Default value: `1` (merged).
    * `metric_c_name` - Chinese metric name.
    * `metric_name` - Metric name.
    * `password` - Kafka password.
    * `protocol_type` - Dispense protocol type. Valid values: `0` stormRetPb, `1` tcbDispensePb, `2` stormRetJson, `3` ADPPb (deprecated), `4` zhongtai pb.
    * `topic` - Kafka topic to dispense to.
    * `type` - Target type.
    * `unit` - Metric unit.
    * `username` - Kafka user name.
  * `dispense_conditions` - Dispense filter conditions.
  * `dispense_regions` - Dispense region list.
  * `ext_metric` - External metric list.
  * `ext_namespace` - External namespace.
  * `name` - Rule name.
  * `period` - Metric statistical period list.
  * `producer` - Producer information.
  * `rule_id` - Rule id.
  * `status` - Rule trigger status.
  * `update_time` - Update time.
* `id` - Id of the resource, which is the string form of `rule_id`.
* `rule_id` - Dispense rule id assigned by the cloud side on creation.
* `rule` - Rule detail reported by the describe API.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateDispenseExternalRule](https://www.tencentcloud.com/search/?wd=CreateDispenseExternalRule&type=document&lang=en)
  - [DescribeDispenseExternalRule](https://www.tencentcloud.com/search/?wd=DescribeDispenseExternalRule&type=document&lang=en)
  - [ModifyDispenseExternalRuleStatus](https://www.tencentcloud.com/search/?wd=ModifyDispenseExternalRuleStatus&type=document&lang=en)
* Read: [DescribeDispenseExternalRule](https://www.tencentcloud.com/search/?wd=DescribeDispenseExternalRule&type=document&lang=en)
* Update:
  - [DescribeDispenseExternalRule](https://www.tencentcloud.com/search/?wd=DescribeDispenseExternalRule&type=document&lang=en)
  - [ModifyDispenseExternalRule](https://www.tencentcloud.com/search/?wd=ModifyDispenseExternalRule&type=document&lang=en)
  - [ModifyDispenseExternalRuleStatus](https://www.tencentcloud.com/search/?wd=ModifyDispenseExternalRuleStatus&type=document&lang=en)
* Delete: [DeleteDispenseExternalRule](https://www.tencentcloud.com/search/?wd=DeleteDispenseExternalRule&type=document&lang=en)

## Import

Monitor dispense external rule can be imported using the rule id, e.g.

```
terraform import tencentcloud_monitor_dispense_external_rule.example 2918
```

