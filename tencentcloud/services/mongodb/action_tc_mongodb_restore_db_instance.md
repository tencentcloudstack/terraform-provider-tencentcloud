Provides an action to restore a MongoDB instance to a specified point in time via the RestoreDBInstance API. The restore operation rolls back the specified databases and collections to the target time and writes the recovered data into new collections. This is a one-time operation; no cloud-side state is persisted after the action completes.

~> **NOTE:** Actions are supported in HashiCorp Terraform version 1.14 and later.

Example Usage

```hcl
action "tencentcloud_mongodb_restore_db_instance" "example" {
  config {
    instance_id  = "cmgo-xxxxxxxx"
    restore_time = "2024-09-01 12:00:00"

    databases {
      db = "db1"

      collections {
        old_collection = "col_old_1"
        new_collection = "col_new_1"
      }

      collections {
        old_collection = "col_old_2"
        new_collection = "col_new_2"
      }
    }
  }
}
```