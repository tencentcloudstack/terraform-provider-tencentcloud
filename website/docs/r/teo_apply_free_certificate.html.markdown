---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_apply_free_certificate"
sidebar_current: "docs-tencentcloud-resource-teo_apply_free_certificate"
description: |-
  Provides a resource to apply TEO (EdgeOne) free certificate for a domain.
---

# tencentcloud_teo_apply_free_certificate

Provides a resource to apply TEO (EdgeOne) free certificate for a domain.

## Example Usage

### Apply free certificate with DNS verification

```hcl
resource "tencentcloud_teo_apply_free_certificate" "example" {
  zone_id             = "zone-2o3h21ed8bsf"
  domain              = "www.example.com"
  verification_method = "dns_challenge"
}
```

### Apply free certificate with HTTP file verification

```hcl
resource "tencentcloud_teo_apply_free_certificate" "example_http" {
  zone_id             = "zone-2o3h21ed8bsf"
  domain              = "www.example.com"
  verification_method = "http_challenge"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) The target domain for the free certificate application.
* `verification_method` - (Required, String, ForceNew) The verification method for the free certificate application. Valid values: `http_challenge`, `dns_challenge`.
* `zone_id` - (Required, String, ForceNew) Zone ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `dns_verification` - DNS verification information. Returned when `verification_method` is `dns_challenge`.
  * `record_type` - The record type.
  * `record_value` - The record value.
  * `subdomain` - The host record.
* `file_verification` - File verification information. Returned when `verification_method` is `http_challenge`.
  * `content` - The content of the verification file.
  * `path` - The URL path for file verification.


