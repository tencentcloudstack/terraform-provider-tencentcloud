---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_default_certificate"
sidebar_current: "docs-tencentcloud-resource-teo_default_certificate"
description: |-
  Provides a resource to create a teo default_certificate
---

# tencentcloud_teo_default_certificate

Provides a resource to create a teo default_certificate

## Example Usage

```hcl
resource "tencentcloud_teo_default_certificate" "default_certificate" {
  zone_id = ""
  cert_info {
    cert_id = ""
    status  = ""
  }
}
```

## Argument Reference

The following arguments are supported:

* `cert_info` - (Required, List) List of default certificates. Note: This field may return null, indicating that no valid value can be obtained.
* `zone_id` - (Required, String) Site ID.

The `cert_info` object supports the following:

* `cert_id` - (Required, String) Server certificate ID, which is the ID of the default certificate. If you choose to upload an external certificate for SSL certificate management, a certificate ID will be generated.
* `status` - (Optional, String) Certificate status.- `applying`: Application in progress.- `failed`: Application failed.- `processing`: Deploying certificate.- `deployed`: Certificate deployed.- `disabled`: Certificate disabled. Note: This field may return null, indicating that no valid value can be obtained.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

teo default_certificate can be imported using the id, e.g.
```
$ terraform import tencentcloud_teo_default_certificate.default_certificate defaultCertificate_id
```

