---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_callback_template"
sidebar_current: "docs-tencentcloud-resource-css_callback_template"
description: |-
  Provides a resource to create a css callback_template
---

# tencentcloud_css_callback_template

Provides a resource to create a css callback_template

## Example Usage

```hcl
resource "tencentcloud_css_callback_template" "callback_template" {
  template_name              = "tf-test"
  description                = "this is demo"
  stream_begin_notify_url    = "http://www.yourdomain.com/api/notify?action=streamBegin"
  stream_end_notify_url      = "http://www.yourdomain.com/api/notify?action=streamEnd"
  record_notify_url          = "http://www.yourdomain.com/api/notify?action=record"
  snapshot_notify_url        = "http://www.yourdomain.com/api/notify?action=snapshot"
  porn_censorship_notify_url = "http://www.yourdomain.com/api/notify?action=porn"
  callback_key               = "adasda131312"
  push_exception_notify_url  = "http://www.yourdomain.com/api/notify?action=pushException"
}
```

## Argument Reference

The following arguments are supported:

* `template_name` - (Required, String) Template name.Maximum length: 255 bytes. Only `Chinese`, `English`, `numbers`, `_`, `-` are supported.
* `callback_key` - (Optional, String) Callback Key, public callback URL.
* `description` - (Optional, String) Description information.Maximum length: 1024 bytes.Only `Chinese`, `English`, `numbers`, `_`, `-` are supported.
* `porn_censorship_notify_url` - (Optional, String) PornCensorship callback URL.
* `push_exception_notify_url` - (Optional, String) Streaming Exception Callback URL.
* `record_notify_url` - (Optional, String) Recording callback URL.
* `snapshot_notify_url` - (Optional, String) Snapshot callback URL.
* `stream_begin_notify_url` - (Optional, String) Launch callback URL.
* `stream_end_notify_url` - (Optional, String) Cutoff callback URL.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateLiveCallbackTemplate](https://cloud.tencent.com/document/api/267/32637) - This API is used to create a callback template. Up to 50 templates can be created. After the template ID is returned, you need to call the [CreateLiveCallbackRule](https://intl.cloud.tencent.com/document/product/267/32638?
* Read: [DescribeLiveCallbackTemplate](https://cloud.tencent.com/document/api/267/32633) - This API is used to get a single callback template.
* Update: [ModifyLiveCallbackTemplate](https://cloud.tencent.com/document/api/267/32631) - This API is used to modify a callback template.
* Delete: [DeleteLiveCallbackTemplate](https://cloud.tencent.com/document/api/267/32635) - This API deletes a callback template.

## Import

css callback_template can be imported using the id, e.g.

```
terraform import tencentcloud_css_callback_template.callback_template templateId
```

