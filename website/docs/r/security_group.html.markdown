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
data "tencentcloud_security_group" "sglab" {
    name        = "mysg"
    description = "favourite sg"
    project_id  = "Default project"
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Required) Description of the security group.
* `name` - (Required) Name of the security group to be queried.
* `project_id` - (Optional, ForceNew) Project ID of the security group.


