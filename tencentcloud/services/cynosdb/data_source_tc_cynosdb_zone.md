Use this data source to query detailed information of cynosdb zone

Example Usage

```hcl
data "tencentcloud_cynosdb_zone" "zone" {
  include_virtual_zones = true
  show_permission = true
}
```