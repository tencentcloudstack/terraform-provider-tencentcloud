---
subcategory: "TencentDB for MongoDB(mongodb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mongodb_instance_params"
sidebar_current: "docs-tencentcloud-datasource-mongodb_instance_params"
description: |-
  Use this data source to query detailed information of mongodb instance_params
---

# tencentcloud_mongodb_instance_params

Use this data source to query detailed information of mongodb instance_params

## Example Usage

```hcl
data "tencentcloud_mongodb_instance_params" "instance_params" {
  instance_id = "cmgo-gwqk8669"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) InstanceId.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_enum_param` - Enum parameter.
  * `current_value` - current value.
  * `default_value` - default value.
  * `enum_value` - enumvalue.
  * `need_restart` - if need restart.
  * `param_name` - name of parameter.
  * `status` - if is running.
  * `tips` - descripition of parameter.
  * `value_type` - value type.
* `instance_integer_param` - Integer parameter.
  * `current_value` - current value.
  * `default_value` - default value.
  * `max` - max value.
  * `min` - min value.
  * `need_restart` - if need restart.
  * `param_name` - name of parameter.
  * `status` - if is running.
  * `tips` - descripition of parameter.
  * `value_type` - value type.
* `instance_multi_param` - multi parameter.
  * `current_value` - current value.
  * `default_value` - default value.
  * `enum_value` - enum value.
  * `need_restart` - if need restart.
  * `param_name` - name of parameter.
  * `status` - if is running.
  * `tips` - descripition of parameter.
  * `value_type` - value type.
* `instance_text_param` - text parameter.
  * `current_value` - current value.
  * `default_value` - default value.
  * `need_restart` - if need restart.
  * `param_name` - name of parameter.
  * `status` - if is running.
  * `text_value` - text value.
  * `tips` - descripition of parameter.
  * `value_type` - value type.


