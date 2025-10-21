---
subcategory: "Data Transmission Service(DTS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dts_migrate_job_config"
sidebar_current: "docs-tencentcloud-resource-dts_migrate_job_config"
description: |-
  Provides a resource to create a dts migrate_job_config
---

# tencentcloud_dts_migrate_job_config

Provides a resource to create a dts migrate_job_config

## Example Usage

```hcl
resource "tencentcloud_dts_migrate_service" "service" {
  src_database_type = "mysql"
  dst_database_type = "cynosdbmysql"
  src_region        = "ap-guangzhou"
  dst_region        = "ap-guangzhou"
  instance_class    = "small"
  job_name          = "tf_test_xxx"
  tags {
    tag_key   = "aaa"
    tag_value = "bbb"
  }
}

resource "tencentcloud_dts_migrate_job" "job" {
  service_id = tencentcloud_dts_migrate_service.service.id
  run_mode   = "immediate"
  migrate_option {
    database_table {
      object_mode = "partial"
      databases {
        db_name    = "tf_ci_test"
        db_mode    = "partial"
        table_mode = "partial"
        tables {
          table_name      = "test"
          new_table_name  = "test_xxx"
          table_edit_mode = "rename"
        }
      }
    }
  }
  src_info {
    region        = "ap-guangzhou"
    access_type   = "cdb"
    database_type = "mysql"
    node_type     = "simple"
    info {
      user        = "root"
      password    = "xxx"
      instance_id = "id"
    }

  }
  dst_info {
    region        = "ap-guangzhou"
    access_type   = "cdb"
    database_type = "cynosdbmysql"
    node_type     = "simple"
    info {
      user        = "user"
      password    = "xxx"
      instance_id = "id"
    }
  }
  auto_retry_time_range_minutes = 0
}

resource "tencentcloud_dts_migrate_job_start_operation" "start" {
  job_id = tencentcloud_dts_migrate_job.job.id
}

// pause the migration job
resource "tencentcloud_dts_migrate_job_config" "config" {
  job_id = tencentcloud_dts_migrate_job_start_operation.start.id
  action = "pause"
}
```

### Continue the a migration job

```hcl
resource "tencentcloud_dts_migrate_job_config" "config" {
  job_id = tencentcloud_dts_migrate_job_start_operation.start.id
  action = "continue"
}
```

### Complete a migration job when the status is readyComplete

```hcl
resource "tencentcloud_dts_migrate_job_config" "config" {
  job_id = tencentcloud_dts_migrate_job_start_operation.start.id
  action = "continue"
}
```

### Stop a running migration job

```hcl
resource "tencentcloud_dts_migrate_job_config" "config" {
  job_id = tencentcloud_dts_migrate_job_start_operation.start.id
  action = "stop"
}
```

### Isolate a stopped/canceled migration job

```hcl
resource "tencentcloud_dts_migrate_job_config" "config" {
  job_id = tencentcloud_dts_migrate_job_start_operation.start.id
  action = "isolate"
}
```

### Recover a isolated migration job

```hcl
resource "tencentcloud_dts_migrate_job_config" "config" {
  job_id = tencentcloud_dts_migrate_job_start_operation.start.id
  action = "recover"
}
```

## Argument Reference

The following arguments are supported:

* `action` - (Required, String) The operation want to perform. Valid values are: `pause`, `continue`, `complete`, `recover`,`stop`.
* `job_id` - (Required, String) job id.
* `complete_mode` - (Optional, String) complete mode, optional value is waitForSync or immediately.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



