Use this data source to query API gateway throttling services.

Example Usage

```hcl
resource "tencentcloud_api_gateway_service" "service" {
  	service_name     = "niceservice"
  	protocol         = "http&https"
  	service_desc     = "your nice service"
  	net_type         = ["INNER", "OUTER"]
	ip_version       = "IPv4"
	release_limit    = 100
	pre_limit        = 100
	test_limit       = 100
}

data "tencentcloud_api_gateway_throttling_services" "id" {
    service_id = tencentcloud_api_gateway_service.service.id
}
```