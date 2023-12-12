Provides a resource to create a cos bucket_referer

Example Usage

```hcl
resource "tencentcloud_cos_bucket_referer" "bucket_referer" {
  bucket = "mycos-1258798060"
  status = "Enabled"
  referer_type = "Black-List"
  domain_list = ["127.0.0.1", "terraform.com"]
  empty_refer_configuration = "Allow"
}
```

Import

cos bucket_referer can be imported using the id, e.g.

```
terraform import tencentcloud_cos_bucket_referer.bucket_referer bucket_id
```