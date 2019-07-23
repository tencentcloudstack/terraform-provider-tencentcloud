---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_security_group"
sidebar_current: "docs-tencentcloud-resource-security_group"
description: |-
  Provide a resource to create security group.
---

# tencentcloud_security_group

Provide a resource to create security group.

## Example Usage

```hcl
resource "tencentcloud_security_group" "sglab" {
  name        = "mysg"
  description = "favourite sg"
  project_id  = 0
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the security group to be queried.
* `description` - (Optional) Description of the security group.
* `project_id` - (Optional, ForceNew) Project ID of the security group.


## Import

Security group can be imported using the id, e.g.

```
  $ terraform import tencentcloud_security_group.sglab sg-ey3wmiz1
```

