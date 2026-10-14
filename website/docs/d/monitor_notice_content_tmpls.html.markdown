---
subcategory: "Cloud Monitor(Monitor)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_notice_content_tmpls"
sidebar_current: "docs-tencentcloud-datasource-monitor_notice_content_tmpls"
description: |-
  Use this data source to query monitor notice content templates.
---

# tencentcloud_monitor_notice_content_tmpls

Use this data source to query monitor notice content templates.

## Example Usage

### Query all templates

```hcl
data "tencentcloud_monitor_notice_content_tmpls" "example" {}
```

### Query by filter

```hcl
data "tencentcloud_monitor_notice_content_tmpls" "example" {
  tmpl_ids      = ["ntpl-plu46bk5"]
  tmpl_name     = "tf-example"
  notice_id     = "notice-c2af43ee-1a4b-4c4a-ae3e-f81481280101"
  tmpl_language = "en"
  monitor_type  = "MT_QCE"
}
```

## Argument Reference

The following arguments are supported:

* `monitor_type` - (Optional, String) Monitor type for query. Valid value: `MT_QCE`.
* `notice_id` - (Optional, String) Notice template ID for query.
* `result_output_file` - (Optional, String) Used to save results.
* `tmpl_ids` - (Optional, Set: [`String`]) Template ID list for query.
* `tmpl_language` - (Optional, String) Template language for query. Valid values: `en`, `zh`.
* `tmpl_name` - (Optional, String) Template name for query.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `notice_content_tmpl_list` - Notification content template list.
  * `create_time` - Create time.
  * `creator` - Creator uin.
  * `monitor_type` - Monitor type.
  * `tmpl_contents_json` - Template content in JSON format.
  * `tmpl_id` - Template ID.
  * `tmpl_language` - Template language.
  * `tmpl_name` - Template name.
  * `update_time` - Update time.


