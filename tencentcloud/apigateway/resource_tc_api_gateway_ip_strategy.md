Use this resource to create IP strategy of API gateway.

Example Usage

```hcl
resource "tencentcloud_api_gateway_service" "service" {
  	service_name = "niceservice"
  	protocol     = "http&https"
  	service_desc = "your nice service"
  	net_type     = ["INNER", "OUTER"]
  	ip_version   = "IPv4"
}

resource "tencentcloud_api_gateway_ip_strategy" "test"{
    service_id    = tencentcloud_api_gateway_service.service.id
    strategy_name = "tf_test"
    strategy_type = "BLACK"
    strategy_data = "9.9.9.9"
}
```

Import

IP strategy of API gateway can be imported using the id, e.g.

```
$ terraform import tencentcloud_api_gateway_ip_strategy.test service-ohxqslqe#IPStrategy-q1lk8ud2
```