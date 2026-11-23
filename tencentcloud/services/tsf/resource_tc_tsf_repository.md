Provides a resource to create a TSF repository

Example Usage

```hcl
resource "tencentcloud_tsf_repository" "example" {
  repository_name = "tf-example"
  repository_type = "default"
  bucket_name     = "tf-example-1308919341"
  bucket_region   = "ap-guangzhou"
  directory       = "/"
  repository_desc = "desc."
}
```

Import

TSF repository can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_repository.example repository-xxxxxxxx
```
