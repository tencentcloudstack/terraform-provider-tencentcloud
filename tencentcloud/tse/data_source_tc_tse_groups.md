Use this data source to query detailed information of tse groups

Example Usage

```hcl
data "tencentcloud_tse_groups" "groups" {
  gateway_id = "gateway-ddbb709b"
  filters {
    name   = "GroupId"
    values = ["group-013c0d8e"]
  }
}
```