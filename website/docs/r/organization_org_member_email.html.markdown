---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_org_member_email"
sidebar_current: "docs-tencentcloud-resource-organization_org_member_email"
description: |-
  Provides a resource to create a organization member email
---

# tencentcloud_organization_org_member_email

Provides a resource to create a organization member email

## Example Usage

```hcl
resource "tencentcloud_organization_org_member_email" "example" {
  member_uin   = 100033118139
  email        = "example@tencent.com"
  country_code = "86"
  phone        = "18611111111"
}
```

## Argument Reference

The following arguments are supported:

* `country_code` - (Required, String) International region.
* `email` - (Required, String) Email address.
* `member_uin` - (Required, Int) Member Uin.
* `phone` - (Required, String) Phone number.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `apply_time` - Application timeNote: This field may return NULL, indicating that the valid value cannot be obtained.
* `bind_id` - Binding IDNote: This field may return NULL, indicating that the valid value cannot be obtained.
* `bind_status` - Binding status is not binding: unbound, to be activated: value, successful binding: success, binding failure: failedNote: This field may return NULL, indicating that the valid value cannot be obtained.
* `bind_time` - Binding timeNote: This field may return NULL, indicating that the valid value cannot be obtained.
* `description` - FailedNote: This field may return NULL, indicating that the valid value cannot be obtained.
* `phone_bind` - Safe mobile phone binding state is not bound: 0, has been binded: 1Note: This field may return NULL, indicating that the valid value cannot be obtained.


## Import

organization member email can be imported using the id, e.g.

```
terraform import tencentcloud_organization_org_member_email.example 100033118139#132
```

