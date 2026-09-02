Provides a resource to create a DLC internal table

Example Usage

Create a basic internal table with columns

```hcl
resource "tencentcloud_dlc_internal_table" "example" {
  database_name = "tf_example_database"
  table_name    = "tf_example_table"
  table_comment = "terraform test table"
  type          = "TABLE"
  table_format  = "hive"
  columns {
    name    = "id"
    type    = "bigint"
    comment = "id column"
  }
  columns {
    name    = "name"
    type    = "string"
    comment = "name column"
  }
  properties {
    key   = "format"
    value = "TEXTFILE"
  }
}
```

Create an Iceberg internal table with smart policy

```hcl
resource "tencentcloud_dlc_internal_table" "example" {
  database_name = "tf_example_database"
  table_name    = "tf_example_iceberg_table"
  table_comment = "terraform iceberg test table"
  type          = "TABLE"
  table_format  = "iceberg"
  primary_keys  = ["id"]
  columns {
    name    = "id"
    type    = "bigint"
    comment = "id column"
    not_null = true
  }
  columns {
    name    = "name"
    type    = "string"
    comment = "name column"
  }
  smart_policy {
    base_info {
      uin        = "100012345678"
      policy_type = "Table"
      catalog    = "DataLakeCatalog"
      database   = "tf_example_database"
      table      = "tf_example_iceberg_table"
      app_id     = "1300123456"
    }
    policy {
      inherit = "false"
      lifecycle {
        lifecycle_enable = "enable"
        expiration       = 30
        expired_field    = "dt"
        expired_field_format = "yyyy-MM-dd"
      }
      table_expiration {
        enabled    = true
        expiration = 90
      }
    }
  }
}
```

Create a table with partitions

```hcl
resource "tencentcloud_dlc_internal_table" "example" {
  database_name = "tf_example_database"
  table_name    = "tf_example_partitioned_table"
  table_comment = "terraform partitioned test table"
  type          = "TABLE"
  table_format  = "hive"
  columns {
    name = "id"
    type = "bigint"
  }
  columns {
    name = "name"
    type = "string"
  }
  partitions {
    name    = "dt"
    type    = "string"
    comment = "date partition"
  }
  properties {
    key   = "format"
    value = "PARQUET"
  }
}
```

Import

DLC internal table can be imported using the composite id (`database_name#table_name`), e.g.

```
terraform import tencentcloud_dlc_internal_table.example tf_example_database#tf_example_table
```
