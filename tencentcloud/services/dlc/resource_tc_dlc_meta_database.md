Provides a resource to create a DLC meta database

Example Usage

```hcl
resource "tencentcloud_dlc_meta_database" "example" {
  database_name = "tf_example_db"
  comment       = "tf example meta database"
  govern_policy {
    rule_type     = "Customize"
    govern_engine = "engine_name"
  }
}
```

Import

DLC meta database can be imported using the datasourceConnectionName#databaseName when `datasource_connection_name` is set, or just the bare databaseName otherwise, e.g.

```
terraform import tencentcloud_dlc_meta_database.example DataLakeCatalog#tf_example_db
```

or

```
terraform import tencentcloud_dlc_meta_database.example tf_example_db
```
