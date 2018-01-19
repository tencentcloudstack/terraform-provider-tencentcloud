---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_security_group"
sidebar_current: "docs-tencentcloud-datasource-security-group"
description: |-
  Provides details about a specific Security Group.
---

# tencentcloud_security_group

`tencentcloud_security_group` provides details about a specific Security Group.

## Example Usage

Basic usage:

```hcl
variable "security_group_id" {}

data "tencentcloud_security_group" "selected" {
  id = "${var.security_group_id}"
}
```

## Argument Reference

The following arguments are supported:

* `security_group_id` - (Required) The ID of the security group.

## Attributes Reference

The following attributes are exported:

* `name` - The name of the security group.
* `description` - The description of the security group.
* `be_associate_count` - Number of associated instances.
* `create_time` - Creation time of security group, for example: 2018-01-22 17:50:21.
