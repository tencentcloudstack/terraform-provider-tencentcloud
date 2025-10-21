---
subcategory: "TDSQL for MySQL(DCDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dcdb_instance_node_info"
sidebar_current: "docs-tencentcloud-datasource-dcdb_instance_node_info"
description: |-
  Use this data source to query detailed information of dcdb instance_node_info
---

# tencentcloud_dcdb_instance_node_info

Use this data source to query detailed information of dcdb instance_node_info

## Example Usage

```hcl
data "tencentcloud_dcdb_instance_node_info" "instance_node_info" {
  instance_id = local.dcdb_id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID, such as tdsqlshard-6ltok4u9.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `nodes_info` - Node information.
  * `node_id` - Node ID.
  * `role` - Node role. Valid values: `master`, `slave`.
  * `shard_id` - Instance shard ID.


