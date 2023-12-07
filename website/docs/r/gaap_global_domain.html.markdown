---
subcategory: "Global Application Acceleration(GAAP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_global_domain"
sidebar_current: "docs-tencentcloud-resource-gaap_global_domain"
description: |-
  Provides a resource to create a gaap global domain
---

# tencentcloud_gaap_global_domain

Provides a resource to create a gaap global domain

## Example Usage

```hcl
resource "tencentcloud_gaap_global_domain" "global_domain" {
  project_id    = 0
  default_value = "xxxxxx.com"
  alias         = "demo"
  tags = {
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `default_value` - (Required, String) Domain name default entry.
* `project_id` - (Required, Int) Domain Name Project ID.
* `alias` - (Optional, String) alias.
* `status` - (Optional, String) Global domain statue. Available values: open and close, default is open.
* `tags` - (Optional, Map) Instance tags.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

gaap global_domain can be imported using the id, e.g.

```
terraform import tencentcloud_gaap_global_domain.global_domain ${projectId}#${domainId}
```

