---
subcategory: "Provider Meta"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mongodb_restore_db_instance"
sidebar_current: "docs-tencentcloud-action-mongodb_restore_db_instance"
description: |-
  Provides an action to restore a MongoDB instance to a specified point in time via the RestoreDBInstance API. The restore operation rolls back the specified databases and collections to the target time and writes the recovered data into new collections. This is a one-time operation; no cloud-side state is persisted after the action completes.
---

# tencentcloud_mongodb_restore_db_instance

Provides an action to restore a MongoDB instance to a specified point in time via the RestoreDBInstance API. The restore operation rolls back the specified databases and collections to the target time and writes the recovered data into new collections. This is a one-time operation; no cloud-side state is persisted after the action completes.

~> **NOTE:** Actions are supported in HashiCorp Terraform version 1.14 and later.

## Example Usage

```hcl
action "tencentcloud_mongodb_restore_db_instance" "example" {
  config {
    instance_id  = "cmgo-3n1xu3sz"
    restore_time = "2026-09-01 12:00:00"

    databases {
      db = "dbDemo"

      collections {
        old_collection = "col_old_1"
        new_collection = "col_new_1"
      }

      collections {
        old_collection = "col_old_2"
        new_collection = "col_new_2"
      }
    }

    timeouts {
      invoke = "30m"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID, e.g. `cmgo-xxxxxxxx`. Please log in to the MongoDB console and copy the instance ID from the instance list.
* `restore_time` - (Required, String) Target point in time to restore. The time must be within the backup retention period of the instance. Format: `YYYY-MM-DD hh:mm:ss`.
* `databases` - (Optional, List of Object) Database and collection information to restore.
* `timeouts` - (Optional, Object) The timeouts block allows you to specify the timeout for the invoke operation.

The `databases` object supports the following:

* `db` - (Required, String) Database name.

The `collections` object of `databases` supports the following:

* `new_collection` - (Required, String) Collection name after restore.
* `old_collection` - (Required, String) Original collection name to restore.


