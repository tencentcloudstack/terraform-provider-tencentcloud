# Example: Create a monitor notice content template

resource "tencentcloud_monitor_notice_content_tmpl" "example" {
  tmpl_name     = "example-notice-template"
  monitor_type  = "MT_QCE"
  tmpl_language = "zh"

  # Template contents configuration
  # Supports multiple notification channels: qcloud_yehe, we_work_robot, ding_ding_robot, fei_shu_robot
  tmpl_contents {
    # QCloud Yehe notification channel
    qcloud_yehe {
      matching_status = ["Trigger", "Recovery"]
      template {
        email {
          title_tmpl   = base64encode("告警通知：{{.AlarmName}}")
          content_tmpl = base64encode("告警内容：{{.AlarmContent}}\n告警时间：{{.AlarmTime}}")
        }
        sms {
          content_tmpl = base64encode("告警：{{.AlarmName}}")
        }
      }
    }

    # WeWork Robot notification channel
    we_work_robot {
      matching_status = ["Trigger"]
      template {
        content_tmpl = base64encode("企业微信告警：{{.AlarmName}}\n详情：{{.AlarmContent}}")
      }
    }

    # DingDing Robot notification channel
    ding_ding_robot {
      matching_status = ["Trigger"]
      template {
        title_tmpl   = base64encode("钉钉告警：{{.AlarmName}}")
        content_tmpl = base64encode("告警详情：{{.AlarmContent}}")
      }
    }

    # FeiShu Robot notification channel
    fei_shu_robot {
      matching_status = ["Trigger"]
      template {
        title_tmpl   = base64encode("飞书告警：{{.AlarmName}}")
        content_tmpl = base64encode("告警详情：{{.AlarmContent}}")
      }
    }
  }
}

# Output the template ID
output "tmpl_id" {
  value = tencentcloud_monitor_notice_content_tmpl.example.id
}
