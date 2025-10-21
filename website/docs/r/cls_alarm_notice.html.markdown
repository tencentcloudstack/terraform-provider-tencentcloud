---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_alarm_notice"
sidebar_current: "docs-tencentcloud-resource-cls_alarm_notice"
description: |-
  Provides a resource to create a cls alarm notice
---

# tencentcloud_cls_alarm_notice

Provides a resource to create a cls alarm notice

## Example Usage

```hcl
resource "tencentcloud_cls_alarm_notice" "example" {
  name = "tf-example"
  type = "All"

  notice_receivers {
    receiver_type = "Uin"
    receiver_ids = [
      100037718139,
    ]
    receiver_channels = [
      "Email",
      "Sms",
    ]
    notice_content_id = "noticetemplate-b417f32a-bdf9-46c5-933e-28c23cd7a6b7"
    start_time        = "00:00:00"
    end_time          = "23:59:59"
  }

  web_callbacks {
    callback_type     = "Http"
    url               = "example.com"
    method            = "POST"
    notice_content_id = "noticetemplate-b417f32a-bdf9-46c5-933e-28c23cd7a6b7"
    remind_type       = 1
  }

  tags = {
    createdBy = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Alarm notice name.
* `type` - (Required, String) Notice type. Value: Trigger, Recovery, All.
* `notice_receivers` - (Optional, List) Notice receivers.
* `tags` - (Optional, Map) Tag description list.
* `web_callbacks` - (Optional, List) Callback info.

The `notice_receivers` object supports the following:

* `receiver_channels` - (Required, Set) Receiver channels, Value: Email, Sms, WeChat, Phone.
* `receiver_ids` - (Required, Set) Receiver id list.
* `receiver_type` - (Required, String) Receiver type, Uin or Group.
* `end_time` - (Optional, String) End time allowed to receive messages.
* `index` - (Optional, Int) Index. The input parameter is invalid, but the output parameter is valid.
* `notice_content_id` - (Optional, String) Notice content ID.
* `start_time` - (Optional, String) Start time allowed to receive messages.

The `web_callbacks` object supports the following:

* `callback_type` - (Required, String) Callback type, Values: Http, WeCom, DingTalk, Lark.
* `url` - (Required, String) Callback url.
* `body` - (Optional, String, **Deprecated**) This parameter is deprecated. Please use `notice_content_id`. Request body.
* `headers` - (Optional, Set, **Deprecated**) This parameter is deprecated. Please use `notice_content_id`. Request headers.
* `index` - (Optional, Int) Index. The input parameter is invalid, but the output parameter is valid.
* `method` - (Optional, String) Method, POST or PUT.
* `mobiles` - (Optional, Set) Telephone list.
* `notice_content_id` - (Optional, String) Notice content ID.
* `remind_type` - (Optional, Int) Remind type. 0: Do not remind; 1: Specified person; 2: Everyone.
* `user_ids` - (Optional, Set) User ID list.
* `web_callback_id` - (Optional, String) Integration configuration ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cls alarm notice can be imported using the id, e.g.

```
terraform import tencentcloud_cls_alarm_notice.example notice-19076f96-0f9a-4206-b308-b478737cab66
```

