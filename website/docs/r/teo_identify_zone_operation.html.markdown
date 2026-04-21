---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_identify_zone_operation"
sidebar_current: "docs-tencentcloud-resource-teo_identify_zone_operation"
description: |-
  Provides a resource to identify TEO zone or subdomain ownership.
---

# tencentcloud_teo_identify_zone_operation

Provides a resource to identify TEO zone or subdomain ownership.

## Example Usage

```hcl
resource "tencentcloud_teo_identify_zone_operation" "example" {
  zone_name = "example.com"
}
```

### With subdomain

```hcl
resource "tencentcloud_teo_identify_zone_operation" "example_sub" {
  zone_name = "example.com"
  domain    = "www.example.com"
}
```

## Argument Reference

The following arguments are supported:

* `zone_name` - (Required, String, ForceNew) Zone name.
* `domain` - (Optional, String, ForceNew) Subdomain under the zone. Required only when verifying a subdomain.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `ascription` - DNS verification information.
  * `record_type` - DNS record type.
  * `record_value` - DNS record value.
  * `subdomain` - DNS record host.
* `file_ascription` - File verification information.
  * `identify_content` - File verification content.
  * `identify_path` - File verification path.


