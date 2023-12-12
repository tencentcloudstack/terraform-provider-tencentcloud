Use this data source to query detailed information of NAT gateways.

Example Usage

```hcl
data "tencentcloud_nat_gateways" "foo" {
  name   = "main"
  vpc_id = "vpc-xfqag"
  id     = "nat-xfaq1"
}
```