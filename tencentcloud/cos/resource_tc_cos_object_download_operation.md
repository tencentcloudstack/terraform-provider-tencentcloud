Provides a resource to download object

Example Usage

```hcl
resource "tencentcloud_cos_object_download_operation" "object_download" {
    bucket = "xxxxxxx"
    key = "test.txt"
    download_path = "/tmp/test.txt"
}
```