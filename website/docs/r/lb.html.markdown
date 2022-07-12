---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lb"
sidebar_current: "docs-tencentcloud-resource-lb"
description: |-
  Provides a Load Balancer resource.
---

# tencentcloud_lb

Provides a Load Balancer resource.

~> **NOTE:** It has been deprecated and replaced by `tencentcloud_clb_instance`.

## Example Usage

```hcl
resource "tencentcloud_lb" "classic" {
  type       = "OPEN"
  forward    = "APPLICATION"
  name       = "tf-test-classic"
  project_id = 0
}
```

## Argument Reference

The following arguments are supported:

* `type` - (Required, String, ForceNew) The network type of the LB. Valid value: 'OPEN', 'INTERNAL'.
* `forward` - (Optional, String, ForceNew) The type of the LB. Valid value: 'CLASSIC', 'APPLICATION'.
* `name` - (Optional, String) The name of the LB.
* `project_id` - (Optional, Int, ForceNew) The project id of the LB, unspecified or 0 stands for default project.
* `vpc_id` - (Optional, String, ForceNew) The VPC ID of the LB, unspecified or 0 stands for CVM basic network.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `status` - The status of the LB.


