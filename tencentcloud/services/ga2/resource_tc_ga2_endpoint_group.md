Provides a resource to create a GA2 (Global Accelerator 2) endpoint group.

Example Usage

```hcl
resource "tencentcloud_ga2_endpoint_group" "example" {
  global_accelerator_id = "ga2-xxxxxxxx"
  listener_id           = "lis-xxxxxxxx"
  endpoint_group_type   = "DEFAULT"

  endpoint_group_configuration {
    name                  = "tf-example"
    endpoint_group_region = "ap-guangzhou"
    description           = "tf example endpoint group"
    enable_health_check   = true
    check_type            = "HTTP"
    check_port            = "80"
    check_path            = "/"
    check_method          = "GET"
    connect_timeout       = 5000
    health_check_interval = 30
    healthy_threshold     = 3
    unhealthy_threshold   = 3
    forward_protocol      = "HTTP"

    endpoint_configurations {
      endpoint_type    = "PublicIp"
      endpoint_service = "1.1.1.1"
      weight           = 10
    }

    endpoint_configurations {
      endpoint_type    = "Domain"
      endpoint_service = "example.com"
      weight           = 20
    }
  }
}
```

Import

GA2 endpoint group can be imported using the composite ID `<global_accelerator_id>#<listener_id>#<endpoint_group_id>`, e.g.

```
terraform import tencentcloud_ga2_endpoint_group.example ga2-xxxxxxxx#lis-xxxxxxxx#eg-xxxxxxxx
```
