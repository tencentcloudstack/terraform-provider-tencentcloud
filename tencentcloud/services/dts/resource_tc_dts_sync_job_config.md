Provides a resource to create a DTS sync job config

~> **NOTE:** Only `Running` sync jobs can be configured.

Example Usage

```hcl
resource "tencentcloud_dts_sync_job_config" "example" {
  job_id = "sync-6zm6nj70"

  dynamic_options {
    op_types = ["Insert", "Update", "Delete", "PartialDDL"]
    ddl_options {
      ddl_object = "Database"
      ddl_value  = ["Create", "Drop", "Alter"]
    }

    ddl_options {
      ddl_object = "Table"
      ddl_value  = ["Alter", "Create", "Drop", "Rename", "Truncate"]
    }
    ddl_options {
      ddl_object = "View"
      ddl_value  = ["Create", "Drop"]
    }
    ddl_options {
      ddl_object = "Index"
      ddl_value  = ["Create", "Drop"]
    }

    kafka_option {
      data_option {
        key   = "skipLargeMessage"
        value = "off"
      }
    }

    conflict_handle_type     = "ReportError"
    deal_of_exist_same_table = "ReportErrorAfterCheck"
  }

  dynamic_objects {
    mode = "Partial"
    databases {
      db_name        = "db1"
      db_mode        = "All"
      table_mode     = "All"
      view_mode      = "All"
      procedure_mode = "All"
      function_mode  = "All"
    }

    advanced_objects = ["procedure", "function"]
  }
}
```

Import

DTS sync job config can be imported using the id, e.g.

```
terraform import tencentcloud_dts_sync_job_config.example sync-6zm6nj70
```
