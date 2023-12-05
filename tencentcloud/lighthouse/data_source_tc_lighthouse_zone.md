Use this data source to query detailed information of lighthouse zone

Example Usage

```hcl
data "tencentcloud_lighthouse_zone" "zone" {
  order_field = "ZONE"
  order = "ASC"
}
```