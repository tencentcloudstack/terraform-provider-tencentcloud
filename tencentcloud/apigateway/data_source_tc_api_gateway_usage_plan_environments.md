Used to query the environment list bound by the plan.

Example Usage

```hcl
resource "tencentcloud_api_gateway_usage_plan" "plan" {
	usage_plan_name         = "my_plan"
	usage_plan_desc         = "nice plan"
	max_request_num         = 100
	max_request_num_pre_sec = 10
}

resource "tencentcloud_api_gateway_service" "service" {
	service_name = "niceservice"
	protocol     = "http&https"
	service_desc = "your nice service"
	net_type     = ["INNER", "OUTER"]
	ip_version   = "IPv4"
}

resource "tencentcloud_api_gateway_usage_plan_attachment" "attach_service" {
	usage_plan_id  = tencentcloud_api_gateway_usage_plan.plan.id
	service_id     = tencentcloud_api_gateway_service.service.id
	environment    = "test"
	bind_type      = "SERVICE"
}

data "tencentcloud_api_gateway_usage_plan_environments" "environment_test" {
	usage_plan_id = tencentcloud_api_gateway_usage_plan_attachment.attach_service.usage_plan_id
	bind_type     = "SERVICE"
}
```