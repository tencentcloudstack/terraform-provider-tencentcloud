---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_domain_referer"
sidebar_current: "docs-tencentcloud-resource-css_domain_referer"
description: |-
  Provides a resource to create a css domain_referer
---

# tencentcloud_css_domain_referer

Provides a resource to create a css domain_referer

## Example Usage

```hcl
resource "tencentcloud_css_domain_referer" "domain_referer" {
  allow_empty = 1
  domain_name = "test122.jingxhu.top"
  enable      = 0
  rules       = "example.com"
  type        = 1
}
```

## Argument Reference

The following arguments are supported:

* `allow_empty` - (Required, Int) Allow blank referers, 0: not allowed, 1: allowed.
* `domain_name` - (Required, String) Domain Name.
* `enable` - (Required, Int) Whether to enable the referer blacklist authentication of the current domain name,`0`: off, `1`: on.
* `rules` - (Required, String) The list of referers to; separate.
* `type` - (Required, Int) List type: 0: blacklist, 1: whitelist.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

css domain_referer can be imported using the id, e.g.

```
terraform import tencentcloud_css_domain_referer.domain_referer domainName
```

