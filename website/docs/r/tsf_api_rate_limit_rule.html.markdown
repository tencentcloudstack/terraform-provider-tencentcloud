---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_api_rate_limit_rule"
sidebar_current: "docs-tencentcloud-resource-tsf_api_rate_limit_rule"
description: |-
  Provides a resource to create a tsf api_rate_limit_rule
---

# tencentcloud_tsf_api_rate_limit_rule

Provides a resource to create a tsf api_rate_limit_rule

## Example Usage

```hcl
resource "tencentcloud_tsf_api_rate_limit_rule" "api_rate_limit_rule" {
  api_id        = "api-xxxxxx"
  max_qps       = 10
  usable_status = "enable"
}
```

## Argument Reference

The following arguments are supported:

* `api_id` - (Required, String) Api Id.
* `max_qps` - (Required, Int) qps value.
* `usable_status` - (Optional, String) Enabled/disabled, enabled/disabled, if not passed, it is enabled by default.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_time` - creation time.
* `description` - describe.
* `rule_content` - Rule content.
* `rule_id` - rule Id.
* `rule_name` - Current limit name.
* `tsf_rule_id` - Tsf Rule ID.
* `updated_time` - update time.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateApiRateLimitRuleWithDetailResp](https://cloud.tencent.com/document/api/649/89211)
* Read: [DescribeApiRateLimitRules](https://cloud.tencent.com/document/api/649/50635)
* Update: [UpdateApiRateLimitRule](https://cloud.tencent.com/document/api/649/50624)
* Delete: [DeleteApiRateLimitRule](https://cloud.tencent.com/document/api/649/89208)

## Import

tsf api_rate_limit_rule can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_api_rate_limit_rule.api_rate_limit_rule api_rate_limit_rule_id
```

