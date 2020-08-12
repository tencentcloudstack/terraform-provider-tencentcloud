---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_acl_attachment"
sidebar_current: "docs-tencentcloud-resource-vpc_acl_attachment"
description: |-
  Provide a resource to attach an existing subnet to Network ACL.
---

# tencentcloud_vpc_acl_attachment

Provide a resource to attach an existing subnet to Network ACL.

## Example Usage



## Argument Reference

The following arguments are supported:

* `acl_id` - (Required, ForceNew) Id of the attached ACL.
* `subnet_id` - (Required, ForceNew) The Subnet instance ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



