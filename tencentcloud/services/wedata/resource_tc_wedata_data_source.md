Provides a resource to create a WeData data source

Example Usage

```hcl
resource "tencentcloud_wedata_data_source" "example" {
  project_id = "2983848457986924544"
  name       = "tf_example"
  type       = "MYSQL"
  prod_con_properties = jsonencode({
    "deployType" : "CONNSTR_PUBLICDB",
    "url" : "jdbc:mysql://1.1.1.1:1111/database",
    "username" : "root",
    "password" : "root"
  })

  display_name = "display_name"
  description  = "description"

  lifecycle {
    ignore_changes = [
      prod_con_properties,
    ]
  }
}
```
