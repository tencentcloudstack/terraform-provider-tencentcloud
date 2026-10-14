---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_callback_rule_attachment"
sidebar_current: "docs-tencentcloud-resource-css_callback_rule_attachment"
description: |-
  Provides a resource to create a css callback_rule
---

# tencentcloud_css_callback_rule_attachment

Provides a resource to create a css callback_rule

## Example Usage

```hcl
resource "tencentcloud_css_callback_rule_attachment" "callback_rule" {
  domain_name = "177154.push.tlivecloud.com"
  template_id = 434039
  app_name    = "live"
}
```

## Argument Reference

The following arguments are supported:

* `app_name` - (Required, String, ForceNew) The streaming path is consistent with the AppName in the streaming and playback addresses. The default is live.
* `domain_name` - (Required, String, ForceNew) Streaming domain name.
* `template_id` - (Required, Int, ForceNew) Template ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateLiveCallbackRule](https://www.tencentcloud.com/document/product/267/30816) - To create a callback rule, you need to first call the [CreateLiveCallbackTemplate](https://intl.cloud.tencent.com/document/product/267/32637?from_cn_redirect=1) API to create a callback template and bind the returned template ID to the dom
* Read: [DescribeLiveCallbackRules](https://www.tencentcloud.com/document/product/267/30812) - This API is used to get the callback rule list.
* Delete: [DeleteLiveCallbackRule](https://www.tencentcloud.com/document/product/267/30814) - This API is used to delete a callback rule.
* Delete: [DescribeLiveCallbackRules](https://www.tencentcloud.com/document/product/267/30812) - This API is used to get the callback rule list.

## Import

css callback_rule can be imported using the id, e.g.

```
terraform import tencentcloud_css_callback_rule_attachment.callback_rule templateId#domainName
```

