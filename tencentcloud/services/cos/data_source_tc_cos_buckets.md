Use this data source to query the COS buckets of the current Tencent Cloud user.

Example Usage

```hcl
data "tencentcloud_cos_buckets" "cos_buckets" {
  bucket_prefix      = "tf-bucket-"
  result_output_file = "mytestpath"
}
```