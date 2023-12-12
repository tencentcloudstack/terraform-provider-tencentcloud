Use this resource to attach API gateway usage plan to service.

~> **NOTE:** If the `auth_type` parameter of API is not `SECRET`, it cannot be bound access key.

Example Usage

Normal creation

```hcl
resource "tencentcloud_api_gateway_usage_plan" "example" {
  usage_plan_name         = "tf_example"
  usage_plan_desc         = "desc."
  max_request_num         = 100
  max_request_num_pre_sec = 10
}

resource "tencentcloud_api_gateway_service" "example" {
  service_name = "tf_example"
  protocol     = "http&https"
  service_desc = "desc."
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
}

resource "tencentcloud_api_gateway_api" "example" {
  service_id            = tencentcloud_api_gateway_service.example.id
  api_name              = "tf_example"
  api_desc              = "my hello api update"
  auth_type             = "SECRET"
  protocol              = "HTTP"
  enable_cors           = true
  request_config_path   = "/user/info"
  request_config_method = "POST"
  request_parameters {
    name          = "email"
    position      = "QUERY"
    type          = "string"
    desc          = "desc."
    default_value = "test@qq.com"
    required      = true
  }
  service_config_type      = "HTTP"
  service_config_timeout   = 10
  service_config_url       = "http://www.tencent.com"
  service_config_path      = "/user"
  service_config_method    = "POST"
  response_type            = "XML"
  response_success_example = "<note>success</note>"
  response_fail_example    = "<note>fail</note>"
  response_error_codes {
    code           = 500
    msg            = "system error"
    desc           = "system error code"
    converted_code = 5000
    need_convert   = true
  }
}

resource "tencentcloud_api_gateway_usage_plan_attachment" "example" {
  usage_plan_id = tencentcloud_api_gateway_usage_plan.example.id
  service_id    = tencentcloud_api_gateway_service.example.id
  environment   = "release"
  bind_type     = "API"
  api_id        = tencentcloud_api_gateway_api.example.id
}
```

Bind the key to a usage plan

```hcl
resource "tencentcloud_api_gateway_api_key" "example" {
  secret_name = "tf_example"
  status      = "on"
}

resource "tencentcloud_api_gateway_usage_plan_attachment" "example" {
  usage_plan_id = tencentcloud_api_gateway_usage_plan.example.id
  service_id    = tencentcloud_api_gateway_service.example.id
  environment   = "release"
  bind_type     = "API"
  api_id        = tencentcloud_api_gateway_api.example.id

  access_key_ids = [tencentcloud_api_gateway_api_key.example.id]
}
```

Import

API gateway usage plan attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_api_gateway_usage_plan_attachment.attach_service usagePlan-pe7fbdgn#service-kuqd6xqk#release#API#api-p8gtanvy
```