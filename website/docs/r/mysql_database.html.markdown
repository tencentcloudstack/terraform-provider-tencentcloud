---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_database"
sidebar_current: "docs-tencentcloud-resource-mysql_database"
description: |-
  Provides a resource to create a mysql database
---

# tencentcloud_mysql_database

Provides a resource to create a mysql database

## Example Usage

```hcl
resource "tencentcloud_mysql_database" "database" {
  instance_id        = "cdb-i9xfdf7z"
  db_name            = "for_tf_test"
  character_set_name = "utf8"
}
```

## Argument Reference

The following arguments are supported:

* `character_set_name` - (Required, String) Character set. Valid values:  `utf8`, `gbk`, `latin1`, `utf8mb4`.
* `db_name` - (Required, String, ForceNew) Name of Database.
* `instance_id` - (Required, String, ForceNew) Instance ID in the format of `cdb-c1nl9rpv`,  which is the same as the one displayed in the TencentDB console.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateDatabase](https://cloud.tencent.com/document/api/236/94238) - This API is used to create a database in a TencentDB instance.
* Read: [DescribeDatabases](https://cloud.tencent.com/document/api/236/17493) - This API is used to query the information of databases in a TencentDB instance which must be a source or disaster recovery instance.
* Delete: [DeleteDatabase](https://cloud.tencent.com/document/api/236/101808) - This API is used to delete a database in a cloud database instance.

## Import

mysql database can be imported using the id, e.g.

```
terraform import tencentcloud_mysql_database.database instanceId#dbName
```

