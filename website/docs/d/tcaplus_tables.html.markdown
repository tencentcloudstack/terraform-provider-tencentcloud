---
subcategory: "TcaplusDB"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcaplus_tables"
sidebar_current: "docs-tencentcloud-datasource-tcaplus_tables"
description: |-
  Use this data source to query TcaplusDB tables.
---

# tencentcloud_tcaplus_tables

Use this data source to query TcaplusDB tables.

## Example Usage

```hcl
data "tencentcloud_tcaplus_tables" "null" {
  cluster_id = "19162256624"
}

data "tencentcloud_tcaplus_tables" "tablegroup" {
  cluster_id    = "19162256624"
  tablegroup_id = "19162256624:3"
}

data "tencentcloud_tcaplus_tables" "name" {
  cluster_id    = "19162256624"
  tablegroup_id = "19162256624:3"
  table_name    = "guagua"
}

data "tencentcloud_tcaplus_tables" "id" {
  cluster_id = "19162256624"
  table_id   = "tcaplus-faa65eb7"
}
data "tencentcloud_tcaplus_tables" "all" {
  cluster_id    = "19162256624"
  tablegroup_id = "19162256624:3"
  table_id      = "tcaplus-faa65eb7"
  table_name    = "guagua"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) Id of the TcaplusDB cluster to be query.
* `result_output_file` - (Optional) File for saving results.
* `table_id` - (Optional) Table id to be query.
* `table_name` - (Optional) Table name to be query.
* `tablegroup_id` - (Optional) Id of the table group to be query.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of TcaplusDB tables. Each element contains the following attributes.
  * `create_time` - Create time of the TcaplusDB table.
  * `description` - Description of the TcaplusDB table.
  * `error` - Error message for creating TcaplusDB table.
  * `idl_id` - IDL file id of the TcaplusDB table.
  * `reserved_read_cu` - Reserved read capacity units of the TcaplusDB table.
  * `reserved_volume` - Reserved storage capacity of the TcaplusDB table (unit:GB).
  * `reserved_write_cu` - Reserved write capacity units of the TcaplusDB table.
  * `status` - Status of the TcaplusDB table.
  * `table_id` - Id of the TcaplusDB table.
  * `table_idl_type` - IDL type of  the TcaplusDB table.
  * `table_name` - Name of  the TcaplusDB table.
  * `table_size` - Size of the TcaplusDB table.
  * `table_type` - Type of the TcaplusDB table.
  * `tablegroup_id` - Table group id of the TcaplusDB table.


