Provides a resource to create a dc_gateway_attachment

Example Usage

```hcl
resource "tencentcloud_dc_gateway_attachment" "dc_gateway_attachment" {
  vpc_id = "vpc-4h9v4mo3"
  nat_gateway_id = "nat-7kanjc6y"
  direct_connect_gateway_id = "dcg-dmbhf7jf"
}
```

Import

dc_gateway_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_dc_gateway_attachment.dc_gateway_attachment vpcId#dcgId#ngId
```