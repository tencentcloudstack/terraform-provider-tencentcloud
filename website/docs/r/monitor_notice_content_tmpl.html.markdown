---
subcategory: "Cloud Monitor(Monitor)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_notice_content_tmpl"
sidebar_current: "docs-tencentcloud-resource-monitor_notice_content_tmpl"
description: |-
  Use this resource to create Monitor notice content template.
---

# tencentcloud_monitor_notice_content_tmpl

Use this resource to create Monitor notice content template.

## Example Usage

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
          title_tmpl   = base64encode("AlarmTitle: {{.AlarmTitle}}")
        }
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `monitor_type` - (Required, String, ForceNew) Monitor type, e.g. MT_QCE.
* `tmpl_contents` - (Required, List) Template content configuration for different notification channels.
* `tmpl_language` - (Required, String) Template language, zh for Chinese, en for English.
* `tmpl_name` - (Required, String) Template name.

The `andon` object of `template` supports the following:

* `content_tmpl` - (Optional, String) Content template.
* `title_tmpl` - (Optional, String) Title template.

The `ding_ding_robot` object of `tmpl_contents` supports the following:

* `matching_status` - (Optional, List) Matching status list, e.g. Trigger, Recovery.
* `template` - (Optional, List) Template configuration.

The `email` object of `template` supports the following:

* `content_tmpl` - (Optional, String) Content template.
* `title_tmpl` - (Optional, String) Title template.

The `fei_shu_robot` object of `tmpl_contents` supports the following:

* `matching_status` - (Optional, List) Matching status list, e.g. Trigger, Recovery.
* `template` - (Optional, List) Template configuration.

The `headers` object of `template` supports the following:

* `key` - (Optional, String) Header key.
* `values` - (Optional, List) Header values.

The `pager_duty_robot` object of `tmpl_contents` supports the following:

* `matching_status` - (Optional, List) Matching status list.
* `template` - (Optional, List) PagerDuty template.

The `qcloud_yehe` object of `tmpl_contents` supports the following:

* `matching_status` - (Optional, List) Matching status list, e.g. Trigger, Recovery.
* `template` - (Optional, List) Template configuration.

The `qywx` object of `template` supports the following:

* `content_tmpl` - (Optional, String) Content template.
* `title_tmpl` - (Optional, String) Title template.

The `site` object of `template` supports the following:

* `content_tmpl` - (Optional, String) Content template.
* `title_tmpl` - (Optional, String) Title template.

The `sms` object of `template` supports the following:

* `content_tmpl` - (Optional, String) Content template.
* `title_tmpl` - (Optional, String) Title template.

The `teams_robot` object of `tmpl_contents` supports the following:

* `matching_status` - (Optional, List) Matching status list, e.g. Trigger, Recovery.
* `template` - (Optional, List) Template configuration.

The `template` object of `ding_ding_robot` supports the following:

* `content_tmpl` - (Optional, String) Content template.
* `title_tmpl` - (Optional, String) Title template.

The `template` object of `fei_shu_robot` supports the following:

* `content_tmpl` - (Optional, String) Content template.
* `title_tmpl` - (Optional, String) Title template.

The `template` object of `pager_duty_robot` supports the following:

* `body` - (Optional, String) Request body template in JSON.
* `headers` - (Optional, List) Request headers.
* `title_tmpl` - (Optional, String) Title template.

The `template` object of `qcloud_yehe` supports the following:

* `andon` - (Optional, List) Andon notification.
* `email` - (Optional, List) Email notification.
* `qywx` - (Optional, List) Enterprise WeChat notification.
* `site` - (Optional, List) Site notification.
* `sms` - (Optional, List) SMS notification.
* `voice` - (Optional, List) Voice notification.
* `wechat` - (Optional, List) WeChat notification.

The `template` object of `teams_robot` supports the following:

* `content_tmpl` - (Optional, String) Content template.

The `template` object of `we_work_robot` supports the following:

* `content_tmpl` - (Optional, String) Content template.

The `template` object of `webhook` supports the following:

* `body_content_type` - (Optional, String) Body content type.
* `body` - (Optional, String) Request body.
* `headers` - (Optional, List) Request headers.

The `tmpl_contents` object supports the following:

* `ding_ding_robot` - (Optional, List) DingDing Robot notification channel configuration.
* `fei_shu_robot` - (Optional, List) FeiShu Robot notification channel configuration.
* `pager_duty_robot` - (Optional, List) PagerDuty Robot notification channel configuration.
* `qcloud_yehe` - (Optional, List) QCloud Yehe notification channel configuration.
* `teams_robot` - (Optional, List) Teams Robot notification channel configuration.
* `we_work_robot` - (Optional, List) WeWork Robot notification channel configuration.
* `webhook` - (Optional, List) Webhook notification channel configuration.

The `voice` object of `template` supports the following:

* `content_tmpl` - (Optional, String) Content template.
* `title_tmpl` - (Optional, String) Title template.

The `we_work_robot` object of `tmpl_contents` supports the following:

* `matching_status` - (Optional, List) Matching status list, e.g. Trigger, Recovery.
* `template` - (Optional, List) Template configuration.

The `webhook` object of `tmpl_contents` supports the following:

* `matching_status` - (Optional, List) Matching status list.
* `template` - (Optional, List) Webhook template.

The `wechat` object of `template` supports the following:

* `alarm_content_tmpl` - (Optional, String) Alarm content template.
* `alarm_object_tmpl` - (Optional, String) Alarm object template.
* `alarm_region_tmpl` - (Optional, String) Alarm region template.
* `alarm_time_tmpl` - (Optional, String) Alarm time template.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `tmpl_id` - Template ID.


## Import

Monitor notice content template can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_notice_content_tmpl.example ntpl-3r1spzjn
```

