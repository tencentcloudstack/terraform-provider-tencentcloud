Provides a resource to create a dts compare_task

Example Usage

```hcl
resource "tencentcloud_dts_compare_task" "compare_task" {
  job_id = ""
  task_name = ""
  object_mode = ""
  objects {
			object_mode = ""
		object_items {
				db_name = ""
				db_mode = ""
				schema_name = ""
				table_mode = ""
			tables {
					table_name = ""
			}
				view_mode = ""
			views {
					view_name = ""
			}
		}

  }
  }

```