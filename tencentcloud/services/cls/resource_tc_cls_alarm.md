Provides a resource to create a cls alarm

Example Usage

Use single condition with alarm_notice_ids

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

  classifications = {
    env     = "production"
    service = "api-gateway"
  }

  tags = {
    createdBy = "terraform"
  }
}
```

Use monitor_notice with flexible notification configuration

~> **NOTE:** `alarm_notice_ids` and `monitor_notice` are mutually exclusive. You can only use one of them.

```hcl
resource "tencentcloud_cls_alarm" "example_monitor_notice" {
  name = "tf-example-monitor-notice"

  # Use monitor_notice for more flexible notification configuration
  # Allows different notices for different alarm levels
  monitor_notice {
    notices {
      notice_id       = "notice-c2af43ee-1a4b-4c4a-ae3e-f81481280101"
      content_tmpl_id = "tmpl-5f7c8a9b-1234-5678-90ab-cdef12345678"
      alarm_levels    = [1, 2]  # Alert levels: 1=Critical, 2=Warning
    }

    notices {
      notice_id       = "notice-d3bf54ff-2b5c-5d5b-bf4f-f92582391202"
      content_tmpl_id = "tmpl-6g8d9b0c-2345-6789-01bc-def123456789"
      alarm_levels    = [3]  # Only send Info level alerts
    }
  }

  alarm_period     = 15
  condition        = "$1.errorCounts > 100"
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

  classifications = {
    env     = "production"
    service = "data-pipeline"
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

  classifications = {
    env     = "staging"
    service = "data-pipeline"
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