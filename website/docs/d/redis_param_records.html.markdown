---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_param_records"
sidebar_current: "docs-tencentcloud-datasource-redis_param_records"
description: |-
  Use this data source to query detailed information of redis param records
---

# tencentcloud_redis_param_records

Use this data source to query detailed information of redis param records

## Example Usage

```hcl
data "tencentcloud_redis_param_records" "param_records" {
  instance_id = "crs-c1nl9rpv"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) The ID of instance.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_param_history` - The parameter name.
  * `modify_time` - Modification time.
  * `new_value` - The modified value.
  * `param_name` - The parameter name.
  * `pre_value` - Modify the previous value.
  * `status` - Parameter status:1: parameter configuration modification.2: The parameter configuration is modified successfully.3: Parameter configuration modification failed.


