---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_external_saml_identity_provider"
sidebar_current: "docs-tencentcloud-resource-organization_external_saml_identity_provider"
description: |-
  Provides a resource to create a Organization external saml identity provider
---

# tencentcloud_organization_external_saml_identity_provider

Provides a resource to create a Organization external saml identity provider

~> **NOTE:** When creating it for the first time, you must set one of `encoded_metadata_document` and `x509_certificate`, `another_x509_certificate` cannot be set alone.

## Example Usage

```hcl
resource "tencentcloud_organization_external_saml_identity_provider" "example" {
  zone_id                   = "z-1os7c9znogct"
  encoded_metadata_document = "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz......RGVzY3JpcHRvcj4="
  another_x509_certificate  = <<-EOF
-----BEGIN CERTIFICATE-----
MIIC8DCCAdigAwIBAgIQPCotiH/l8K1K6kBgL4mBfzANBgkqhkiG9w0BAQsFADA0MTIwMAYDVQQD
......
qs39KP9jOtSzEzc1YhiX
-----END CERTIFICATE-----
EOF
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String, ForceNew) Space ID.
* `another_x509_certificate` - (Optional, String) Another X509 certificate in PEM format. If this parameter is specified, all existing certificates will be replaced.
* `encoded_metadata_document` - (Optional, String) IdP metadata document (Base64 encoded). Provided by an IdP that supports the SAML 2.0 protocol.
* `entity_id` - (Optional, String, ForceNew) IdP identifier.
* `login_url` - (Optional, String, ForceNew) IdP login URL.
* `sso_status` - (Optional, String, ForceNew) SSO enabling status. Valid values: Enabled, Disabled (default).
* `x509_certificate` - (Optional, String) X509 certificate in PEM format. If this parameter is specified, all existing certificates will be replaced.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `another_certificate_id` - Another certificate ID.
* `certificate_id` - Certificate ID.
* `create_time` - Create time.
* `update_time` - Update time.


