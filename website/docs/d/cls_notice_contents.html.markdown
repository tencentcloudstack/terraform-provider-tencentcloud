---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_notice_contents"
sidebar_current: "docs-tencentcloud-datasource-cls_notice_contents"
description: |-
  Use this data source to query detailed information of CLS notice content templates.
---

# tencentcloud_cls_notice_contents

Use this data source to query detailed information of CLS notice content templates.

## Example Usage

### Query all notice content templates

```hcl
data "tencentcloud_cls_notice_contents" "example" {}
```

### Query by template name

```hcl
data "tencentcloud_cls_notice_contents" "example" {
  filters {
    key    = "name"
    values = ["DefaultTemplate(English)"]
  }
}
```

### Query by template ID

```hcl
data "tencentcloud_cls_notice_contents" "example" {
  filters {
    key    = "noticeContentId"
    values = ["Default-en"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter conditions. Supported filter names: name (notice content template name), noticeContentId (notice content template ID). Each request supports up to 10 filters, and each filter value list supports up to 100 values.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `key` - (Required, String) Filter field name. Valid values: name, noticeContentId.
* `values` - (Required, List) Filter field values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `notice_content_list` - Notice content template list.
  * `create_time` - Creation time (Unix timestamp in seconds).
  * `flag` - Template flag. 0: user-defined, 1: system built-in.
  * `name` - Notice content template name.
  * `notice_content_id` - Notice content template ID.
  * `notice_contents` - Notice content template details.
    * `recovery_content` - Alarm recovery notification content template.
      * `content` - Notification content template body.
      * `headers` - Request headers (only for custom callback channel).
      * `title` - Notification content template title.
    * `trigger_content` - Alarm trigger notification content template.
      * `content` - Notification content template body.
      * `headers` - Request headers (only for custom callback channel).
      * `title` - Notification content template title.
    * `type` - Channel type. Valid values: Email, Sms, WeChat, Phone, WeCom, DingTalk, Lark, Http.
  * `sub_uin` - Creator/modifier sub-account ID.
  * `type` - Language type. 0: Chinese, 1: English.
  * `uin` - Creator primary account ID.
  * `update_time` - Update time (Unix timestamp in seconds).


