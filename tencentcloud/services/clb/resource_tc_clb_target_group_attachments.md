Provides a resource to create a clb target_group_attachments

This resource supports bidirectional binding (target group binding to the load balancer, load balancer binding to the target group). When choosing either the load balancer or the target group as the binding target, up to 20 combinations can be bound at most.

Example Usage

Load balancer binding to the target group

```hcl
resource "tencentcloud_clb_target_group_attachments" "target_group_attachments" {
  load_balancer_id = "lb-phbx2420"
  associations {
    listener_id = "lbl-m2q6sp9m"
    target_group_id = "lbtg-5xunivs0"
    location_id = "loc-jjqr0ric"
  }
}

```
Target group binding to the load balancer
```hcl
resource "tencentcloud_clb_target_group_attachments" "target_group_attachments" {
  target_group_id = "lbtg-5xunivs0"
  associations { 
    listener_id = "lbl-m2q6sp9m"
    load_balancer_id = "lb-phbx2420"
    location_id = "loc-jjqr0ric"
  }
}

```