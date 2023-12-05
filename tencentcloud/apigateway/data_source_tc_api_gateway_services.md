Use this data source to query API gateway services.

Example Usage

```hcl
resource "tencentcloud_api_gateway_service" "service" {
  service_name = "niceservice"
  protocol     = "http&https"
  service_desc = "your nice service"
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
}

data "tencentcloud_api_gateway_services" "name" {
    service_name = tencentcloud_api_gateway_service.service.service_name
}

data "tencentcloud_api_gateway_services" "id" {
    service_id = tencentcloud_api_gateway_service.service.id
}
```