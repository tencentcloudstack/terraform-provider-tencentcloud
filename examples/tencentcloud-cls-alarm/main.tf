# Example: Create a CLS alarm with alarm_notice_ids

resource "tencentcloud_cls_alarm" "with_notice_ids" {
  name = "tf-example-alarm-notice-ids"
  
  # Use alarm_notice_ids (traditional way)
  alarm_notice_ids = [
    "notice-c2af43ee-1a4b-4c4a-ae3e-f81481280101",
  ]
  
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
    service = "api-gateway"
  }

  tags = {
    createdBy = "terraform"
  }
}

# Example: Create a CLS alarm with monitor_notice (new way)
# Note: alarm_notice_ids and monitor_notice are mutually exclusive

resource "tencentcloud_cls_alarm" "with_monitor_notice" {
  name = "tf-example-monitor-notice"

  # Use monitor_notice for more flexible notification configuration
  monitor_notice {
    notices {
      notice_id       = "notice-c2af43ee-1a4b-4c4a-ae3e-f81481280101"
      content_tmpl_id = "tmpl-5f7c8a9b-1234-5678-90ab-cdef12345678"
      alarm_levels    = [1, 2]  # Alert levels: 1=Critical, 2=Warning, 3=Info
    }

    notices {
      notice_id       = "notice-d3bf54ff-2b5c-5d5b-bf4f-f92582391202"
      content_tmpl_id = "tmpl-6g8d9b0c-2345-6789-01bc-def123456789"
      alarm_levels    = [3]  # Only send Info level alerts to this notice
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

# Example: Create a CLS alarm with multi_conditions and monitor_notice

resource "tencentcloud_cls_alarm" "multi_conditions_example" {
  name = "tf-example-multi-conditions"

  monitor_notice {
    notices {
      notice_id    = "notice-c2af43ee-1a4b-4c4a-ae3e-f81481280101"
      alarm_levels = [1, 2, 3]
    }
  }

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
    condition   = "[$1.__QUERYCOUNT__]> 100"
    alarm_level = 1
  }

  multi_conditions {
    condition   = "[$1.__QUERYCOUNT__]> 50"
    alarm_level = 2
  }

  monitor_time {
    time = 1
    type = "Period"
  }

  classifications = {
    env     = "staging"
    service = "web-service"
  }

  tags = {
    createdBy = "terraform"
  }
}
