Use this data source to query the COS bucket multipart uploads.

Example Usage

```hcl
data "tencentcloud_cos_bucket_multipart_uploads" "cos_bucket_multipart_uploads" {
	bucket = "xxxxxx"
}
```