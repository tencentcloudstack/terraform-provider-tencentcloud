---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_modify_access_period"
sidebar_current: "docs-tencentcloud-resource-waf_modify_access_period"
description: |-
  Provides a resource to create a waf modify_access_period
---

# tencentcloud_waf_modify_access_period

Provides a resource to create a waf modify_access_period

## Example Usage

```hcl
resource "tencentcloud_waf_modify_access_period" "example" {
  topic_id = "1ae37c76-df99-4e2b-998c-20f39eba6226"
  period   = 30
}
```

## Argument Reference

The following arguments are supported:

* `period` - (Required, Int, ForceNew) Access log retention period, range is [1, 180].
* `topic_id` - (Required, String, ForceNew) Log topic, new version does not need to be uploaded.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



