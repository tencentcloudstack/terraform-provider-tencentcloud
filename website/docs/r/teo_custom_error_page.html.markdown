---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_custom_error_page"
sidebar_current: "docs-tencentcloud-resource-teo_custom_error_page"
description: |-
  Provides a resource to create a teo custom_error_page
---

# tencentcloud_teo_custom_error_page

Provides a resource to create a teo custom_error_page

## Example Usage

```hcl
resource "tencentcloud_teo_custom_error_page" "custom_error_page" {
  zone_id = ""
  entity  = ""
  name    = ""
  content = ""
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required, String) Page content.
* `entity` - (Required, String) Subdomain.
* `name` - (Required, String) Page name.
* `zone_id` - (Required, String) Site ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `page_id` - Page ID.


