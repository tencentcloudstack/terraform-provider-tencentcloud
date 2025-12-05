---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_identity_center_external_saml_identity_provider"
sidebar_current: "docs-tencentcloud-resource-identity_center_external_saml_identity_provider"
description: |-
  Provides a resource to create a Organization identity center external saml identity provider
---

# tencentcloud_identity_center_external_saml_identity_provider

Provides a resource to create a Organization identity center external saml identity provider

## Example Usage

```hcl
resource "tencentcloud_identity_center_external_saml_identity_provider" "example" {
  zone_id                   = "z-1os7c9znogct"
  sso_status                = "Enabled"
  encoded_metadata_document = "PD94bWwgdmVyc2lvbj0iM......VzY3JpcHRvcj4="
}
```

### Or

```hcl
resource "tencentcloud_identity_center_external_saml_identity_provider" "example" {
  zone_id          = "z-1os7c9znogct"
  entity_id        = "https://sts.windows.net/d513d5bc-9f39-4069-ba9a-1eeab2ca58c1/"
  login_url        = "https://login.microsoftonline.com/d513d5bc-9f39-4069-ba9a-1eeab2ca58c1/saml2"
  sso_status       = "Enabled"
  x509_certificate = <<-EOF
-----BEGIN CERTIFICATE-----
MIIC8DCCAdigAwIBAgIQVbznAx6JSrhKG7gfJdx+jDANBgkqhkiG9w0BAQsFADA0MTIwMAYDVQQD
......
......
......
8hRskP2V6CH9PS0Zz2Zq
-----END CERTIFICATE-----
  EOF
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String, ForceNew) Space ID.
* `encoded_metadata_document` - (Optional, String) IdP metadata document (Base64 encoded). Provided by an IdP that supports the SAML 2.0 protocol.
* `entity_id` - (Optional, String) IdP identifier.
* `login_url` - (Optional, String) IdP login URL.
* `sso_status` - (Optional, String) SSO enabling status. Valid values: Enabled, Disabled (default).
* `x509_certificate` - (Optional, String) X509 certificate in PEM format. If this parameter is specified, all existing certificates will be replaced.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `acs_url` - Acs url.
* `certificate_ids` - Certificate ids.
* `create_time` - Create time.
* `update_time` - Update time.


## Import

Organization identity center external saml identity provider can be imported using the id, e.g.

```
terraform import tencentcloud_identity_center_external_saml_identity_provider.example z-1os7c9znogct
```

