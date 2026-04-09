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
  name                = "tf-example"
  jump_domain         = "https://console.cloud.tencent.com"
  deliver_status      = 2
  alarm_shield_status = 2
  callback_prioritize = true
  notice_rules {
    escalate = true
    interval = 10
    rule = jsonencode(
      {
        Children = [
          {
            Children = [
              {
                Type  = "Compare"
                Value = "In"
              },
              {
                Type = "Value"
                Value = jsonencode(
                  [
                    1,
                  ]
                )
              },
            ]
            Type  = "Condition"
            Value = "NotifyType"
          },
          {
            Children = [
              {
                Type  = "Compare"
                Value = "In"
              },
              {
                Type = "Value"
                Value = jsonencode(
                  [
                    0,
                    2,
                  ]
                )
              },
            ]
            Type  = "Condition"
            Value = "Level"
          },
        ]
        Type  = "Operation"
        Value = "AND"
      }
    )
    type = 1

    escalate_notices {
      escalate = true
      interval = 10
      type     = 1

      notice_receivers {
        end_time          = "23:59:59"
        index             = 1
        notice_content_id = "Default-zh"
        receiver_channels = [
          "Phone",
          "Sms",
        ]
        receiver_ids = [
          19284382,
        ]
        receiver_type = "Uin"
        start_time    = "00:00:00"
      }
    }
    escalate_notices {
      escalate = false
      interval = 10
      type     = 1

      notice_receivers {
        end_time          = "23:59:59"
        index             = 1
        notice_content_id = "Default-en"
        receiver_channels = [
          "Email",
          "Phone",
          "Sms",
        ]
        receiver_ids = [
          19284382,
        ]
        receiver_type = "Uin"
        start_time    = "00:00:00"
      }
    }

    notice_receivers {
      end_time          = "23:59:59"
      index             = 1
      notice_content_id = "Default-en"
      receiver_channels = [
        "Sms",
      ]
      receiver_ids = [
        19284382,
      ]
      receiver_type = "Uin"
      start_time    = "00:00:00"
    }
  }

  deliver_config {
    region   = "ap-guangzhou"
    topic_id = "898016cf-7e17-426f-9167-9b56fcfc603e"
    scope    = 0
  }

  tags = {
    createdBy = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Alarm notice name.
* `alarm_shield_status` - (Optional, Int) Alarm shield status (no-login operation). Valid values: 1 (off), 2 (on, default).
* `callback_prioritize` - (Optional, Bool) Callback prioritize. true: use custom callback params from notice content template; false: use params from alarm policy.
* `deliver_config` - (Optional, List) Deliver log configuration. Required when deliver_status is 2.
* `deliver_status` - (Optional, Int) Deliver log switch. Valid values: 1 (off, default), 2 (on). When set to 2, deliver_config is required.
* `jump_domain` - (Optional, String) Jump domain. Must start with http:// or https://, cannot end with /.
* `notice_receivers` - (Optional, List) Notice receivers.
* `notice_rules` - (Optional, List) Notice rules (advanced mode). Mutually exclusive with type/notice_receivers/web_callbacks (simple mode).
* `tags` - (Optional, Map) Tag description list.
* `type` - (Optional, String) Notice type. Value: Trigger, Recovery, All.
* `web_callbacks` - (Optional, List) Callback info.

The `deliver_config` object supports the following:

* `region` - (Required, String) Region of the target log topic. e.g. ap-guangzhou.
* `topic_id` - (Required, String) Target log topic ID.
* `scope` - (Optional, Int) Deliver data scope. 0: all logs (default); 1: only alarm trigger and recovery logs.

The `escalate_notices` object of `notice_rules` supports the following:

* `escalate` - (Optional, Bool) Whether to continue escalating from this level. true: enable; false: disable.
* `interval` - (Optional, Int) Escalate interval in minutes. Range: [1, 14400].
* `notice_receivers` - (Optional, List) Notice receivers for this escalation level.
* `type` - (Optional, Int) Escalate condition. 1: unclaimed and unresolved (default); 2: unresolved.
* `web_callbacks` - (Optional, List) Web callbacks for this escalation level.

The `notice_receivers` object of `escalate_notices` supports the following:

* `receiver_channels` - (Required, Set) Receiver channels, Value: Email, Sms, WeChat, Phone.
* `receiver_ids` - (Required, Set) Receiver id list.
* `receiver_type` - (Required, String) Receiver type, Uin or Group.
* `end_time` - (Optional, String) End time allowed to receive messages.
* `index` - (Optional, Int) Index. The input parameter is invalid, but the output parameter is valid.
* `notice_content_id` - (Optional, String) Notice content ID.
* `start_time` - (Optional, String) Start time allowed to receive messages.

The `notice_receivers` object of `notice_rules` supports the following:

* `receiver_channels` - (Required, Set) Receiver channels, Value: Email, Sms, WeChat, Phone.
* `receiver_ids` - (Required, Set) Receiver id list.
* `receiver_type` - (Required, String) Receiver type, Uin or Group.
* `end_time` - (Optional, String) End time allowed to receive messages.
* `index` - (Optional, Int) Index. The input parameter is invalid, but the output parameter is valid.
* `notice_content_id` - (Optional, String) Notice content ID.
* `start_time` - (Optional, String) Start time allowed to receive messages.

The `notice_receivers` object supports the following:

* `receiver_channels` - (Required, Set) Receiver channels, Value: Email, Sms, WeChat, Phone.
* `receiver_ids` - (Required, Set) Receiver id list.
* `receiver_type` - (Required, String) Receiver type, Uin or Group.
* `end_time` - (Optional, String) End time allowed to receive messages.
* `index` - (Optional, Int) Index. The input parameter is invalid, but the output parameter is valid.
* `notice_content_id` - (Optional, String) Notice content ID.
* `start_time` - (Optional, String) Start time allowed to receive messages.

The `notice_rules` object supports the following:

* `escalate_notices` - (Optional, List) Alarm escalate notice chain, ordered from level 1 to level 5 (max). Each element represents the next escalation level.
* `escalate` - (Optional, Bool) Alarm escalate switch. true: enable; false: disable (default).
* `interval` - (Optional, Int) Alarm escalate interval in minutes. Range: [1, 14400].
* `notice_receivers` - (Optional, List) Notice receivers for this rule.
* `rule` - (Optional, String) Matching rule JSON string.
* `type` - (Optional, Int) Alarm escalate condition. 1: unclaimed and unresolved (default); 2: unresolved.
* `web_callbacks` - (Optional, List) Web callbacks for this rule.

The `web_callbacks` object of `escalate_notices` supports the following:

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

The `web_callbacks` object of `notice_rules` supports the following:

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

