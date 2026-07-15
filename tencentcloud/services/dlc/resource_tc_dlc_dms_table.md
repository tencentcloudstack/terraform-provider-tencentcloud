Provides a resource to create a DLC DMS table

Example Usage

```hcl
resource "tencentcloud_dlc_dms_table" "example" {
  db_name = "tf_example_db"
  name    = "tf_example_table"
  type    = "EXTERNAL_TABLE"

  asset {
    name        = "tf_example_table"
    description = "tf example dlc dms table"
    owner       = "root"
  }

  columns {
    name     = "id"
    type     = "bigint"
    position = 1
  }

  columns {
    name     = "name"
    type     = "string"
    position = 2
  }

  sds {
    location      = "cosn://tf-example-bucket/example/"
    input_format  = "org.apache.hadoop.hive.ql.io.avro.AvroContainerInputFormat"
    output_format = "org.apache.hadoop.hive.ql.io.avro.AvroContainerOutputFormat"
    serde_lib     = "org.apache.hadoop.hive.serde2.avro.AvroSerDe"
    serde_params {
      key   = "serialization.format"
      value = "1"
    }
  }
}
```

Import

DLC DMS table can be imported using the db_name#name, e.g.

```
terraform import tencentcloud_dlc_dms_table.example tf_example_db#tf_example_table
```
