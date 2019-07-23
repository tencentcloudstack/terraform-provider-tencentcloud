---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_security_group"
sidebar_current: "docs-tencentcloud-datasource-security_group"
description: |-
  Use this data source to query detailed information of security group.
---

# tencentcloud_security_group

Use this data source to query detailed information of security group.

~> **NOTE:** It has been deprecated and replaced by tencentcloud_security_groups.

## Example Usage

```hcl
resource "tencentcloud_security_group" "sglab" {
  name        = "mysg"
  description = "favourite sg"
  project_id  = "Default project"
}
data "tencentcloud_security_group" "sglab" {
  security_group_id = "${tencentcloud_security_group.sglab.id}"
}
```

## Argument Reference

The following arguments are supported:

* `security_group_id` - (Required) ID of the security group to be queried.
* `description` - (Optional) Description of the security group.
* `name` - (Optional) Name of the security group to be queried.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `be_associate_count` - Number of security group binding resources.
* `create_time` - Creation time of security group.


