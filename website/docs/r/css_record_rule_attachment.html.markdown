---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_record_rule_attachment"
sidebar_current: "docs-tencentcloud-resource-css_record_rule_attachment"
description: |-
  Provides a resource to create a css record_rule
---

# tencentcloud_css_record_rule_attachment

Provides a resource to create a css record_rule

## Example Usage

```hcl
resource "tencentcloud_css_record_rule_attachment" "record_rule" {
  domain_name = "177154.push.tlivecloud.com"
  template_id = 1262818
  app_name    = "qqq"
  stream_name = "ppp"
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, String, ForceNew) Streaming domain name.
* `template_id` - (Required, Int, ForceNew) Template ID.
* `app_name` - (Optional, String, ForceNew) The streaming path is consistent with the AppName in the streaming and playback addresses. The default is live.
* `stream_name` - (Optional, String, ForceNew) Stream name. Note: If this parameter is set to a non empty string, the rule will only work on this streaming.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

css record_rule can be imported using the id, e.g.

```
terraform import tencentcloud_css_record_rule_attachment.record_rule templateId#domainName
```

