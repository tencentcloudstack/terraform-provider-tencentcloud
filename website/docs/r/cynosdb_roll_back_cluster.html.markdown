---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_roll_back_cluster"
sidebar_current: "docs-tencentcloud-resource-cynosdb_roll_back_cluster"
description: |-
  Provides a resource to create a cynosdb roll_back_cluster
---

# tencentcloud_cynosdb_roll_back_cluster

Provides a resource to create a cynosdb roll_back_cluster

## Example Usage

```hcl
resource "tencentcloud_cynosdb_roll_back_cluster" "roll_back_cluster" {
  cluster_id        = "cynosdbmysql-bws8h88b"
  rollback_strategy = "snapRollback"
  rollback_id       = 732725
  # expect_time = "2022-01-20 00:00:00"
  expect_time_thresh = 0
  rollback_databases {
    old_database = "users"
    new_database = "users_bak_1"
  }
  rollback_tables {
    database = "tf_ci_test"
    tables {
      old_table = "test"
      new_table = "test_bak_111"
    }

  }
  rollback_mode = "full"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) The ID of cluster.
* `rollback_id` - (Required, Int, ForceNew) Rollback ID.
* `rollback_strategy` - (Required, String, ForceNew) Backfile policy timeRollback - Backfile by point in time snapRollback - Backfile by backup file.
* `expect_time_thresh` - (Optional, Int, ForceNew) Expected Threshold (Obsolete).
* `expect_time` - (Optional, String, ForceNew) Expected rollback Time.
* `rollback_databases` - (Optional, List, ForceNew) Database list.
* `rollback_mode` - (Optional, String, ForceNew) Rollback mode by time point, full: normal; Db: fast; Table: Extreme speed (default is normal).
* `rollback_tables` - (Optional, List, ForceNew) Table list.

The `rollback_databases` object supports the following:

* `new_database` - (Required, String) New database name.
* `old_database` - (Required, String) Old database name.

The `rollback_tables` object supports the following:

* `database` - (Required, String) New database name.
* `tables` - (Required, List) Tables.

The `tables` object of `rollback_tables` supports the following:

* `new_table` - (Required, String) New table name.
* `old_table` - (Required, String) Old table name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



