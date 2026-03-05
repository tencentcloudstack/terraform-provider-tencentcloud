Use this resource to create Monitor notice content template.

Example Usage

```hcl
resource "tencentcloud_monitor_notice_content_tmpl" "example" {
  tmpl_name     = "tf-example"
  monitor_type  = "MT_QCE"
  tmpl_language = "en"
  tmpl_contents {
    qcloud_yehe {
      matching_status = ["Trigger", "Recovery"]
      template {
        email {
          title_tmpl   = base64encode("AlarmTitle{{.AlarmName}}")
          content_tmpl = base64encode("AlarmContent{{.AlarmContent}}")
        }

        sms {
          content_tmpl = base64encode("Alarm: {{.AlarmName}}")
        }
      }
    }

    we_work_robot {
      matching_status = ["Trigger"]
      template {
        content_tmpl = base64encode("AlarmContent: {{.AlarmName}}")
      }
    }

    ding_ding_robot {
      matching_status = ["Trigger"]
      template {
        title_tmpl   = base64encode("AlarmTitle: {{.AlarmName}}")
        content_tmpl = base64encode("AlarmContent: {{.AlarmContent}}")
      }
    }

    fei_shu_robot {
      matching_status = ["Trigger"]
      template {
        title_tmpl   = base64encode("AlarmTitle: {{.AlarmName}}")
        content_tmpl = base64encode("AlarmContent: {{.AlarmContent}}")
      }
    }
  }
}
```

Import

Monitor notice content template can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_notice_content_tmpl.example ntpl-3r1spzjn
```
