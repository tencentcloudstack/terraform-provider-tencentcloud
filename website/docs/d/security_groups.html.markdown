---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_security_groups"
sidebar_current: "docs-tencentcloud-datasource-security_groups"
description: |-
  Use this data source to query detailed information of security groups.
---

# tencentcloud_security_groups

Use this data source to query detailed information of security groups.

## Example Usage

```hcl
resource "tencentcloud_security_group" "sglab" {
  name        = "mysg"
  description = "favourite sg"
  project_id  = "Default project"
}
data "tencentcloud_security_groups" "sglab" {
  security_group_id = "${tencentcloud_security_group.sglab.id}"
  name              = "mysg"
  project_id        = "Default project"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Name of the security group to be queried. Conflict with 'security_group_id'.
* `project_id` - (Optional) Project ID of the security group. Conflict with 'security_group_id'.
* `security_group_id` - (Optional) ID of the security group to be queried. Conflict with 'name' and 'project_id'.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `security_groups` - Information list of security group.
  * `be_associate_count` - Number of security group binding resources.
  * `create_time` - Creation time of security group.
  * `description` - Description of the security group.
  * `id` - Inquired ID of the security group.
  * `name` - Inquired name of the security group.
  * `project_id` - Inquired project ID of the security group.


