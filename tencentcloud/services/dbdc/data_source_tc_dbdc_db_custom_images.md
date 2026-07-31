Use this data source to query available DB Custom OS images from the TencentCloud DBDC product.

Example Usage

Query all dbdc db custom images

```hcl
data "tencentcloud_dbdc_db_custom_images" "example" {}
```

Query dbdc db custom images by os-type filter

```hcl
data "tencentcloud_dbdc_db_custom_images" "example" {
  name   = "os-type"
  values = ["linux"]
}
```
