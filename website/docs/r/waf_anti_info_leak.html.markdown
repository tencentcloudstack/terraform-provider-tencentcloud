---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_anti_info_leak"
sidebar_current: "docs-tencentcloud-resource-waf_anti_info_leak"
description: |-
  Provides a resource to create a waf anti_info_leak
---

# tencentcloud_waf_anti_info_leak

Provides a resource to create a waf anti_info_leak

## Example Usage

```hcl
resource "tencentcloud_waf_anti_info_leak" "example" {
  domain      = "tf.example.com"
  name        = "tf_example"
  action_type = 0
  strategies {
    field   = "information"
    content = "phone"
  }
  uri    = "/anti_info_leak_url"
  status = 1
}
```

## Argument Reference

The following arguments are supported:

* `action_type` - (Required, Int) Rule Action. 0: alarm; 1: replacement; 2: only displaying the first four digits; 3: only displaying the last four digits; 4: blocking.
* `domain` - (Required, String) Domain.
* `name` - (Required, String) Rule Name.
* `strategies` - (Required, List) Strategies detail.
* `uri` - (Required, String) Uri.
* `status` - (Optional, Int) status.

The `strategies` object supports the following:

* `content` - (Required, String) Matching Content. If field is returncode support: 400, 403, 404, 4xx, 500, 501, 502, 504, 5xx; If field is information support: idcard, phone, bankcard; If field is keywords users input matching content themselves.
* `field` - (Required, String) Matching Fields. support: returncode, keywords, information.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

waf anti_info_leak can be imported using the id, e.g.

```
terraform import tencentcloud_waf_anti_info_leak.example 3100077499#tf.example.com
```

