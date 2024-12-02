Provides a resource to create a cls alarm

Example Usage

Use single condition

```hcl
resource "tencentcloud_cls_alarm" "example" {
  name = "tf-example"
  alarm_notice_ids = [
    "notice-c2af43ee-1a4b-4c4a-ae3e-f81481280101",
  ]
  alarm_period     = 15
  condition        = "$1.source='10.0.0.1'"
  alarm_level      = 1
  message_template = "{{.Label}}"
  status           = true
  trigger_count    = 1

  alarm_targets {
    logset_id         = "e74efb8e-f647-48b2-a725-43f11b122081"
    topic_id          = "59cf3ec0-1612-4157-be3f-341b2e7a53cb"
    query             = "status:>500 | select count(*) as errorCounts"
    start_time_offset = -15
    end_time_offset   = 0
    number            = 1
    syntax_rule       = 1
  }

  analysis {
    content = "__FILENAME__"
    name    = "terraform"
    type    = "field"

    config_info {
      key   = "QueryIndex"
      value = "1"
    }
  }

  monitor_time {
    time = 1
    type = "Period"
  }

  tags = {
    createdBy = "terraform"
  }
}
```

Use multi conditions

```hcl
resource "tencentcloud_cls_alarm" "example" {
  name = "tf-example"
  alarm_notice_ids = [
    "notice-c2af43ee-1a4b-4c4a-ae3e-f81481280101",
  ]
  alarm_period     = 15
  message_template = "{{.Label}}"
  status           = true
  trigger_count    = 1

  alarm_targets {
    logset_id         = "e74efb8e-f647-48b2-a725-43f11b122081"
    topic_id          = "59cf3ec0-1612-4157-be3f-341b2e7a53cb"
    query             = "status:>500 | select count(*) as errorCounts"
    start_time_offset = -15
    end_time_offset   = 0
    number            = 1
    syntax_rule       = 1
  }

  analysis {
    content = "__FILENAME__"
    name    = "terraform"
    type    = "field"

    config_info {
      key   = "QueryIndex"
      value = "1"
    }
  }

  multi_conditions {
    condition   = "[$1.__QUERYCOUNT__]> 0"
    alarm_level = 1
  }

  multi_conditions {
    condition   = "$1.source='10.0.0.1'"
    alarm_level = 2
  }

  monitor_time {
    time = 1
    type = "Period"
  }

  tags = {
    createdBy = "terraform"
  }
}
```

Import

cls alarm can be imported using the id, e.g.

```
terraform import tencentcloud_cls_alarm.example alarm-d8529662-e10f-440c-ba80-50f3dcf215a3
```