---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_domain"
sidebar_current: "docs-tencentcloud-resource-css_domain"
description: |-
  Provides a resource to create a css domain
---

# tencentcloud_css_domain

Provides a resource to create a css domain

## Example Usage

```hcl
resource "tencentcloud_css_domain" "domain" {
  domain_name          = "iac-tf.cloud"
  domain_type          = 0
  play_type            = 1
  is_delay_live        = 0
  is_mini_program_live = 0
  verify_owner_type    = "dbCheck"
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, String) Domain Name.
* `domain_type` - (Required, Int) Domain type: `0`: push stream. `1`: playback.
* `enable` - (Optional, Bool) Switch. true: enable the specified domain, false: disable the specified domain.
* `is_delay_live` - (Optional, Int) Whether it is LCB: `0`: LVB. `1`: LCB. Default value is 0.
* `is_mini_program_live` - (Optional, Int) `0`: LVB. `1`: LVB on Mini Program. Note: this field may return null, indicating that no valid values can be obtained. Default value is 0.
* `play_type` - (Optional, Int) Play Type. This parameter is valid only if `DomainType` is 1. Available values: `1`: in Mainland China. `2`: global. `3`: outside Mainland China. Default value is 1.
* `verify_owner_type` - (Optional, String) Domain name attribution verification type. `dnsCheck`, `fileCheck`, `dbCheck`. The default is `dbCheck`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

css domain can be imported using the id, e.g.

```
terraform import tencentcloud_css_domain.domain domain_name
```

