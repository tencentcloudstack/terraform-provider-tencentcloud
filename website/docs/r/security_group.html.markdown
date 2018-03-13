---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_security_group"
sidebar_current: "docs-tencentcloud-resource-vpc-security-group-x"
description: |-
  Provides a security group resource.
---

# tencentcloud_security_group

Provides a security group resource.

## Example Usage

Basic usage:

```hcl
resource "tencentcloud_security_group" "sg" {
  name        = "test security group"
  description = "For testing security groups"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the security group. Name should be unique in each project, and no more than 60 characters.
* `description` - (Optional) The security group's description, maximum length is 100 characters.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the security group.
* `name` - The name of the security group.
* `description` - The description of the security group.
