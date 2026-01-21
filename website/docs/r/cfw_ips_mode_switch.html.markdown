---
subcategory: "Cloud Firewall(CFW)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cfw_ips_mode_switch"
sidebar_current: "docs-tencentcloud-resource-cfw_ips_mode_switch"
description: |-
  Provides a resource to create a CFW ips mode switch
---

# tencentcloud_cfw_ips_mode_switch

Provides a resource to create a CFW ips mode switch

## Example Usage

```hcl
resource "tencentcloud_cfw_ips_mode_switch" "example" {
  mode = 1
}
```

## Argument Reference

The following arguments are supported:

* `mode` - (Optional, Int) Protection mode: 0-observation mode, 1-interception mode, 2-strict mode.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

CFW ips mode switch can be imported using the customId(like uuid or base64 string), e.g.

```
terraform import tencentcloud_cfw_ips_mode_switch.example FTNxVFqU1BeA5JKfQlmkPg==
```

