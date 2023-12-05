Provides a resource to abort multipart upload

Example Usage

```hcl
resource "tencentcloud_cos_object_abort_multipart_upload_operation" "abort_multipart_upload" {
    bucket = "keep-test-xxxxxx"
    key = "object"
    upload_id = "xxxxxx"
}
```