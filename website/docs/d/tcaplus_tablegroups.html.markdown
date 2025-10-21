---
subcategory: "TcaplusDB"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcaplus_tablegroups"
sidebar_current: "docs-tencentcloud-datasource-tcaplus_tablegroups"
description: |-
  Use this data source to query table groups of the TcaplusDB cluster.
---

# tencentcloud_tcaplus_tablegroups

Use this data source to query table groups of the TcaplusDB cluster.

## Example Usage

```hcl
data "tencentcloud_tcaplus_tablegroups" "null" {
  cluster_id = "19162256624"
}
data "tencentcloud_tcaplus_tablegroups" "id" {
  cluster_id    = "19162256624"
  tablegroup_id = "19162256624:1"
}
data "tencentcloud_tcaplus_tablegroups" "name" {
  cluster_id      = "19162256624"
  tablegroup_name = "test"
}
data "tencentcloud_tcaplus_tablegroups" "all" {
  cluster_id      = "19162256624"
  tablegroup_id   = "19162256624:1"
  tablegroup_name = "test"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) Id of the TcaplusDB cluster to be query.
* `result_output_file` - (Optional, String) File for saving results.
* `tablegroup_id` - (Optional, String) Id of the table group to be query.
* `tablegroup_name` - (Optional, String) Name of the table group to be query.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of table group. Each element contains the following attributes.
  * `create_time` - Create time of the table group..
  * `table_count` - Number of tables.
  * `tablegroup_id` - Id of the table group.
  * `tablegroup_name` - Name of the table group.
  * `total_size` - Total storage size (MB).


