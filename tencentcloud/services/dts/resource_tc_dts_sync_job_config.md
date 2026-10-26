Provides a resource to create a DTS sync job config

Example Usage

```hcl
resource "tencentcloud_dts_sync_job_config" "example" {
  job_id = "sync-werwfs23"

  dynamic_options {
    op_types     = ["Insert", "Update", "Delete", "DDL"]
    ddl_options {
      ddl_object = "Table"
      ddl_value     = ["Create", "Drop"]
    }

    conflict_handle_type = "ReportError"

    conflict_handle_option {
      condition_column               = "id"
      condition_operator             = ">"
      condition_order_in_src_and_dst = ">"
    }

    deal_of_exist_same_table = "ReportErrorAfterCheck"
  }

  dynamic_objects {
    mode = "Partial"

    databases {
      db_name = "test_db"
      db_mode = "Partial"

      tables {
        table_name = "test_table"
      }
    }

    online_ddl {
      status = "OFF"
    }
  }
}
```

Import

DTS sync job config can be imported using the id, e.g.

```
terraform import tencentcloud_dts_sync_job_config.example sync-werwfs23
```