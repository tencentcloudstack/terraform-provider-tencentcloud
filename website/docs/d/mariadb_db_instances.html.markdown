---
subcategory: "TencentDB for MariaDB(MariaDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mariadb_db_instances"
sidebar_current: "docs-tencentcloud-datasource-mariadb_db_instances"
description: |-
  Use this data source to query detailed information of mariadb dbInstances
---

# tencentcloud_mariadb_db_instances

Use this data source to query detailed information of mariadb dbInstances

## Example Usage

```hcl
data "tencentcloud_mariadb_db_instances" "dbInstances" {
  instance_ids = ["tdsql-ijxtqk5p"]
  project_ids  = ["0"]
  vpc_id       = "5556791"
  subnet_id    = "3454730"
}
```

## Argument Reference

The following arguments are supported:

* `instance_ids` - (Optional, Set: [`String`]) instance ids.
* `project_ids` - (Optional, Set: [`Int`]) project ids.
* `result_output_file` - (Optional, String) Used to save results.
* `search_name` - (Optional, String) instance name or vip.
* `subnet_id` - (Optional, String) subnet id.
* `vpc_id` - (Optional, String) vpc id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instances` - instances info.
  * `db_version_id` - db version id.
  * `instance_id` - instance id.
  * `instance_name` - instance name.
  * `memory` - meory of instance.
  * `project_id` - project id.
  * `region` - region.
  * `resource_tags` - resource tags.
    * `tag_key` - tag key.
    * `tag_value` - tag value.
  * `storage` - storage of instance.
  * `subnet_id` - subnet id.
  * `vpc_id` - vpc id.
  * `zone` - available zone.


