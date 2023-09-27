---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_watermark_rule_attachment"
sidebar_current: "docs-tencentcloud-resource-css_watermark_rule_attachment"
description: |-
  Provides a resource to create a css watermark_rule
---

# tencentcloud_css_watermark_rule_attachment

Provides a resource to create a css watermark_rule

## Example Usage

### Binding watermark rule with a css stream

```hcl
resource "tencentcloud_css_pull_stream_task" "example" {
  stream_name = "tf_example_stream_name"
  source_type = "PullLivePushLive"
  source_urls = ["rtmp://xxx.com/live/stream"]
  domain_name = "test.domain.com"
  app_name    = "live"
  start_time  = "2023-09-27T10:28:21Z"
  end_time    = "2023-09-27T17:28:21Z"
  operator    = "tf_admin"
  comment     = "This is a e2e test case."
}

resource "tencentcloud_css_watermark" "example" {
  picture_url    = "picture_url"
  watermark_name = "watermark_name"
  x_position     = 0
  y_position     = 0
  width          = 0
  height         = 0
}

resource "tencentcloud_css_watermark_rule_attachment" "watermark_rule_attachment" {
  domain_name = tencentcloud_css_pull_stream_task.example.domain_name
  app_name    = tencentcloud_css_pull_stream_task.example.app_name
  stream_name = tencentcloud_css_pull_stream_task.example.stream_name
  template_id = tencentcloud_css_watermark.example.id
}
```

## Argument Reference

The following arguments are supported:

* `app_name` - (Required, String, ForceNew) rule app name.
* `domain_name` - (Required, String, ForceNew) rule domain name.
* `stream_name` - (Required, String, ForceNew) rule stream name.
* `template_id` - (Required, Int, ForceNew) The template Id can be acquired by the Id of `tencentcloud_css_watermark`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - create time.
* `update_time` - update time.


## Import

css watermark_rule_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_css_watermark_rule_attachment.watermark_rule domain_name#app_name#stream_name#template_id
```

