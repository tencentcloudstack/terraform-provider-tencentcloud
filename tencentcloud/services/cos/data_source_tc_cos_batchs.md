Use this data source to query the COS batch.

~> **NOTE:** The current resource does not support `cos_domain`.

Example Usage

```hcl
data "tencentcloud_cos_batchs" "cos_batchs" {
  uin = "xxxxxx"
  appid = "xxxxxx"
}
```