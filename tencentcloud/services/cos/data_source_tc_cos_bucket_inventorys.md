Use this data source to query the COS bucket inventorys.

~> **NOTE:** The current resource does not support cdc.

Example Usage

```hcl
data "tencentcloud_cos_bucket_inventorys" "cos_bucket_inventorys" {
	bucket = "xxxxxx"
}
```