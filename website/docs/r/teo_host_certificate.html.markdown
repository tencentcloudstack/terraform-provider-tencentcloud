---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_host_certificate"
sidebar_current: "docs-tencentcloud-resource-teo_host_certificate"
description: |-
  Provides a resource to create a teo hostCertificate
---

# tencentcloud_teo_host_certificate

Provides a resource to create a teo hostCertificate

## Example Usage

```hcl
resource "tencentcloud_teo_host_certificate" "host_certificate" {
  zone_id = tencentcloud_teo_zone.zone.id
  host    = tencentcloud_teo_dns_record.dns_record.name

  cert_info {
    cert_id = "yqWPPbs7"
    status  = "deployed"
  }
}
```

## Argument Reference

The following arguments are supported:

* `host` - (Required, String) Domain name.
* `zone_id` - (Required, String) Site ID.
* `cert_info` - (Optional, List) Server certificate configuration.Note: This field may return null, indicating that no valid value can be obtained.

The `cert_info` object supports the following:

* `cert_id` - (Required, String) Server certificate ID, which is the ID of the default certificate. If you choose to upload an external certificate for SSL certificate management, a certificate ID will be generated.
* `status` - (Optional, String) Certificate deployment status.- processing: Deploying- deployed: DeployedNote: This field may return null, indicating that no valid value can be obtained.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

teo hostCertificate can be imported using the id, e.g.
```
$ terraform import tencentcloud_teo_host_certificate.host_certificate zoneId#host
```

