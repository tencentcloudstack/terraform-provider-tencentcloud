Provides a resource to copy object

Example Usage

```hcl
resource "tencentcloud_cos_object_copy_operation" "object_copy" {
    bucket = "keep-copy-xxxxxxx"
    key = "copy-acl.txt"
    source_url = "keep-test-xxxxxx.cos.ap-guangzhou.myqcloud.com/acl.txt"
}
```