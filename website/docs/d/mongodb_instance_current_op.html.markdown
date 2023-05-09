---
subcategory: "TencentDB for MongoDB(mongodb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mongodb_instance_current_op"
sidebar_current: "docs-tencentcloud-datasource-mongodb_instance_current_op"
description: |-
  Use this data source to query detailed information of mongodb instance_current_op
---

# tencentcloud_mongodb_instance_current_op

Use this data source to query detailed information of mongodb instance_current_op

## Example Usage

```hcl
data "tencentcloud_mongodb_instance_current_op" "instance_current_op" {
  instance_id   = "cmgo-b43i3wkj"
  op            = "command"
  order_by_type = "desc"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID, the format is: cmgo-9d0p6umb.Same as the instance ID displayed in the cloud database console page.
* `millisecond_running` - (Optional, Int) Filter condition, the time that the operation has been executed (unit: millisecond),the result will return the operation that exceeds the set time, the default value is 0,and the value range is [0, 3600000].
* `ns` - (Optional, String) Filter condition, the namespace namespace to which the operation belongs, in the format of db.collection.
* `op` - (Optional, String) Filter condition, operation type, possible values: none, update, insert, query, command, getmore,remove and killcursors.
* `order_by_type` - (Optional, String) Returns the sorting method of the result set, possible values: ASC/asc or DESC/desc.
* `order_by` - (Optional, String) Returns the sorted field of the result set, currently supports: MicrosecsRunning/microsecsrunning,the default is ascending sort.
* `replica_set_name` - (Optional, String) filter condition, shard name.
* `result_output_file` - (Optional, String) Used to save results.
* `state` - (Optional, String) Filter condition, node status, possible value: primary, secondary.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `current_ops` - current operation list.
  * `microsecs_running` - running time(ms).
  * `node_name` - Node name.
  * `ns` - operation namespace.
  * `op_id` - operation id.
  * `op` - operation value.
  * `operation` - operation info.
  * `query` - operation query.
  * `replica_set_name` - Replication name.
  * `state` - operation state.


