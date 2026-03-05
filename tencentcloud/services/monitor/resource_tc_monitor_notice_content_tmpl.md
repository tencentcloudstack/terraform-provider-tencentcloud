Use this resource to create Monitor notice content template.

Example Usage

```hcl
resource "tencentcloud_monitor_notice_content_tmpl" "example" {
  tmpl_name     = "tf-example"
  monitor_type  = "MT_QCE"
  tmpl_language = "en"
  tmpl_contents {
    qcloud_yehe {
      matching_status = ["Trigger"]
      template {
        email {
          content_tmpl = base64encode("AlarmContent: {{.AlarmContent}}")
          title_tmpl   =  base64encode("AlarmTitle: {{.AlarmTitle}}")
        }
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
