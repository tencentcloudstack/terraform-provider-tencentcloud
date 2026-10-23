---
subcategory: "Cloud Firewall(FWM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_fwm_cfw_share_role"
sidebar_current: "docs-tencentcloud-resource-fwm_cfw_share_role"
description: |-
  Provides a resource to create a fwm cfw_share_role, used to manage the cloud firewall (fwm) share role (sharer/user/none) expectation state of an organization member.
---

# tencentcloud_fwm_cfw_share_role

Provides a resource to create a fwm cfw_share_role, used to manage the cloud firewall (fwm) share role (sharer/user/none) expectation state of an organization member.

## Example Usage

```hcl
resource "tencentcloud_fwm_cfw_share_role" "example" {
  member_id      = "1234567890"
  app_id         = "1234567890"
  cfw_share_role = "none"
}
```

## Argument Reference

The following arguments are supported:

* `app_id` - (Required, String, ForceNew) Target member AppId, used to locate the member in Read.
* `cfw_share_role` - (Required, String) Cloud firewall share role, valid values: sharer, user, none.
* `member_id` - (Required, String) Target member ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cfw_sharer_app_id` - Sharer AppId, only available when the role is user.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Read: [DescribeOrganMembers](https://www.tencentcloud.com/search/?wd=DescribeOrganMembers&type=document&lang=en)
* Update:
  - [DescribeOrganMembers](https://www.tencentcloud.com/search/?wd=DescribeOrganMembers&type=document&lang=en)
  - [ModifyCfwShareRole](https://www.tencentcloud.com/search/?wd=ModifyCfwShareRole&type=document&lang=en)

## Import

fwm cfw share role can be imported using the appId#memberId, e.g.

```
terraform import tencentcloud_fwm_cfw_share_role.example 1234567890#1234567890
```

