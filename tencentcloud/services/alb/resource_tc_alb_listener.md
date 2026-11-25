Provides a resource to manage an ALB listener.

Example Usage

```hcl
resource "tencentcloud_alb_listener" "example" {
  load_balancer_id  = "alb-12345678"
  listener_protocol = "HTTPS"
  listener_port     = 443
  listener_name     = "tf-example-listener"

  default_actions {
    type = "TargetGroup"
    target_group_config {
      target_groups {
        target_group_id = "lbtg-12345678"
        weight          = 10
      }
    }
  }

  tags {
    tag_key   = "env"
    tag_value = "tf-example"
  }
}
```

Import

ALB listener can be imported using the id, e.g.

```
terraform import tencentcloud_alb_listener.example lst-12345678
```
