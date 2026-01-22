---
subcategory: "Application Performance Management(APM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_apm_prometheus_rule"
sidebar_current: "docs-tencentcloud-resource-apm_prometheus_rule"
description: |-
  Provides a resource to create a APM prometheus rule
---

# tencentcloud_apm_prometheus_rule

Provides a resource to create a APM prometheus rule

## Example Usage

```hcl
resource "tencentcloud_apm_prometheus_rule" "example" {
  instance_id       = "apm-lhqHyRBuA"
  name              = "tf-example"
  service_name      = "java-market-service"
  metric_match_type = 0
  metric_name_rule  = "task.duration"
  status            = 1
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Business system ID.
* `metric_match_type` - (Required, Int) Match type: 0 - precision match, 1 - prefix match, 2 - suffix match.
* `metric_name_rule` - (Required, String) Specifies the rule for customer-defined metric names with cache hit.
* `name` - (Required, String) Metric match rule name.
* `service_name` - (Required, String) Applications where the rule takes effect. input an empty string for all applications.
* `status` - (Optional, Int) Rule status. 1 - enabled, 2 - disabled. Default value: 1.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `rule_id` - ID of the indicator matching rule.


## Import

APM prometheus rule can be imported using the instanceId#ruleId, e.g.

```
terraform import tencentcloud_apm_prometheus_rule.example apm-lhqHyRBuA#140
```

