---
subcategory: "TencentDB for MariaDB(MariaDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mariadb_operate_hour_db_instance"
sidebar_current: "docs-tencentcloud-resource-mariadb_operate_hour_db_instance"
description: |-
  Provides a resource to create a mariadb activate_hour_db_instance
---

# tencentcloud_mariadb_operate_hour_db_instance

Provides a resource to create a mariadb activate_hour_db_instance

## Example Usage

```hcl
resource "tencentcloud_mariadb_operate_hour_db_instance" "activate_hour_db_instance" {
  instance_id = "tdsql-9vqvls95"
  operate     = "activate"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID.
* `operate` - (Required, String) Operation, `activate`- activate the hour db instance, `isolate`- isolate the hour db instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [ActivateHourDBInstance](https://cloud.tencent.com/document/api/237/70426) - 解隔离MariaDB按量计费实例
* Create: [IsolateHourDBInstance](https://cloud.tencent.com/document/api/237/70425) - 隔离MariaDB按量计费实例
* Read: [DescribeDBInstanceDetail](https://cloud.tencent.com/document/api/237/89390) - 查询实例详细信息
* Update: [ActivateHourDBInstance](https://cloud.tencent.com/document/api/237/70426) - 解隔离MariaDB按量计费实例
* Update: [IsolateHourDBInstance](https://cloud.tencent.com/document/api/237/70425) - 隔离MariaDB按量计费实例

