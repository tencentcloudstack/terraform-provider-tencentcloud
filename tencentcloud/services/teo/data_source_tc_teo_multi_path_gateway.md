Use this data source to query detailed information of TEO multi-path gateways

Example Usage

Query all gateways by zone_id

```hcl
data "tencentcloud_teo_multi_path_gateway" "example" {
  zone_id = "zone-2noq7st5t3t6"
}
```

Query gateways by zone_id with filters

```hcl
data "tencentcloud_teo_multi_path_gateway" "example" {
  zone_id = "zone-2noq7st5t3t6"
  filters {
    name = "gateway-type"
    values = ["cloud"]
  }
}
```
