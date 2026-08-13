---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_message_receiver"
sidebar_current: "docs-tencentcloud-resource-cam_message_receiver"
description: |-
  Provides a resource to create a CAM message receiver
---

# tencentcloud_cam_message_receiver

Provides a resource to create a CAM message receiver

~> **NOTE:** For security reasons, the CAM will return the `email` and `phone_number` parameter values in encrypted form. Please use the `ignore_changes` function in Terraform's `lifecycle` to include these two parameters.

## Example Usage

```hcl
resource "tencentcloud_cam_message_receiver" "example" {
  name         = "tf-example"
  remark       = "remark."
  country_code = "86"
  phone_number = "18123456789"
  email        = "demo@qq.com"

  lifecycle {
    ignore_changes = [email, phone_number]
  }
}
```

## Argument Reference

The following arguments are supported:

* `country_code` - (Required, String, ForceNew) The international area code for mobile phone numbers is 86 for domestic areas.
* `email` - (Required, String, ForceNew) Email address, for example: 57*****@qq.com.
* `name` - (Required, String, ForceNew) Username of the message recipient.
* `phone_number` - (Required, String, ForceNew) Mobile phone number, for example: 132****2492.
* `remark` - (Optional, String, ForceNew) Recipient's notes.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `email_flag` - Whether the email is verified.
* `is_receiver_owner` - Whether it is the primary contact person.
* `phone_flag` - Whether the mobile phone number is verified.
* `uid` - UID.
* `uin` - Account uin.
* `wechat_flag` - Whether WeChat is allowed to receive notifications.


## Import

CAM message receiver can be imported using the id, e.g.

```
terraform import tencentcloud_cam_message_receiver.example tf-example
```

