Use this data source to query detailed information of tse gateway_routes

Example Usage

```hcl
data "tencentcloud_tse_gateway_routes" "gateway_routes" {
  gateway_id   = "gateway-ddbb709b"
  service_name = "test"
  route_name   = "keep-routes"
}
```