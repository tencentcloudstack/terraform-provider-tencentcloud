Provides a resource to create a VPC end point service

Example Usage

```hcl
resource "tencentcloud_vpc_end_point_service" "example" {
  end_point_service_name = "tf-example"
  vpc_id                 = "vpc-9r35gtih"
  auto_accept_flag       = false
  service_type           = "CLB"
  service_instance_id    = "lb-jvb31e26"
}
```

Import

VPC end point service can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_end_point_service.example vpcsvc-l770dxs5
```