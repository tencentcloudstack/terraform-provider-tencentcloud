Provides a resource to create a cynosdb roll_back_cluster

Example Usage

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