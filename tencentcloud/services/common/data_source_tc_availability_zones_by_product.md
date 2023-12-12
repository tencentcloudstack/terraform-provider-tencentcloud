Use this data source to get the available zones in current region. Must set product param to fetch the product infomations(e.g. => cvm, vpc). By default only `AVAILABLE` zones will be returned, but `UNAVAILABLE` zones can also be fetched when `include_unavailable` is specified.

Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "all" {
  product="cvm"
}
```