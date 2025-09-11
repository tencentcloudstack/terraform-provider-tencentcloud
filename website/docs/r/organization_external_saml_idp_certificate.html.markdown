---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_external_saml_idp_certificate"
sidebar_current: "docs-tencentcloud-resource-organization_external_saml_idp_certificate"
description: |-
  Provides a resource to create a Organization external saml identity provider certificate
---

# tencentcloud_organization_external_saml_idp_certificate

Provides a resource to create a Organization external saml identity provider certificate

## Example Usage

```hcl
resource "tencentcloud_organization_external_saml_idp_certificate" "example" {
  zone_id          = "z-dsj3ieme"
  x509_certificate = "MIIBtjCCAVugAwIBAgITBmyf1XSXNmY/Owua2eiedgPySjAKBggqhkj********"
}
```

## Argument Reference

The following arguments are supported:

* `x509_certificate` - (Required, String, ForceNew) X509 certificate in PEM format, provided by the SAML identity provider.
* `zone_id` - (Required, String, ForceNew) Space ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `certificate_id` - Certificate ID.
* `issuer` - Certificate issuer.
* `not_after` - Certificate expiration date.
* `not_before` - Certificate creation date.
* `serial_number` - Certificate serial number.
* `signature_algorithm` - Certificate signature algorithm.
* `version` - Certificate version.


## Import

Organization external saml identity provider certificate can be imported using the id, e.g.

```
terraform import tencentcloud_organization_external_saml_idp_certificate.example z-dsj3ieme#idp-c-2jd8923je29dr34
```

