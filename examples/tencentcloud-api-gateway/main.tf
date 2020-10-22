resource "tencentcloud_api_gateway_api_key" "test" {
  secret_name = var.secret_name
  status      = var.status
}

resource "tencentcloud_api_gateway_service" "service" {
  service_name = var.service_name
  protocol     = var.protocol
  service_desc = var.service_desc
  net_type     = ["INNER", "OUTER"]
  ip_version   = var.ip_version
  release_limit    = 100
  pre_limit        = 100
  test_limit       = 100
}

resource "tencentcloud_api_gateway_api" "api" {
    service_id            = tencentcloud_api_gateway_service.service.id
    api_name              = var.api_name
    api_desc              = var.api_desc
    auth_type             = "SECRET"
    protocol              = "HTTP"
    enable_cors           = true
    request_config_path   = "/user/info"
    request_config_method = "POST"

    request_parameters {
    	name          = "email"
        position      = "QUERY"
        type          = "string"
        desc          = "your email please?"
        default_value = "tom@qq.com"
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
    	code           = 10
        msg            = "system error"
       	desc           = "system error code"
       	converted_code = -10
        need_convert   = true
    }

    release_limit    = 100
    pre_limit        = 100
    test_limit       = 100
}

resource "tencentcloud_api_gateway_custom_domain" "foo" {
	service_id         = "service-ohxqslqe"
	sub_domain         = "tic-test.dnsv1.com"
	protocol           = "http"
	net_type           = "OUTER"
	is_default_mapping = "false"
	default_domain     = "service-ohxqslqe-1259649581.gz.apigw.tencentcs.com"
	path_mappings      = ["/good#test","/root#release"]
}

resource "tencentcloud_api_gateway_ip_strategy" "test"{
	service_id    = tencentcloud_api_gateway_service.service.id
	strategy_name	= var.strategy_name
	strategy_type	= "BLACK"
	strategy_data	= "9.9.9.9"
}

resource "tencentcloud_api_gateway_api_key_attachment" "attach" {
  api_key_id    = tencentcloud_api_gateway_api_key.test.id
  usage_plan_id = tencentcloud_api_gateway_usage_plan.plan.id
}

resource "tencentcloud_api_gateway_usage_plan" "plan" {
	usage_plan_name         = var.usage_plan_name
	usage_plan_desc         = var.usage_plan_desc
	max_request_num         = 100
	max_request_num_pre_sec = 10
}

resource "tencentcloud_api_gateway_usage_plan_attachment" "attach_service" {
	usage_plan_id  = tencentcloud_api_gateway_usage_plan.plan.id
	service_id     = tencentcloud_api_gateway_service.service.id
	environment    = "test"
	bind_type      = "SERVICE"
}

resource "tencentcloud_api_gateway_service_release" "service" {
  service_id       = tencentcloud_api_gateway_throttling_api.service.service_id
  environment_name = "release"
  release_desc     = var.release_desc
}

resource "tencentcloud_api_gateway_strategy_attachment" "test"{
   service_id       = tencentcloud_api_gateway_service_release.service.service_id
   strategy_id      = tencentcloud_api_gateway_ip_strategy.test.strategy_id
   environment_name = "release"
   bind_api_id      = tencentcloud_api_gateway_api.api.id
}

data "tencentcloud_api_gateway_api_keys" "name" {
  secret_name = tencentcloud_api_gateway_api_key.test.secret_name
}

data "tencentcloud_api_gateway_api_keys" "id" {
  api_key_id = tencentcloud_api_gateway_api_key.test.id
}

data "tencentcloud_api_gateway_apis" "id" {
  service_id = tencentcloud_api_gateway_service.service.id
  api_id     = tencentcloud_api_gateway_api.api.id
}

data "tencentcloud_api_gateway_apis" "name" {
  service_id = tencentcloud_api_gateway_service.service.id
  api_name   = tencentcloud_api_gateway_api.api.api_name
}

data "tencentcloud_api_gateway_customer_domains" "id" {
	service_id = tencentcloud_api_gateway_custom_domain.foo.service_id
}

data "tencentcloud_api_gateway_ip_strategies" "id" {
	service_id = tencentcloud_api_gateway_ip_strategy.test.service_id
}

data "tencentcloud_api_gateway_ip_strategies" "name" {
    service_id = tencentcloud_api_gateway_ip_strategy.test.service_id
	strategy_name = tencentcloud_api_gateway_ip_strategy.test.strategy_name
}

data "tencentcloud_api_gateway_services" "name" {
    service_name = tencentcloud_api_gateway_service.service.service_name
}

data "tencentcloud_api_gateway_services" "id" {
    service_id = tencentcloud_api_gateway_service.service.id
}

data "tencentcloud_api_gateway_throttling_apis" "id" {
    service_id = tencentcloud_api_gateway_throttling_api.service.service_id
}

data "tencentcloud_api_gateway_throttling_apis" "foo" {
	service_id        = tencentcloud_api_gateway_throttling_api.service.service_id
	environment_names = ["release", "test"]
}

data "tencentcloud_api_gateway_throttling_services" "id" {
    service_id = tencentcloud_api_gateway_throttling_service.service.service_id
}

data "tencentcloud_api_gateway_usage_plan_environments" "environment_test" {
	usage_plan_id = tencentcloud_api_gateway_usage_plan_attachment.attach_service.usage_plan_id
	bind_type     = "SERVICE"
}

data "tencentcloud_api_gateway_usage_plans" "name" {
  usage_plan_name = tencentcloud_api_gateway_usage_plan.plan.usage_plan_name
}

data "tencentcloud_api_gateway_usage_plans" "id" {
  usage_plan_id = tencentcloud_api_gateway_usage_plan.plan.id
}