Use this resource to create monitor notice content template.

Example Usage

```hcl
resource "tencentcloud_monitor_notice_content_tmpl" "example" {
  tmpl_name     = "tf-example-template"
  monitor_type  = "MT_QCE"
  tmpl_language = "zh"
  tmpl_contents = jsonencode({
    "MatchingStatus" : ["Trigger"],
    "Template" : {
      "WeWorkRobot" : {
        "TitleTmpl" : "告警通知",
        "ContentTmpl" : "告警详情：{{.Content}}"
      },
      "DingDingRobot" : {
        "TitleTmpl" : "告警通知",
        "ContentTmpl" : "告警详情：{{.Content}}"
      }
    }
  })
}
```

Import

Monitor notice content template can be imported using the id (format: `tmplID#tmplName`), e.g.

```
$ terraform import tencentcloud_monitor_notice_content_tmpl.example ntpl-3r1spzjn#tf-example-template
```
