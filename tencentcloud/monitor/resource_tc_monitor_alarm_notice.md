Provides a alarm notice resource for monitor.

Example Usage

```hcl
resource "tencentcloud_monitor_alarm_notice" "example" {
  name            = "test_alarm_notice"
  notice_language = "zh-CN"
  notice_type     = "ALL"

  url_notices {
    end_time   = 86399
    is_valid = 0
    start_time = 0
    url        = "https://www.mytest.com/validate"
    weekday    = [
      1,
      2,
      3,
      4,
      5,
      6,
      7,
    ]
  }

  user_notices {
    end_time                 = 86399
    group_ids                = []
    need_phone_arrive_notice = 1
    notice_way               = [
      "EMAIL",
      "SMS",
    ]
    phone_call_type       = "CIRCLE"
    phone_circle_interval = 180
    phone_circle_times    = 2
    phone_inner_interval  = 180
    phone_order           = []
    receiver_type         = "USER"
    start_time            = 0
    user_ids              = [
      11082189,
      11082190,
    ]
    weekday = [
      1,
      2,
      3,
      4,
      5,
      6,
      7,
    ]
  }
}
```

Import

Monitor Alarm Notice can be imported, e.g.

```
$ terraform import tencentcloud_monitor_alarm_notice.import-test noticeId
```