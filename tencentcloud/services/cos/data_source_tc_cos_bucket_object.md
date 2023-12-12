Use this data source to query the metadata of an object stored inside a bucket.

Example Usage

```hcl
data "tencentcloud_cos_bucket_object" "mycos" {
  bucket             = "mycos-test-1258798060"
  key                = "hello-world.py"
  result_output_file = "TFresults"
}
```