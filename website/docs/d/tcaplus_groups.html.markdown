---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcaplus_groups"
sidebar_current: "docs-tencentcloud-datasource-tcaplus_groups"
description: |-
  Use this data source to query tcaplus table groups
---

# tencentcloud_tcaplus_groups

Use this data source to query tcaplus table groups

## Example Usage

```hcl
data "tencentcloud_tcaplus_groups" "null" {
  cluster_id = "19162256624"
}
data "tencentcloud_tcaplus_groups" "id" {
  cluster_id = "19162256624"
  group_id   = "19162256624:1"
}
data "tencentcloud_tcaplus_groups" "name" {
  cluster_id = "19162256624"
  group_name = "test"
}
data "tencentcloud_tcaplus_groups" "all" {
  cluster_id = "19162256624"
  group_id   = "19162256624:1"
  group_name = "test"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) Id of the tcaplus cluster to be query.
* `group_id` - (Optional) Group id to be query.
* `group_name` - (Optional) Group name to be query.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of tcaplus table groups. Each element contains the following attributes.
  * `create_time` - Create time of the tcaplus group.
  * `group_id` - Id of the tcaplus group.
  * `group_name` - Name of the tcaplus group.
  * `table_count` - Number of tables.
  * `total_size` - The total storage(MB).


