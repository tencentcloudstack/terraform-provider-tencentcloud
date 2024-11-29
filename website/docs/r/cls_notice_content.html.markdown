---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_notice_content"
sidebar_current: "docs-tencentcloud-resource-cls_notice_content"
description: |-
  Provides a resource to create a cls notice content
---

# tencentcloud_cls_notice_content

Provides a resource to create a cls notice content

## Example Usage

```hcl
resource "tencentcloud_cls_notice_content" "example" {
  name = "tf-example"
  type = 0
  notice_contents {
    type = "Email"

    trigger_content {
      title   = "title"
      content = "This is content."
      headers = [
        "Content-Type:application/json"
      ]
    }

    recovery_content {
      title   = "title"
      content = "This is content."
      headers = [
        "Content-Type:application/json"
      ]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Notice content name.
* `notice_contents` - (Optional, List) Template detailed configuration.
* `type` - (Optional, Int) Template content language. 0: Chinese 1: English.

The `notice_contents` object supports the following:

* `type` - (Required, String) Channel type. Email: Email; Sms: SMS; WeChat: WeChat; Phone: Telephone; WeCom: Enterprise WeChat; DingTalk: DingTalk; Lark: Feishu; HTTP: Custom callback.
* `recovery_content` - (Optional, List) Template for Alarm Recovery Notification Content.
* `trigger_content` - (Optional, List) Alarm triggered notification content template.

The `recovery_content` object of `notice_contents` supports the following:

* `content` - (Optional, String) Notification content template body information.
* `headers` - (Optional, Set) Request headers: In HTTP requests, request headers contain additional information sent by the client to the server, such as user agent, authorization credentials, expected response format, etc. Only `custom callback` supports this configuration.
* `title` - (Optional, String) Notification content template title information. Some notification channel types do not support 'title', please refer to the Tencent Cloud Console page.

The `trigger_content` object of `notice_contents` supports the following:

* `content` - (Optional, String) Notification content template body information.
* `headers` - (Optional, Set) Request headers: In HTTP requests, request headers contain additional information sent by the client to the server, such as user agent, authorization credentials, expected response format, etc. Only `custom callback` supports this configuration.
* `title` - (Optional, String) Notification content template title information. Some notification channel types do not support 'title', please refer to the Tencent Cloud Console page.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cls notice content can be imported using the id, e.g.

```
terraform import tencentcloud_cls_notice_content.example noticetemplate-b417f32a-bdf9-46c5-933e-28c23cd7a6b7
```

