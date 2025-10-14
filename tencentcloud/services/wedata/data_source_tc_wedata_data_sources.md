Use this data source to query detailed information of WeData data sources

Example Usage

```hcl
data "tencentcloud_wedata_data_sources" "example" {
  project_id   = "2982667120655491072"
  name         = "tf_example"
  display_name = "display_name"
  type         = ["MYSQL", "ORACLE"]
  creator      = "user"
}
```
