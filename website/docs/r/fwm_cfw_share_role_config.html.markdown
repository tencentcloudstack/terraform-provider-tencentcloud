---
subcategory: "Firewall Manager(FWM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_fwm_cfw_share_role_config"
sidebar_current: "docs-tencentcloud-resource-fwm_cfw_share_role_config"
description: |-
  Provides a resource to create a fwm cfw_share_role, used to manage the cloud firewall (fwm) share role (sharer/user/none) expectation state of an organization member.
---

# tencentcloud_fwm_cfw_share_role_config

Provides a resource to create a fwm cfw_share_role, used to manage the cloud firewall (fwm) share role (sharer/user/none) expectation state of an organization member.

## Example Usage

### Configure the share role as user

```hcl
resource "tencentcloud_fwm_cfw_share_role_config" "example" {
  cfw_share_role   = "user"
  user_member_id   = "mem-tencent-e46e42e79411caf0"
  sharer_member_id = "mem-tencent-d55a312b431535c3"
}
```

### Configure the share role as none

```hcl
resource "tencentcloud_fwm_cfw_share_role_config" "example" {
  cfw_share_role = "none"
  user_member_id = "mem-tencent-e46e42e79411caf0"
}
```

## Argument Reference

The following arguments are supported:

* `cfw_share_role` - (Required, String) Cloud firewall shared roles, with values of sharer, user, and none. Sharer -> user prohibited from direct changes, must first be set to none.
* `user_member_id` - (Required, String, ForceNew) User ID.
* `sharer_member_id` - (Optional, String) Sharer account ID, required when CfwShareRole is user.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [ModifyCfwShareRole](https://www.tencentcloud.com/search/?wd=ModifyCfwShareRole&type=document&lang=en)
* Read: [DescribeOrganMembers](https://www.tencentcloud.com/search/?wd=DescribeOrganMembers&type=document&lang=en)
* Update: [ModifyCfwShareRole](https://www.tencentcloud.com/search/?wd=ModifyCfwShareRole&type=document&lang=en)

## Import

fwm cfw share role can be imported using id, e.g.

```
terraform import tencentcloud_fwm_cfw_share_role_config.example mem-tencent-e46e42e79411caf0
```

