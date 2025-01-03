Provides a resource to create a VPC end point

Example Usage

```hcl
resource "tencentcloud_vpc_end_point" "example" {
  vpc_id               = "vpc-391sv4w3"
  subnet_id            = "subnet-ljyn7h30"
  end_point_name       = "tf-example"
  end_point_service_id = "vpcsvc-69y13tdb"
  end_point_vip        = "10.0.2.1"

  security_groups_ids = [
    "sg-ghvp9djf",
    "sg-if748odn",
    "sg-3k7vtgf7",
  ]
}
```

Import

VPC end point can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_end_point.example vpce-ntv3vy9k
```