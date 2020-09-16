---
subcategory: "CynosDB"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_instances"
sidebar_current: "docs-tencentcloud-datasource-cynosdb_instances"
description: |-
  Use this data source to query detailed information of Cynosdb instances.
---

# tencentcloud_cynosdb_instances

Use this data source to query detailed information of Cynosdb instances.

## Example Usage

```hcl
data "tencentcloud_cynosdb_instances" "foo" {
  instance_id   = "cynosdbmysql-ins-0wln9u6w"
  project_id    = 0
  db_type       = "MYSQL"
  instance_name = "test"
}
```

## Argument Reference

The following arguments are supported:

* `db_type` - (Optional) Type of CynosDB, and available values include `MYSQL`, `POSTGRESQL`.
* `instance_id` - (Optional) ID of the Cynosdb instance to be queried.
* `instance_name` - (Optional) Name of the Cynosdb instance to be queried.
* `project_id` - (Optional) ID of the project to be queried.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_list` - A list of instances. Each element contains the following attributes:
  * `create_time` - Creation time of the CynosDB instance.
  * `instance_name` - Name of CynosDB instance.
  * `instance_status` - Status of the Cynosdb instance.
  * `instance_storage_size` - Storage size of the Cynosdb instance, unit in GB.
  * `instance_type` - Instance type. `ro` for readonly instance, `rw` for read and write instance.


