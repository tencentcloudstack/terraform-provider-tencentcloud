---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_instance_param_records"
sidebar_current: "docs-tencentcloud-datasource-sqlserver_instance_param_records"
description: |-
  Use this data source to query detailed information of sqlserver instance_param_records
---

# tencentcloud_sqlserver_instance_param_records

Use this data source to query detailed information of sqlserver instance_param_records

## Example Usage

```hcl
data "tencentcloud_sqlserver_instance_param_records" "example" {
  instance_id = "mssql-qelbzgwf"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID in the format of mssql-dj5i29c5n. It is the same as the instance ID displayed in the TencentDB console and the response parameter InstanceId of the DescribeDBInstances API.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - Parameter modification records.
  * `instance_id` - Instance ID.
  * `modify_time` - Modification time.
  * `new_value` - Parameter value after modification.
  * `old_value` - Parameter value before modification.
  * `param_name` - Parameter name.
  * `status` - Parameter modification status. Valid values: 1 (initializing and waiting for modification), 2 (modification succeed), 3 (modification failed), 4 (modifying).


