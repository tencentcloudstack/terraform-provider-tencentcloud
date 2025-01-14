---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_web_callback"
sidebar_current: "docs-tencentcloud-resource-cls_web_callback"
description: |-
  Provides a resource to create a CLS web callback
---

# tencentcloud_cls_web_callback

Provides a resource to create a CLS web callback

## Example Usage

### If type is WeCom

```hcl
resource "tencentcloud_cls_web_callback" "example" {
  name    = "tf-example"
  type    = "WeCom"
  webhook = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=7ac695f9-8af1-443f-b2c9-9a112f0647b5"
}
```

### If type is Http

```hcl
resource "tencentcloud_cls_web_callback" "example" {
  name    = "tf-example"
  type    = "Http"
  webhook = "https://demo.com"
  method  = "POST"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Notification content name.
* `type` - (Required, String) Channel type. WeCom: WeChat for Enterprise; DingTalk: DingTalk; Lark: Feishu; Http: custom callback.
* `webhook` - (Required, String) Webhook address.
* `key` - (Optional, String) Secret key.
* `method` - (Optional, String) Request method. Supports POST and PUT. Required when Type is Http.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

CLS web callback can be imported using the id, e.g.

```
terraform import tencentcloud_cls_web_callback.example webcallback-f2124b3d-e1e5-412c-9034-8e2fedeec952
```

