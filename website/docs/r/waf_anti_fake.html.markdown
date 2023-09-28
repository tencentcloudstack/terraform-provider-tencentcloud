---
subcategory: "Waf"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_anti_fake"
sidebar_current: "docs-tencentcloud-resource-waf_anti_fake"
description: |-
  Provides a resource to create a waf anti_fake
---

# tencentcloud_waf_anti_fake

Provides a resource to create a waf anti_fake

~> **NOTE:** Uri: Please configure static resources such as. html,. shtml,. txt,. js,. css,. jpg,. png, or access paths for static resources..

## Example Usage

```hcl
resource "tencentcloud_waf_anti_fake" "example" {
  domain = "www.waf.com"
  name   = "tf_example"
  uri    = "/anti_fake_url.html"
  status = 1
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String) Domain.
* `name` - (Required, String) Name.
* `uri` - (Required, String) Uri.
* `status` - (Optional, Int) status. 0: Turn off rules and log switches, 1: Turn on the rule switch and Turn off the log switch; 2: Turn off the rule switch and turn on the log switch;3: Turn on the log switch.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `protocol` - protocol.
* `rule_id` - rule id.


## Import

waf anti_fake can be imported using the id, e.g.

```
terraform import tencentcloud_waf_anti_fake.example 3200035516#www.waf.com
```

