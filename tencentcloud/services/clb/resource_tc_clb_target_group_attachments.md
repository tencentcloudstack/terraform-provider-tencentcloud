Provides a resource to create a CLB target group attachments

~> **NOTE:** This resource supports bidirectional binding (target group binding to the load balancer, load balancer binding to the target group). When choosing either the load balancer or the target group as the binding target, up to 20 combinations can be bound at most.

Example Usage

Load balancer binding to the target group

```hcl
resource "tencentcloud_clb_target_group_attachments" "example" {
  load_balancer_id = "lb-lmgp1eis"
  associations {
    listener_id     = "lbl-jbdfcswy"
    target_group_id = "lbtg-bjosq37w"
    location_id     = "loc-bjl41tpc"
    weight          = "10"
  }
}
```

Target group binding to the load balancer

```hcl
resource "tencentcloud_clb_target_group_attachments" "example" {
  load_balancer_id = "lb-lmgp1eis"
  associations {
    listener_id      = "lbl-jbdfcswy"
    load_balancer_id = "lb-phbx2420"
    location_id      = "loc-bjl41tpc"
    weight           = "10"
  }
}
```
