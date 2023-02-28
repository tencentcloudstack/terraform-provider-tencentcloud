---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_domain_config"
sidebar_current: "docs-tencentcloud-resource-css_domain_config"
description: |-
  Provides a resource to configure(enable/disable) the css domain.
---

# tencentcloud_css_domain_config

Provides a resource to configure(enable/disable) the css domain.

## Example Usage

```hcl
resource "tencentcloud_css_domain_config" "enable_domain" {
  domain_name   = "your_domain_name"
  enable_domain = true
}
```
```hcl
resource "tencentcloud_css_domain_config" "forbid_domain" {
  domain_name   = "your_domain_name"
  enable_domain = false
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, String) Domain Name.
* `enable_domain` - (Required, Bool) Switch. true: enable the specified domain, false: disable the specified domain.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



