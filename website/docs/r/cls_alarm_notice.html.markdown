---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_alarm_notice"
sidebar_current: "docs-tencentcloud-resource-cls_alarm_notice"
description: |-
  Provides a resource to create a cls alarm_notice
---

# tencentcloud_cls_alarm_notice

Provides a resource to create a cls alarm_notice

## Example Usage

```hcl
resource "tencentcloud_cls_alarm_notice" "alarm_notice" {
  name = "terraform-alarm-notice-test"
  tags = {
    "createdBy" = "terraform"
  }
  type = "All"

  notice_receivers {
    index = 0
    receiver_channels = [
      "Sms",
    ]
    receiver_ids = [
      13478043,
      15972111,
    ]
    receiver_type = "Uin"
    start_time    = "00:00:00"
    end_time      = "23:59:59"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) alarm notice name.
* `type` - (Required, String) notice type.
* `notice_receivers` - (Optional, List) notice receivers.
* `tags` - (Optional, Map) Tag description list.
* `web_callbacks` - (Optional, List) callback info.

The `notice_receivers` object supports the following:

* `receiver_channels` - (Required, Set) receiver channels, Email,Sms,WeChat or Phone.
* `receiver_ids` - (Required, Set) receiver id.
* `receiver_type` - (Required, String) receiver type, Uin or Group.
* `end_time` - (Optional, String) end time allowed to receive messages.
* `index` - (Optional, Int) index.
* `start_time` - (Optional, String) start time allowed to receive messages.

The `web_callbacks` object supports the following:

* `callback_type` - (Required, String) callback type, WeCom or Http.
* `url` - (Required, String) callback url.
* `body` - (Optional, String) abandoned.
* `headers` - (Optional, Set) abandoned.
* `index` - (Optional, Int) index.
* `method` - (Optional, String) Method, POST or PUT.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cls alarm_notice can be imported using the id, e.g.

```
terraform import tencentcloud_cls_alarm_notice.alarm_notice alarm_notice_id
```

