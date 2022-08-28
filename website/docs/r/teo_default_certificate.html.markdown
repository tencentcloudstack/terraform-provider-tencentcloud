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
resource "tencentcloud_teo_default_certificate" "defaultCertificate" {
  zone_id = ""
  cert_id = ""
}
```

## Argument Reference

The following arguments are supported:

* `cert_id` - (Required, String) Server certificate ID, which is the ID of the default certificate. If you choose to upload an external certificate for SSL certificate management, a certificate ID will be generated.
* `zone_id` - (Required, String) Site ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cert_info` - List of default certificates. Note: This field may return null, indicating that no valid value can be obtained.
  * `alias` - Certificate alias. Note: This field may return null, indicating that no valid value can be obtained.
  * `common_name` - Certificate common name. Note: This field may return null, indicating that no valid value can be obtained.
  * `effective_time` - Time when the certificate takes effect. Note: This field may return null, indicating that no valid value can be obtained.
  * `expire_time` - Time when the certificate expires. Note: This field may return null, indicating that no valid value can be obtained.
  * `message` - Returns a message to display failure causes when `Status` is failed.Note: This field may return null, indicating that no valid value can be obtained.
  * `status` - Certificate status.- applying: Application in progress.- failed: Application failed.- processing: Deploying certificate.- deployed: Certificate deployed.- disabled: Certificate disabled.Note: This field may return null, indicating that no valid value can be obtained.
  * `subject_alt_name` - Domain names added to the SAN certificate. Note: This field may return null, indicating that no valid value can be obtained.
  * `type` - Certificate type.- default: Default certificate.- upload: External certificate.- managed: Tencent Cloud managed certificate.Note: This field may return null, indicating that no valid value can be obtained.


## Import

teo defaultCertificate can be imported using the id, e.g.
```
$ terraform import tencentcloud_teo_default_certificate.defaultCertificate defaultCertificate_id
```

