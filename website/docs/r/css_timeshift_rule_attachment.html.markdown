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

* Create: [CreateLiveTimeShiftRule](https://www.tencentcloud.com/document/product/267/53726) - This API is used to create a time shifting rule. You need to first call the [CreateLiveTranscodeTemplate](https://intl.cloud.tencent.com/document/product/267/86169?
* Read: [DescribeLiveTimeShiftRules](https://www.tencentcloud.com/document/product/267/53722) - This API is used to query time shifting rules.
* Delete: [DeleteLiveTimeShiftRule](https://www.tencentcloud.com/document/product/267/53724) - This API is used to delete a time shifting rule.
* Delete: [DescribeLiveTimeShiftRules](https://www.tencentcloud.com/document/product/267/53722) - This API is used to query time shifting rules.

## Import

css timeshift_rule_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_css_timeshift_rule_attachment.timeshift_rule_attachment templateId#domainName
```

