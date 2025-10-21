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



## Import

css callback_rule can be imported using the id, e.g.

```
terraform import tencentcloud_css_callback_rule_attachment.callback_rule templateId#domainName
```

