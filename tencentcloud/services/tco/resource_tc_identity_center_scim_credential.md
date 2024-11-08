Provides a resource to create an identity center scim credential

Example Usage

```hcl
resource "tencentcloud_identity_center_scim_credential" "identity_center_scim_credential" {
  zone_id = "z-xxxxxx"
}
```

Import

organization identity_center_scim_credential can be imported using the id, e.g.

```
terraform import tencentcloud_identity_center_scim_credential.identity_center_scim_credential ${zone_id}#${credential_id}
```
