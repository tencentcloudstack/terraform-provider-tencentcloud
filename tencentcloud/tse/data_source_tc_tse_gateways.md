Use this data source to query detailed information of tse gateways

Example Usage

```hcl
data "tencentcloud_tse_gateways" "gateways" {
  filters {
    name   = "GatewayId"
    values = ["gateway-ddbb709b"]
  }
}
```