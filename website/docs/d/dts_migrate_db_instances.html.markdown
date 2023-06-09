---
subcategory: "Data Transmission Service(DTS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dts_migrate_db_instances"
sidebar_current: "docs-tencentcloud-datasource-dts_migrate_db_instances"
description: |-
  Use this data source to query detailed information of dts migrate_db_instances
---

# tencentcloud_dts_migrate_db_instances

Use this data source to query detailed information of dts migrate_db_instances

## Example Usage

```hcl
data "tencentcloud_dts_migrate_db_instances" "migrate_db_instances" {
  database_type  = "mysql"
  migrate_role   = "src"
  instance_id    = "cdb-ffulb2sg"
  instance_name  = "cdb_test"
  limit          = 10
  offset         = 10
  account_mode   = "self"
  tmp_secret_id  = "AKIDvBDyVmna9TadcS4YzfBZmkU5TbX12345"
  tmp_secret_key = "ZswjGWWHm24qMeiX6QUJsELDpC12345"
  tmp_token      = "JOqqCPVuWdNZvlVDLxxx"
}
```

## Argument Reference

The following arguments are supported:

* `database_type` - (Required, String) Database type.
* `account_mode` - (Optional, String) The owning account of the resource is null or self(resources in the self account), other(resources in the other account).
* `instance_id` - (Optional, String) Database instance id.
* `instance_name` - (Optional, String) Database instance name.
* `limit` - (Optional, Int) Limit.
* `migrate_role` - (Optional, String) Whether the instance is the migration source or destination,src(for source), dst(for destination).
* `offset` - (Optional, Int) Offset.
* `result_output_file` - (Optional, String) Used to save results.
* `tmp_secret_id` - (Optional, String) temporary secret id, used across account.
* `tmp_secret_key` - (Optional, String) temporary secret key, used across account.
* `tmp_token` - (Optional, String) temporary token, used across account.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instances` - Instance list.
  * `hint` - The reason of can&#39;t used in migration.
  * `instance_id` - Instance Id.
  * `instance_name` - Database instance name.
  * `usable` - Can used in migration, 1-yes, 0-no.
  * `vip` - Instance vip.
  * `vport` - Instance port.
* `request_id` - Unique request id, provide this when encounter a problem.


