---
subcategory: "Config"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_config_rule_evaluation_results"
sidebar_current: "docs-tencentcloud-datasource-config_rule_evaluation_results"
description: |-
  Use this data source to query detailed information of Config rule evaluation results (rule dimension).
---

# tencentcloud_config_rule_evaluation_results

Use this data source to query detailed information of Config rule evaluation results (rule dimension).

## Example Usage

### Query evaluation results by rule ID

```hcl
data "tencentcloud_config_rule_evaluation_results" "example" {
  config_rule_id = "cr-pHmVQS1UpihV4MSTkmIo"
}
```

### Query evaluation results with filters

```hcl
data "tencentcloud_config_rule_evaluation_results" "example" {
  config_rule_id  = "cr-pHmVQS1UpihV4MSTkmIo"
  compliance_type = ["NON_COMPLIANT"]
  resource_type   = ["QCS::CVM::Instance"]
}
```

## Argument Reference

The following arguments are supported:

* `config_rule_id` - (Required, String) Config rule ID.
* `compliance_type` - (Optional, List: [`String`]) Compliance type list for filtering. Valid values: COMPLIANT, NON_COMPLIANT.
* `resource_type` - (Optional, List: [`String`]) Resource type list for filtering (e.g. QCS::CVM::Instance).
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result_list` - Evaluation result list.
  * `annotation` - Evaluation annotation detail.
    * `configuration` - Actual resource configuration (non-compliant configuration).
    * `desired_value` - Expected resource configuration (compliant configuration).
    * `operator` - Comparison operator between actual and expected configuration.
    * `property` - JSON path of the current configuration in the resource attribute structure.
  * `compliance_pack_id` - Compliance pack ID.
  * `compliance_type` - Compliance type. Valid values: COMPLIANT, NON_COMPLIANT.
  * `config_rule_id` - Config rule ID.
  * `config_rule_invoked_time` - Evaluation invocation time.
  * `config_rule_name` - Config rule name.
  * `invoking_event_message_type` - Rule invocation type.
  * `resource_id` - Resource ID.
  * `resource_name` - Resource name.
  * `resource_region` - Resource region.
  * `resource_type` - Resource type.
  * `result_recorded_time` - Evaluation result recorded time.
  * `risk_level` - Risk level. Valid values: 1 (high risk), 2 (medium risk), 3 (low risk).


