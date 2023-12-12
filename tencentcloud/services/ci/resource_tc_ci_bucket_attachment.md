Provides a resource to create a ci bucket

Example Usage

```hcl
resource "tencentcloud_ci_bucket_attachment" "bucket_attachment" {
  bucket = "terraform-ci-xxxxxx"
}
```

Import

ci bucket can be imported using the id, e.g.

```
terraform import tencentcloud_ci_bucket_attachment.bucket_attachment terraform-ci-xxxxxx
```