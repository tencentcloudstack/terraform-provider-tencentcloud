Provides a resource to create a wedata script

Example Usage

```hcl
resource "tencentcloud_wedata_script" "example" {
  file_path           = "/datastudio/project/tf_example.sql"
  project_id          = "1470575647377821696"
  bucket_name         = "wedata-demo-1257305158"
  region              = "ap-guangzhou"
  file_extension_type = "sql"
}
```

Import

wedata script can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_script.example 1470575647377821696#/datastudio/project/tf_example.sql#4147824b-7ba2-432b-8a8b-7e747594c926
```