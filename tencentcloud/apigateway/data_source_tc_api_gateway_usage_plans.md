Use this data source to query API gateway usage plans.

Example Usage

```hcl
resource "tencentcloud_api_gateway_usage_plan" "plan" {
  usage_plan_name         = "my_plan"
  usage_plan_desc         = "nice plan"
  max_request_num         = 100
  max_request_num_pre_sec = 10
}

data "tencentcloud_api_gateway_usage_plans" "name" {
  usage_plan_name = tencentcloud_api_gateway_usage_plan.plan.usage_plan_name
}

data "tencentcloud_api_gateway_usage_plans" "id" {
  usage_plan_id = tencentcloud_api_gateway_usage_plan.plan.id
}
```