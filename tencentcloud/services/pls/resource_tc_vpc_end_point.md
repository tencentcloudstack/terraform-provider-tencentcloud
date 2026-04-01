Provides a resource to create a VPC end point

Example Usage

```hcl
resource "tencentcloud_vpc_end_point" "example" {
  vpc_id               = "vpc-391sv4w3"
  subnet_id            = "subnet-ljyn7h30"
  end_point_name       = "tf-example"
  end_point_service_id = "vpcsvc-69y13tdb"
  end_point_vip        = "10.0.2.1"

  security_group_id    = "sg-ghvp9djf"

  security_groups_ids = [
    "sg-ghvp9djf",
    "sg-if748odn",
    "sg-3k7vtgf7",
  ]

  tags = {
    env     = "test"
    project = "terraform"
  }

  ip_address_type = "Ipv4"
}
```

Argument Reference

The following arguments are supported:

* `vpc_id` - (Required, ForceNew) ID of the VPC.
* `subnet_id` - (Required, ForceNew) ID of the subnet.
* `end_point_name` - (Required) Name of the endpoint.
* `end_point_service_id` - (Required, ForceNew) ID of the endpoint service.
* `end_point_vip` - (Optional, ForceNew) VIP of the endpoint IP.
* `security_group_id` - (Optional, ForceNew) ID of the security group.
* `security_groups_ids` - (Optional) Ordered security groups associated with the endpoint.
* `tags` - (Optional) Tags of the VPC endpoint.
* `ip_address_type` - (Optional, ForceNew) IP address type. Valid values are `Ipv4` and `Ipv6`. Default is `Ipv4`.

Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `end_point_owner` - APPID.
* `state` - State of the endpoint.
* `create_time` - Create time.
* `cdc_id` - CDC instance ID.

Import

VPC end point can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_end_point.example vpce-ntv3vy9k
```