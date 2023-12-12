Use this resource to create dayu DDoS policy case

~> **NOTE:** when a dayu DDoS policy case is created, there will be a dayu DDoS policy created with the same prefix name in the same time. This resource only supports Anti-DDoS of type `bgp`, `bgp-multip` and `bgpip`. One Anti-DDoS resource can only has one DDoS policy case resource. When there is only one Anti-DDoS resource and one policy case, those two resource will be bind automatically.

Example Usage

```hcl
resource "tencentcloud_dayu_ddos_policy_case" "foo" {
  resource_type       = "bgpip"
  name                = "tf_test_policy_case"
  platform_types      = ["PC", "MOBILE"]
  app_type            = "WEB"
  app_protocols       = ["tcp", "udp"]
  tcp_start_port      = "1000"
  tcp_end_port        = "2000"
  udp_start_port      = "3000"
  udp_end_port        = "4000"
  has_abroad          = "yes"
  has_initiate_tcp    = "yes"
  has_initiate_udp    = "yes"
  peer_tcp_port       = "1111"
  peer_udp_port       = "3333"
  tcp_footprint       = "511"
  udp_footprint       = "500"
  web_api_urls        = ["abc.com", "test.cn/aaa.png"]
  min_tcp_package_len = "1000"
  max_tcp_package_len = "1200"
  min_udp_package_len = "1000"
  max_udp_package_len = "1200"
  has_vpn             = "yes"
}
```