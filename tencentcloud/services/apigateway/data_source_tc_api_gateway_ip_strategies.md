Use this data source to query API gateway IP strategy.

Example Usage

```hcl
resource "tencentcloud_api_gateway_service" "service" {
  service_name = "ck"
  protocol     = "http&https"
  service_desc = "your nice service"
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
}

resource "tencentcloud_api_gateway_ip_strategy" "test"{
	service_id      = tencentcloud_api_gateway_service.service.id
	strategy_name	= "tf_test"
	strategy_type	= "BLACK"
	strategy_data	= "9.9.9.9"
}

data "tencentcloud_api_gateway_ip_strategies" "id" {
	service_id = tencentcloud_api_gateway_ip_strategy.test.service_id
}

data "tencentcloud_api_gateway_ip_strategies" "name" {
    service_id = tencentcloud_api_gateway_ip_strategy.test.service_id
	strategy_name = tencentcloud_api_gateway_ip_strategy.test.strategy_name
}
```