Provides a resource to create a cls alarm notice

Example Usage

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

Import

cls alarm notice can be imported using the id, e.g.

```
terraform import tencentcloud_cls_alarm_notice.example notice-19076f96-0f9a-4206-b308-b478737cab66
```