---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_security_group"
sidebar_current: "docs-tencentcloud-resource-security-group"
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

* `name` - (Optional) The name of the security group. Under the same project, the name can not be the same, can be arbitrarily named, but no t more than 60 characters.
* `description` - (Optional) The security group description, The upper limit of 100 characters.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the security group.
* `name` - The name of the security group.
* `description` - The description of the security group.
