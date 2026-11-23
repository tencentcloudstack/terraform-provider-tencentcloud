Use this data source to query detailed information of DC instances.

Example Usage

```hcl
data "tencentcloud_dc_instances" "name_select" {
  name = "t"
}

data "tencentcloud_dc_instances" "id" {
  dc_id = "dc-kax48sg7"
}
```