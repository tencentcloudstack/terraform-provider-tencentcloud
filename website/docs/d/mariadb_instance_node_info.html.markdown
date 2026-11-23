---
subcategory: "TencentDB for MariaDB(MariaDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mariadb_instance_node_info"
sidebar_current: "docs-tencentcloud-datasource-mariadb_instance_node_info"
description: |-
  Use this data source to query detailed information of mariadb instance_node_info
---

# tencentcloud_mariadb_instance_node_info

Use this data source to query detailed information of mariadb instance_node_info

## Example Usage

```hcl
data "tencentcloud_mariadb_instance_node_info" "instance_node_info" {
  instance_id = "tdsql-9vqvls95"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID, such as tdsql-6ltok4u9.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `nodes_info` - Node information.
  * `node_id` - Node ID.
  * `role` - Node role. Valid values: `master`, `slave`.

## API Reference

This data source queries the following TencentCloud APIs:

* Read: [DescribeInstanceNodeInfo](https://www.tencentcloud.com/document/api/237/40253)

