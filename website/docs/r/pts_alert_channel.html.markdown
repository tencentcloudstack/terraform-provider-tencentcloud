---
subcategory: "Performance Testing Service(PTS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_pts_alert_channel"
sidebar_current: "docs-tencentcloud-resource-pts_alert_channel"
description: |-
  Provides a resource to create a pts alert_channel
---

# tencentcloud_pts_alert_channel

Provides a resource to create a pts alert_channel

~> **NOTE:** Modification is not currently supported, please go to the console to modify.

## Example Usage

```hcl
resource "tencentcloud_pts_alert_channel" "alert_channel" {
  notice_id       = ""
  project_id      = ""
  amp_consumer_id = ""
}
```

## Argument Reference

The following arguments are supported:

* `notice_id` - (Required, String) Notice ID.
* `project_id` - (Required, String) Project ID.
* `amp_consumer_id` - (Optional, String) AMP Consumer ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `app_id` - App ID Note: this field may return null, indicating that a valid value cannot be obtained.
* `created_at` - Creation time Note: this field may return null, indicating that a valid value cannot be obtained.
* `status` - Status Note: this field may return null, indicating that a valid value cannot be obtained.
* `sub_account_uin` - Sub-user ID Note: this field may return null, indicating that a valid value cannot be obtained.
* `uin` - User ID Note: this field may return null, indicating that a valid value cannot be obtained.
* `updated_at` - Update time Note: this field may return null, indicating that a valid value cannot be obtained.


## Import

pts alert_channel can be imported using the id, e.g.
```
$ terraform import tencentcloud_pts_alert_channel.alert_channel alertChannel_id
```

