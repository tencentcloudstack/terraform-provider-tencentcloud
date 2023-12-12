Use this resource to create API gateway usage plan.

Example Usage

```hcl
resource "tencentcloud_api_gateway_usage_plan" "example" {
  usage_plan_name         = "tf_example"
  usage_plan_desc         = "desc."
  max_request_num         = 100
  max_request_num_pre_sec = 10
}
```

Import

API gateway usage plan can be imported using the id, e.g.

```
$ terraform import tencentcloud_api_gateway_usage_plan.plan usagePlan-gyeafpab
```