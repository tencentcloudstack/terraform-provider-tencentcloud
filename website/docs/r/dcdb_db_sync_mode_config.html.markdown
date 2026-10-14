---
subcategory: "TDSQL for MySQL(DCDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dcdb_db_sync_mode_config"
sidebar_current: "docs-tencentcloud-resource-dcdb_db_sync_mode_config"
description: |-
  Provides a resource to create a dcdb db_sync_mode_config
---

# tencentcloud_dcdb_db_sync_mode_config

Provides a resource to create a dcdb db_sync_mode_config

## Example Usage

```hcl
resource "tencentcloud_dcdb_db_sync_mode_config" "config" {
  instance_id = "%s"
  sync_mode   = 2
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) ID of the instance for which to modify the sync mode. The ID is in the format of `tdsql-ow728lmc`.
* `sync_mode` - (Required, Int) Sync mode. Valid values: `0` (async), `1` (strong sync), `2` (downgradable strong sync).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DescribeDBSyncMode](https://cloud.tencent.com/document/api/557/19991) - This API is used to query the sync mode of a TencentDB instance.
* Create: [ModifyDBSyncMode](https://cloud.tencent.com/document/api/557/19981) - This API is used to modify the sync mode of a TencentDB instance.
* Read: [DescribeDBSyncMode](https://cloud.tencent.com/document/api/557/19991) - This API is used to query the sync mode of a TencentDB instance.
* Update: [DescribeDBSyncMode](https://cloud.tencent.com/document/api/557/19991) - This API is used to query the sync mode of a TencentDB instance.
* Update: [ModifyDBSyncMode](https://cloud.tencent.com/document/api/557/19981) - This API is used to modify the sync mode of a TencentDB instance.

## Import

dcdb db_sync_mode_config can be imported using the id, e.g.

```
terraform import tencentcloud_dcdb_db_sync_mode_config.db_sync_mode_config db_sync_mode_config_id
```

