Provides a resource to create a WeData lineage attachment

~> **NOTE:** Do not use the same relation parameters for lineage binding, as this will cause overwriting.

Example Usage

```hcl
resource "tencentcloud_wedata_lineage_attachment" "example" {
  relations {
    source {
      resource_unique_id = "2s5veseIo2AXGOHJkKjBvQ"
      resource_type      = "TABLE"
      platform           = "WEDATA"
      resource_name      = "db_demo.1"
      description        = "DLC"
    }

    target {
      resource_unique_id = "fM8OgzE-AM2h4aaJmdXoPg"
      resource_type      = "TABLE"
      platform           = "WEDATA"
      resource_name      = "db_demo.2"
      description        = "DLC"
    }

    processes {
      process_id       = "20241107221758402"
      process_type     = "SCHEDULE_TASK"
      platform         = "WEDATA"
      process_sub_type = "SQL_TASK"
    }
  }
}
```
