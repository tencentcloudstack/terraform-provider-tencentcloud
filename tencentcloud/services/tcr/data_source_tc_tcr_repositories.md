Use this data source to query detailed information of TCR repositories.

Example Usage

```hcl
data "tencentcloud_tcr_repositories" "example" {
  instance_id     = "tcr-12345678"
  namespace_name  = "tf_test_ns"
  repository_name = "test"
}
```