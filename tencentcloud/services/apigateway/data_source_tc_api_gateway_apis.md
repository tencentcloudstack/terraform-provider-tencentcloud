Use this data source to query API gateway APIs.

Example Usage

```hcl
resource "tencentcloud_api_gateway_service" "service" {
  service_name = "ck"
  protocol     = "http&https"
  service_desc = "your nice service"
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
}

resource "tencentcloud_api_gateway_api" "api" {
  service_id               = tencentcloud_api_gateway_service.service.id
  api_name                 = "hello"
  api_desc                 = "my hello api"
  auth_type                = "NONE"
  protocol                 = "HTTP"
  enable_cors              = true
  request_config_path      = "/user/info"
  request_config_method    = "GET"
  service_config_type      = "HTTP"
  service_config_timeout   = 15
  service_config_url       = "http://www.qq.com"
  service_config_path      = "/user"
  service_config_method    = "GET"
  response_type            = "HTML"
  response_success_example = "success"
  response_fail_example    = "fail"
}

data "tencentcloud_api_gateway_apis" "id" {
  service_id = tencentcloud_api_gateway_service.service.id
  api_id     = tencentcloud_api_gateway_api.api.id
}

data "tencentcloud_api_gateway_apis" "name" {
  service_id = tencentcloud_api_gateway_service.service.id
  api_name   = tencentcloud_api_gateway_api.api.api_name
}
```