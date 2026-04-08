Provides a resource to create a cls alarm notice

Example Usage

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

Import

cls alarm notice can be imported using the id, e.g.

```
terraform import tencentcloud_cls_alarm_notice.example notice-19076f96-0f9a-4206-b308-b478737cab66
```
