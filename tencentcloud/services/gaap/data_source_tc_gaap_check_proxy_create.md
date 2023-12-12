Use this data source to query detailed information of gaap check proxy create

Example Usage

```hcl
data "tencentcloud_gaap_check_proxy_create" "check_proxy_create" {
  access_region = "Guangzhou"
  real_server_region = "Beijing"
  bandwidth = 10
  concurrent = 2
  ip_address_version = "IPv4"
  network_type = "normal"
  package_type = "Thunder"
  http3_supported = 0
}
```