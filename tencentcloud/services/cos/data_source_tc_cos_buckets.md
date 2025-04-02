Use this data source to query the COS buckets of the current Tencent Cloud user.

Example Usage

Query all cos buckets

```hcl
data "tencentcloud_cos_buckets" "example" {}
```

Query cos buckets by filters

```hcl
data "tencentcloud_cos_buckets" "example" {
  bucket_prefix = "tf-example-prefix"
  tags = {
    createBy = "Terraform"
  }
}
```
