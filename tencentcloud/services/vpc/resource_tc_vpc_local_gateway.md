Provides a resource to create a vpc local_gateway

Example Usage

```hcl
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_vpc_local_gateway" "example" {
  local_gateway_name = "tf-example"
  vpc_id             = tencentcloud_vpc.vpc.id
  cdc_id             = "cluster-j9gyu1iy"
}
```

Import

vpc local_gateway can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_local_gateway.local_gateway local_gateway_id
```