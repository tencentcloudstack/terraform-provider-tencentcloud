Use this data source to query detailed information of TEO multi path gateways

Example Usage

```hcl
data "tencentcloud_teo_multi_path_gateways" "example" {
  zone_id = "zone-2o1xvpmq7nn"
  filters {
    name   = "gateway-type"
    values = ["cloud"]
  }
}
```
