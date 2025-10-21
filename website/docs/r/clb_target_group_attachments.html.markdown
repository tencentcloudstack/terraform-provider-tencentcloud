---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_target_group_attachments"
sidebar_current: "docs-tencentcloud-resource-clb_target_group_attachments"
description: |-
  Provides a resource to create a clb target_group_attachments
---

# tencentcloud_clb_target_group_attachments

Provides a resource to create a clb target_group_attachments

This resource supports bidirectional binding (target group binding to the load balancer, load balancer binding to the target group). When choosing either the load balancer or the target group as the binding target, up to 20 combinations can be bound at most.

## Example Usage

### Load balancer binding to the target group

```hcl
resource "tencentcloud_clb_target_group_attachments" "target_group_attachments" {
  load_balancer_id = "lb-phbx2420"
  associations {
    listener_id     = "lbl-m2q6sp9m"
    target_group_id = "lbtg-5xunivs0"
    location_id     = "loc-jjqr0ric"
  }
}
```

### Target group binding to the load balancer

```hcl
resource "tencentcloud_clb_target_group_attachments" "target_group_attachments" {
  target_group_id = "lbtg-5xunivs0"
  associations {
    listener_id      = "lbl-m2q6sp9m"
    load_balancer_id = "lb-phbx2420"
    location_id      = "loc-jjqr0ric"
  }
}
```

## Argument Reference

The following arguments are supported:

* `associations` - (Required, Set, ForceNew) Association array, the combination cannot exceed 20.
* `load_balancer_id` - (Optional, String, ForceNew) CLB instance ID, (load_balancer_id and target_group_id require at least one).
* `target_group_id` - (Optional, String, ForceNew) Target group ID, (load_balancer_id and target_group_id require at least one).

The `associations` object supports the following:

* `listener_id` - (Optional, String, ForceNew) Listener ID.
* `load_balancer_id` - (Optional, String, ForceNew) CLB instance ID, when the binding target is target group, load_balancer_id in associations is required.
* `location_id` - (Optional, String, ForceNew) Forwarding rule ID.
* `target_group_id` - (Optional, String, ForceNew) Target group ID, when the binding target is clb, the target_group_id in associations is required.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



