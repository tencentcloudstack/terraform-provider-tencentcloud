---
subcategory: "TDSQL for MySQL(DCDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dcdb_parameters"
sidebar_current: "docs-tencentcloud-datasource-dcdb_parameters"
description: |-
  Use this data source to query detailed information of dcdb parameters
---

# tencentcloud_dcdb_parameters

Use this data source to query detailed information of dcdb parameters

## Example Usage

```hcl
data "tencentcloud_dcdb_parameters" "parameters" {
  instance_id = "your_instance_id"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) instance id.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - parameter list.
  * `constraint` - params constraint.
    * `enum` - a list of optional values of type num.
    * `range` - range constraint.
      * `max` - max value.
      * `min` - min value.
    * `string` - constraint type is string.
    * `type` - type.
  * `default` - default value.
  * `have_set_value` - have set value.
  * `need_restart` - need restart.
  * `param` - parameter name.
  * `value` - parameter value.


