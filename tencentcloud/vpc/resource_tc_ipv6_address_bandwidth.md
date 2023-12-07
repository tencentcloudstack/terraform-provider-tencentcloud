Provides a resource to create a ipv6_address_bandwidth

Example Usage

```hcl
resource "tencentcloud_ipv6_address_bandwidth" "ipv6_address_bandwidth" {
  ipv6_address                = "2402:4e00:1019:9400:0:9905:a90b:2ef0"
  internet_max_bandwidth_out = 6
  internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
#  bandwidth_package_id       = "bwp-34rfgt56"
}
```