---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lb"
sidebar_current: "docs-tencentcloud-resource-lb-x"
description: |-
  Provides a Load Balancer resource.
---

# tencentcloud_lb

Provides a Load Balancer resource.

~> **NOTE:** It has been deprecated and replaced by `tencentcloud_clb_instance`.

## Example Usage

Basic usage:

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

* `type` - (Required)  The network type of the LB, valid choices: "OPEN", "INTERNAL".
* `forward` - (Optional) The type of the LB, valid choices: "CLASSIC", "APPLICATION".
* `name` - (Optional) The name of the LB.
* `vpc_id` - (Optional) The VPC ID of the LB, unspecified or 0 stands for CVM basic network.
* `project_id` - (Optional) The project id of the LB, unspecified or 0 stands for default project.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The status of the LB.
