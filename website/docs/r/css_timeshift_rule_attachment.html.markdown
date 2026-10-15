---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_timeshift_rule_attachment"
sidebar_current: "docs-tencentcloud-resource-css_timeshift_rule_attachment"
description: |-
  Provides a resource to create a css timeshift_rule_attachment
---

# tencentcloud_css_timeshift_rule_attachment

Provides a resource to create a css timeshift_rule_attachment

## Example Usage

```hcl
resource "tencentcloud_css_timeshift_rule_attachment" "timeshift_rule_attachment" {
  domain_name = "177154.push.tlivecloud.com"
  template_id = 252586
  app_name    = "qqq"
  stream_name = "ppp"
}
```

## Argument Reference

The following arguments are supported:

* `app_name` - (Required, String, ForceNew) The push path, which should be the same as `AppName` in the push and playback URLs. The default value is `live`.
* `domain_name` - (Required, String, ForceNew) The push domain.
* `stream_name` - (Required, String, ForceNew) The stream name.Note: If you pass in a non-empty string, the rule will only be applied to the specified stream.
* `template_id` - (Required, Int, ForceNew) The template ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateLiveTimeShiftRule](https://cloud.tencent.com/document/api/267/86170) - 创建直播时移规则
* Read: [DescribeLiveTimeShiftRules](https://cloud.tencent.com/document/api/267/86166) - 获取直播时移规则列表
* Delete: [DeleteLiveTimeShiftRule](https://cloud.tencent.com/document/api/267/86168) - 删除直播时移规则
* Delete: [DescribeLiveTimeShiftRules](https://cloud.tencent.com/document/api/267/86166) - 获取直播时移规则列表

## Import

css timeshift_rule_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_css_timeshift_rule_attachment.timeshift_rule_attachment templateId#domainName
```

