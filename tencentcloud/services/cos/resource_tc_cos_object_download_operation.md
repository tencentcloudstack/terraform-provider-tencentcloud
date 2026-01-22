Provides a resource to download object

Example Usage

```hcl
resource "tencentcloud_cos_object_download_operation" "example" {
  bucket        = "private-bucket-1309116523"
  key           = "demo.txt"
  download_path = "/tmp/demo.txt"

  timeouts {
    create = "10m"
  }
}
```