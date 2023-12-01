---
subcategory: "ClickHouse(CDWCH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clickhouse_backup_strategy"
sidebar_current: "docs-tencentcloud-resource-clickhouse_backup_strategy"
description: |-
  Provides a resource to create a clickhouse backup strategy
---

# tencentcloud_clickhouse_backup_strategy

Provides a resource to create a clickhouse backup strategy

## Example Usage

```hcl
resource "tencentcloud_clickhouse_backup" "backup" {
  instance_id     = "cdwch-xxxxxx"
  cos_bucket_name = "xxxxxx"
}

resource "tencentcloud_clickhouse_backup_strategy" "backup_strategy" {
  instance_id = "cdwch-xxxxxx"
  data_backup_strategy {
    week_days    = "3"
    retain_days  = 2
    execute_hour = 1
    back_up_tables {
      database    = "iac"
      table       = "my_table"
      total_bytes = 0
      v_cluster   = "default_cluster"
      ips         = "10.0.0.35"
    }
  }
  meta_backup_strategy {
    week_days    = "1"
    retain_days  = 2
    execute_hour = 3
  }
}
```

## Argument Reference

The following arguments are supported:

* `data_backup_strategy` - (Required, List) Data backup strategy.
* `instance_id` - (Required, String, ForceNew) Instance id.
* `meta_backup_strategy` - (Optional, List) Metadata backup strategy.

The `back_up_tables` object supports the following:

* `database` - (Required, String) Database.
* `table` - (Required, String) Table.
* `total_bytes` - (Required, Int) Back up the list of tables.
* `ips` - (Optional, String) Table ip.
* `rip` - (Optional, String) Ip address of cvm.
* `v_cluster` - (Optional, String) Virtual clusters.
* `zoo_path` - (Optional, String) ZK path.

The `data_backup_strategy` object supports the following:

* `back_up_tables` - (Required, List) Back up the list of tables.
* `execute_hour` - (Required, Int) Execution hour.
* `retain_days` - (Required, Int) Retention days.
* `week_days` - (Required, String) The day of the week is separated by commas. For example: 2 represents Tuesday.

The `meta_backup_strategy` object supports the following:

* `execute_hour` - (Optional, Int) Execution hour.
* `retain_days` - (Optional, Int) Retention days.
* `week_days` - (Optional, String) The day of the week is separated by commas. For example: 2 represents Tuesday.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

clickhouse backup_strategy can be imported using the id, e.g.

```
terraform import tencentcloud_clickhouse_backup_strategy.backup_strategy instance_id
```

