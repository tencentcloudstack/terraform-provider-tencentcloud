---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_cc_auto_status"
sidebar_current: "docs-tencentcloud-resource-waf_cc_auto_status"
description: |-
  Provides a resource to create a waf cc_auto_status
---

# tencentcloud_waf_cc_auto_status

Provides a resource to create a waf cc_auto_status

## Example Usage

```hcl
resource "tencentcloud_waf_cc_auto_status" "example" {
  domain  = "www.demo.com"
  edition = "sparta-waf"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) Domain.
* `edition` - (Required, String, ForceNew) Waf edition. clb-waf means clb-waf, sparta-waf means saas-waf.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `status` - cc auto status, 1 means open, 0 means close.


## Import

waf cc_auto_status can be imported using the id, e.g.

```
terraform import tencentcloud_waf_cc_auto_status.example www.demo.com#sparta-waf
```

