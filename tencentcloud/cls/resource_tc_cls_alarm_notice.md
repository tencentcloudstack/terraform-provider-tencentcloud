Provides a resource to create a cls alarm_notice

Example Usage

```hcl
resource "tencentcloud_cls_alarm_notice" "alarm_notice" {
  name = "terraform-alarm-notice-test"
  tags = {
    "createdBy" = "terraform"
  }
  type = "All"

  notice_receivers {
    index             = 0
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

Import

cls alarm_notice can be imported using the id, e.g.

```
terraform import tencentcloud_cls_alarm_notice.alarm_notice alarm_notice_id
```