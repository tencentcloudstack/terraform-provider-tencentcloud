Use this data source to query the list of indexes of a COS vector bucket.

Example Usage

```hcl
data "tencentcloud_cos_vector_indexes" "example" {
  vector_bucket_name = "tf-example-1309101522"
  prefix             = "tf-example"
}
```
