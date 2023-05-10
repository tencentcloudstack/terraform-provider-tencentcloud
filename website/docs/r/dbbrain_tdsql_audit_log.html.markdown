---
subcategory: "TencentDB for DBbrain(dbbrain)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbbrain_tdsql_audit_log"
sidebar_current: "docs-tencentcloud-resource-dbbrain_tdsql_audit_log"
description: |-
  Provides a resource to create a dbbrain tdsql_audit_log
---

# tencentcloud_dbbrain_tdsql_audit_log

Provides a resource to create a dbbrain tdsql_audit_log

## Example Usage

```hcl
resource "tencentcloud_dbbrain_tdsql_audit_log" "my_log" {
  product           = "dcdb"
  node_request_type = "dcdb"
  instance_id       = "%s"
  start_time        = "%s"
  end_time          = "%s"
  filter {
    host = ["%%", "127.0.0.1"]
    user = ["tf_test", "mysql"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String, ForceNew) Deadline time, such as `2019-09-11 10:13:14`.
* `instance_id` - (Required, String, ForceNew) Instance ID.
* `node_request_type` - (Required, String, ForceNew) Consistent with Product. For example: dcdb, mariadb.
* `product` - (Required, String, ForceNew) Service product type, supported values include: dcdb - cloud database Tdsql, mariadb - cloud database MariaDB for MariaDB..
* `start_time` - (Required, String, ForceNew) Start time, such as `2019-09-10 12:13:14`.
* `filter` - (Optional, List, ForceNew) Filter conditions. Logs can be filtered according to the filter conditions set.

The `filter` object supports the following:

* `affect_rows` - (Optional, Int, ForceNew) Number of affected rows. Indicates filtering audit logs whose affected rows are greater than this value.
* `db_name` - (Optional, Set, ForceNew) Database name.
* `exec_time` - (Optional, Int, ForceNew) Execution time. The unit is: us. It means to filter the audit logs whose execution time is greater than this value.
* `host` - (Optional, Set, ForceNew) Client Address.
* `sent_rows` - (Optional, Int, ForceNew) Return the number of rows. It means to filter the audit log with the number of returned rows greater than this value.
* `user` - (Optional, Set, ForceNew) Username.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



