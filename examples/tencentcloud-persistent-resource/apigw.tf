resource "tencentcloud_api_gateway_service" "service" {
  service_name = "keep_apigw_service"
  protocol     = "http&https"
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
}