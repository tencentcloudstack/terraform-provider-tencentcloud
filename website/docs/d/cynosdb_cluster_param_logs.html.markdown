---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_cluster_param_logs"
sidebar_current: "docs-tencentcloud-datasource-cynosdb_cluster_param_logs"
description: |-
  Use this data source to query detailed information of cynosdb cluster_param_logs
---

# tencentcloud_cynosdb_cluster_param_logs

Use this data source to query detailed information of cynosdb cluster_param_logs

## Example Usage

```hcl
data "tencentcloud_cynosdb_cluster_param_logs" "cluster_param_logs" {
  cluster_id    = "cynosdbmysql-bws8h88b"
  instance_ids  = ["cynosdbmysql-ins-afqx1hy0"]
  order_by      = "CreateTime"
  order_by_type = "DESC"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) Cluster ID.
* `instance_ids` - (Optional, Set: [`String`]) Instance ID list, used to record specific instances of operations.
* `order_by_type` - (Optional, String) Define specific sorting rules, limited to one of desc, asc, DESC, or ASC.
* `order_by` - (Optional, String) Sort field, defining which field to sort based on when returning results.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cluster_param_logs` - Parameter modification record note: This field may return null, indicating that a valid value cannot be obtained.
  * `cluster_id` - Cluster ID.
  * `create_time` - Creation time.
  * `current_value` - Current value.
  * `instance_id` - Instance ID.
  * `param_name` - Parameter Name.
  * `status` - modify state.
  * `update_time` - Update time.
  * `update_value` - Modified value.


