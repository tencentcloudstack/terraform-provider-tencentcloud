Use this data source to query the available mongodb specifications for different zone.

Example Usage

```hcl
data "tencentcloud_mongodb_zone_config" "mongodb" {
  available_zone = "ap-guangzhou-2"
}
```