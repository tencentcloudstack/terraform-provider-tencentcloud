Use this data source to query detailed information of VPN gateways.

Example Usage

```hcl
data "tencentcloud_nat_gateway_snats" "snat" {
  nat_gateway_id     = tencentcloud_nat_gateway.my_nat.id
  subnet_id          = tencentcloud_nat_gateway_snat.my_subnet.id
  public_ip_addr     = ["50.29.23.234"]
  description        = "snat demo"
  result_output_file = "./snat.txt"
}
```