---
subcategory: "Firewall Manager(FWM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_fwm_account_group"
sidebar_current: "docs-tencentcloud-resource-fwm_account_group"
description: |-
  Provides a resource to create a FWM (Edge firewall account group management) account group
---

# tencentcloud_fwm_account_group

Provides a resource to create a FWM (Edge firewall account group management) account group

## Example Usage

```hcl
resource "tencentcloud_fwm_account_group" "example" {
  group_name  = "test-account-group"
  description = "test account group"
}
```

## Argument Reference

The following arguments are supported:

* `group_name` - (Required, String) Account group name, up to 50 characters, unique under the same admin AppId.
* `description` - (Optional, String, ForceNew) Description of the account group, only settable on create.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `can_delete` - Whether the account group can be deleted.
* `create_time` - Create time of the account group.
* `group_id` - Account group ID, e.g. acg-xxxxxxxx.
* `member_count` - Member count of the account group.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateAccountGroup](https://www.tencentcloud.com/search/?wd=CreateAccountGroup&type=document&lang=en)
* Read: [DescribeAccountGroups](https://www.tencentcloud.com/search/?wd=DescribeAccountGroups&type=document&lang=en)
* Update: [ModifyAccountGroup](https://www.tencentcloud.com/search/?wd=ModifyAccountGroup&type=document&lang=en)
* Delete: [DeleteAccountGroup](https://www.tencentcloud.com/search/?wd=DeleteAccountGroup&type=document&lang=en)

## Import

FWM account group can be imported using the GroupId, e.g.

```
terraform import tencentcloud_fwm_account_group.example acg-1a2b3c4d
```

