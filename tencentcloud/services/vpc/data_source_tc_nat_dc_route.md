Use this data source to query detailed information of vpc nat_dc_route

Example Usage

```hcl
data "tencentcloud_nat_dc_route" "nat_dc_route" {
  nat_gateway_id = "nat-gnxkey2e"
  vpc_id         = "vpc-pyyv5k3v"
}
```