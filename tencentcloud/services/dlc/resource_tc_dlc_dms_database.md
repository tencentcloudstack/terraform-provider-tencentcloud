Provides a resource to create a DMS database for Tencent Cloud DLC (Data Lake Compute).

Example Usage

```hcl
resource "tencentcloud_dlc_dms_database" "example" {
  name                      = "tf_example_dms_database"
  schema_name               = "DataLake"
  datasource_connection_name = "tf_example_connection"
  location                  = "cosn://tf-example-bucket-1300000000/data/"

  asset {
    name        = "tf_example_asset"
    catalog     = "DataLake"
    description = "example dms database asset."
    owner       = "root"

    params {
      key   = "param_key"
      value = "param_value"
    }

    biz_params {
      key   = "biz_key"
      value = "biz_value"
    }
  }
}
```

Import

DLC DMS database can be imported using the compound id `name#schema_name#datasource_connection_name`, e.g.

```
terraform import tencentcloud_dlc_dms_database.example tf_example_dms_database#DataLake#tf_example_connection
```
