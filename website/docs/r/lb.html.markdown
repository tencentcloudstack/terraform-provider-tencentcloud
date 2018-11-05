---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lb"
sidebar_current: "docs-tencentcloud-resource-lb-x"
description: |-
  Provides a Load Balancer resource.
---

# tencentcloud_lb

Provides a Load Balancer resource.

## Example Usage

Basic usage:

```hcl
resource "tencentcloud_lb" "classic" {
  type       = 2
  forward    = 0
  name       = "tf-test-classic"
  project_id = 0
}
```

## Argument Reference

The following arguments are supported:

* `type` - (Required) (Int) The network type of the LB, 2 for public network and 1 for private network.
* `forward` - (Optional) (Int) The type of the LB, 0 for classic and 1 for application.
* `name` - (Optional) The name of the LB.
* `vpc_id` - (Optional) The VPC ID of the LB, unspecified or 0 stands for CVM basic network.
* `project_id` - (Optional) The project id of the LB, unspecified or 0 stands for default project.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The status of the LB.
