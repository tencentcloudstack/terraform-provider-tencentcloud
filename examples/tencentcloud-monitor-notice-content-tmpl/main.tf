# Example: Create a monitor notice content template

resource "tencentcloud_monitor_notice_content_tmpl" "example" {
  tmpl_name     = "example-notice-template"
  monitor_type  = "MT_QCE"
  tmpl_language = "zh"

  # Template contents in JSON format
  # Supports multiple notification channels: QCloudYehe, WeWorkRobot, DingDingRobot, FeiShuRobot
  tmpl_contents = jsonencode({
    "QCloudYehe" : [
      {
        "MatchingStatus" : ["Trigger"],
        "Template" : {
          "Email" : {
            "ContentTmpl" : base64encode("告警内容：{{.AlarmContent}}\n告警时间：{{.AlarmTime}}"),
            "TitleTmpl" : base64encode("告警通知：{{.AlarmName}}")
          },
          "SMS" : {
            "ContentTmpl" : base64encode("告警：{{.AlarmName}}"),
            "TitleTmpl" : base64encode("")
          }
        }
      }
    ],
    "DingDingRobot" : [
      {
        "MatchingStatus" : ["Trigger"],
        "Template" : {
          "TitleTmpl" : base64encode("钉钉告警：{{.AlarmName}}"),
          "ContentTmpl" : base64encode("告警详情：{{.AlarmContent}}")
        }
      }
    ],
    "FeiShuRobot" : [
      {
        "MatchingStatus" : ["Trigger"],
        "Template" : {
          "TitleTmpl" : base64encode("飞书告警：{{.AlarmName}}"),
          "ContentTmpl" : base64encode("告警详情：{{.AlarmContent}}")
        }
      }
    ],
    "WeWorkRobot" : [
      {
        "MatchingStatus" : ["Trigger"],
        "Template" : {
          "ContentTmpl" : base64encode("企业微信告警：{{.AlarmName}}\n详情：{{.AlarmContent}}")
        }
      }
    ]
  })
}

# Output the template ID
output "tmpl_id" {
  value = tencentcloud_monitor_notice_content_tmpl.example.id
}
