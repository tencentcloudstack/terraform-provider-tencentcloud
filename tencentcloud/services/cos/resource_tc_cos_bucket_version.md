Provides a resource to create a cos bucket_version

Example Usage

```hcl
resource "tencentcloud_cos_bucket_version" "bucket_version" {
  bucket = "mycos-1258798060"
  status = "Enabled"
}
```

Import

cos bucket_version can be imported using the id, e.g.

```
terraform import tencentcloud_cos_bucket_version.bucket_version bucket_id
```