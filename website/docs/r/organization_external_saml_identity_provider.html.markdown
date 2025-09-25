---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_external_saml_identity_provider"
sidebar_current: "docs-tencentcloud-resource-organization_external_saml_identity_provider"
description: |-
  Provides a resource to create a organization organization_external_saml_identity_provider
---

# tencentcloud_organization_external_saml_identity_provider

Provides a resource to create a organization organization_external_saml_identity_provider

## Example Usage

```hcl
resource "tencentcloud_organization_external_saml_identity_provider" "organization_external_saml_identity_provider" {
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String, ForceNew) Space ID.
* `encoded_metadata_document` - (Optional, String, ForceNew) IdP metadata document (Base64 encoded). Provided by an IdP that supports the SAML 2.0 protocol.
* `entity_id` - (Optional, String, ForceNew) IdP identifier.
* `login_url` - (Optional, String, ForceNew) IdP login URL.
* `sso_status` - (Optional, String, ForceNew) SSO enabling status. Valid values: Enabled, Disabled (default).
* `x509_certificate` - (Optional, String, ForceNew) X509 certificate in PEM format. If this parameter is specified, all existing certificates will be replaced.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time.
* `update_time` - Update time.


## Import

organization organization_external_saml_identity_provider can be imported using the id, e.g.

```
terraform import tencentcloud_organization_external_saml_identity_provider.organization_external_saml_identity_provider organization_external_saml_identity_provider_id
```

