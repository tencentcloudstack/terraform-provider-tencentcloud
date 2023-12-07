Use this data source to get the available zones in current region. By default only `AVAILABLE` zones will be returned, but `UNAVAILABLE` zones can also be fetched when `include_unavailable` is specified.

~> **NOTE:** It has been deprecated and replaced by tencentcloud_availability_zones_by_product.

Example Usage

```hcl
data "tencentcloud_availability_zones" "my_favourite_zone" {
  name = "ap-guangzhou-3"
}
```