Use this data source to query the list of COS vector buckets.

Example Usage

Query all buckets

```hcl
data "tencentcloud_cos_vector_buckets" "example" {}
```

Query buckets by prefix

```hcl
data "tencentcloud_cos_vector_buckets" "example" {
  prefix = "tf-example"
}
```
