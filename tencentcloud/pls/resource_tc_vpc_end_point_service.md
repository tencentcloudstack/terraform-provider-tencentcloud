Provides a resource to create a vpc end_point_service

Example Usage

```hcl
resource "tencentcloud_vpc_end_point_service" "end_point_service" {
  vpc_id = "vpc-391sv4w3"
  end_point_service_name = "terraform-endpoint-service"
  auto_accept_flag = false
  service_instance_id = "lb-o5f6x7ke"
  service_type = "CLB"
}
```

Import

vpc end_point_service can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_end_point_service.end_point_service end_point_service_id
```