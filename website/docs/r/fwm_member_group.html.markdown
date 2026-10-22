---
subcategory: "Cloud Firewall Management(FWM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_fwm_member_group"
sidebar_current: "docs-tencentcloud-resource-fwm_member_group"
description: |-
  Provides a resource to manage which FWM (Cloud Firewall Management) organization account group a member belongs to.
---

# tencentcloud_fwm_member_group

Provides a resource to manage which FWM (Cloud Firewall Management) organization account group a member belongs to.

~> NOTE: When `group_id` is non-empty, the member must have been managed by cloud firewall (a `tencentcloud_fwm_cfw_managed` resource with `cfw_managed = 1`), so the example below uses `depends_on` to ensure the member is managed first. There is a reverse coupling with `tencentcloud_fwm_cfw_managed`: disabling management (cfw_managed 1->0) will automatically move the member out of the group on the server side. If you manage both resources at the same time, set `group_id` to empty before disabling management, otherwise unexpected state drift may occur.

## Example Usage

### Move a managed member into an account group

```hcl
resource "tencentcloud_fwm_cfw_managed" "example" {
  member_id   = "mem-xxxx"
  cfw_managed = 1
}

resource "tencentcloud_fwm_account_group" "example" {
  group_name = "tf-example-group"
}

resource "tencentcloud_fwm_member_group" "example" {
  member_id = "mem-xxxx"
  app_id    = "1234567890"
  group_id  = tencentcloud_fwm_account_group.example.id

  depends_on = [tencentcloud_fwm_cfw_managed.example]
}
```

### Move a member out of an account group

```hcl
resource "tencentcloud_fwm_member_group" "example" {
  member_id = "mem-xxxx"
  app_id    = "1234567890"
  group_id  = ""
}
```

## Argument Reference

The following arguments are supported:

* `app_id` - (Required, String, ForceNew) AppId of the FWM organization member (numeric string), used to locate the member because the cloud API filter does not support MemberId.
* `member_id` - (Required, String) Member ID of the FWM organization member, such as `mem-xxxx`.
* `group_id` - (Optional, String) Account group ID that the member belongs to. Empty string means the member is not grouped (moved out). A member can belong to at most one group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Read: [DescribeOrganMembers](https://www.tencentcloud.com/search/?wd=DescribeOrganMembers&type=document&lang=en)
* Update: [ModifyMemberGroup](https://www.tencentcloud.com/search/?wd=ModifyMemberGroup&type=document&lang=en)

## Import

FWM member group resource can be imported using the two-segment composite id, e.g. `app_id#member_id`, where the first segment is the member AppId used to locate the member and the second segment is the member MemberId.

```sh
terraform import tencentcloud_fwm_member_group.example 1234567890#mem-xxxx
```

