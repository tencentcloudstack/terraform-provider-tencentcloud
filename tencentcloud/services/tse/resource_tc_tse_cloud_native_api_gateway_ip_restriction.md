Provides a resource to manage a TSE cloud native API gateway IP restriction.

Example Usage

```hcl
resource "tencentcloud_tse_cloud_native_api_gateway_ip_restriction" "example" {
  gateway_id       = "gateway-ddbb709b"
  source_type      = "route"
  source_id        = "route-xxxx"
  enabled          = true
  restriction_type = "whiteList"
  address_list     = ["10.0.0.0/8", "192.168.1.1"]
}
```

Import

tse cloud native api gateway ip restriction can be imported using the composite id, e.g.

```
terraform import tencentcloud_tse_cloud_native_api_gateway_ip_restriction.example gatewayId#sourceType#sourceId
```
