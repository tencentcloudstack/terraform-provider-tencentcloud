---
subcategory: "Teo"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_default_certificate"
sidebar_current: "docs-tencentcloud-resource-teo_default_certificate"
description: |-
  Provides a resource to create a teo defaultCertificate
---

# tencentcloud_teo_default_certificate

Provides a resource to create a teo defaultCertificate

## Example Usage

```hcl
resource "tencentcloud_teo_default_certificate" "default_certificate" {
  zone_id = tencentcloud_teo_zone.zone.id

  cert_info {
    cert_id = "teo-28i46c1gtmkl"
    status  = "deployed"
  }
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String) Site ID.
* `cert_info` - (Optional, List) List of default certificates. Note: This field may return null, indicating that no valid value can be obtained.

The `cert_info` object supports the following:

* `cert_id` - (Required, String) Server certificate ID, which is the ID of the default certificate. If you choose to upload an external certificate for SSL certificate management, a certificate ID will be generated.
* `status` - (Optional, String) Certificate status.- applying: Application in progress.- failed: Application failed.- processing: Deploying certificate.- deployed: Certificate deployed.- disabled: Certificate disabled.Note: This field may return null, indicating that no valid value can be obtained.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

teo default_certificate can be imported using the id, e.g.
```
$ terraform import tencentcloud_teo_default_certificate.default_certificate zoneId
```

