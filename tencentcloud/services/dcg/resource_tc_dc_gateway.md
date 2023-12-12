Provides a resource to creating direct connect gateway instance.

Example Usage

```hcl
resource "tencentcloud_vpc" "main" {
  name       = "ci-vpc-instance-test"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_dc_gateway" "vpc_main" {
  name                = "ci-cdg-vpc-test"
  network_instance_id = tencentcloud_vpc.main.id
  network_type        = "VPC"
  gateway_type        = "NAT"
}
```

Import

Direct connect gateway instance can be imported, e.g.

```
$ terraform import tencentcloud_dc_gateway.instance dcg-id
```