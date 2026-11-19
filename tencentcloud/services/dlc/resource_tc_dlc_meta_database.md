Provides a resource to create a DLC meta database

Example Usage

```hcl
resource "tencentcloud_dlc_meta_database" "example" {
  datasource_connection_name = "DataLakeCatalog"
  meta_database_info {
    database_name = "tf_example_db"
    comment       = "tf example meta database"
  }
}
```

Import

DLC meta database can be imported using the datasourceConnectionName#databaseName, e.g.

```
terraform import tencentcloud_dlc_meta_database.example DataLakeCatalog#tf_example_db
```
