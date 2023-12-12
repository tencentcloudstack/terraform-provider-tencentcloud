Provides a resource to create a wedata resource

Example Usage

```hcl
resource "tencentcloud_wedata_resource" "example" {
  file_path       = "/datastudio/resource/demo"
  project_id      = "1612982498218618880"
  file_name       = "tf_example"
  cos_bucket_name = "wedata-demo-1314991481"
  cos_region      = "ap-guangzhou"
  files_size      = "8165"
}
```

Import

wedata resource can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_resource.example 1612982498218618880#/datastudio/resource/demo#75431931-7d27-4034-b3de-3dc3348a220e
```