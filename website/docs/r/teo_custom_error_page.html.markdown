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
resource "tencentcloud_teo_custom_error_page" "error_page_0" {
  zone_id = data.tencentcloud_teo_zone_ddos_policy.zone_policy.zone_id
  entity  = data.tencentcloud_teo_zone_ddos_policy.zone_policy.shield_areas[0].application[0].host

  name    = "test"
  content = "<html lang='en'><body><div><p>test content</p></div></body></html>"
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


