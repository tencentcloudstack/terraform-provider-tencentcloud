Use this data source to get the available regions. By default only `AVAILABLE` regions will be returned, but `UNAVAILABLE` regions can also be fetched when `include_unavailable` is specified.

Example Usage

```hcl
data "tencentcloud_availability_regions" "my_favourite_region" {
  name = "ap-guangzhou"
}
```