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



