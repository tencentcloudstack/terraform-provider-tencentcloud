Use this data source to query detailed information of DLC data engine network

Example Usage

```hcl
data "tencentcloud_dlc_data_engine_network" "example" {
  sort_by = "create-time"
  sorting = "desc"
  filters {
    name   = "engine-network-id"
    values = ["DataEngine_Network-g1sxyw8v"]
  }
}
```
