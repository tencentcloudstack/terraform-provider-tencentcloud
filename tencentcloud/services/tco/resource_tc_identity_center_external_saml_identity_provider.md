Provides a resource to create a organization identity_center_external_saml_identity_provider

Example Usage

```hcl
resource "tencentcloud_identity_center_external_saml_identity_provider" "identity_center_external_saml_identity_provider" {
    zone_id = "z-xxxxxx"
    sso_status = "Enabled"
}
```

Import

organization identity_center_external_saml_identity_provider can be imported using the id, e.g.

```
terraform import tencentcloud_identity_center_external_saml_identity_provider.identity_center_external_saml_identity_provider ${zoneId}
```
