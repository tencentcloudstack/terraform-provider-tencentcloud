---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_check_free_certificate_verification"
sidebar_current: "docs-tencentcloud-resource-teo_check_free_certificate_verification"
description: |-
  Provides a resource to check free certificate verification for TEO (EdgeOne) domains.
---

# tencentcloud_teo_check_free_certificate_verification

Provides a resource to check free certificate verification for TEO (EdgeOne) domains.

## Example Usage

```hcl
resource "tencentcloud_teo_check_free_certificate_verification" "example" {
  zone_id = "zone-2o3h21ed8bsu"
  domain  = "example.com"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) The domain name to verify, which is the domain used when applying for a free certificate.
* `zone_id` - (Required, String, ForceNew) Zone ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `common_name` - The domain name to which the certificate is issued when the free certificate is successfully applied.
* `expire_time` - The expiration time of the certificate when the free certificate is successfully applied, in ISO 8601 format.
* `signature_algorithm` - The signature algorithm used by the certificate when the free certificate is successfully applied.


