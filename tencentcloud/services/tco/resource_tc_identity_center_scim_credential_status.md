Provides a resource to manage identity center scim credential status

Example Usage

```hcl
resource "tencentcloud_identity_center_scim_credential_status" "identity_center_scim_credential_status" {
  zone_id = "z-xxxxxx"
  credential_id = "scimcred-xxxxxx"
  status = "Enabled"
}
```

Import

organization identity_center_scim_credential_status can be imported using the id, e.g.

```
terraform import tencentcloud_identity_center_scim_credential_status.identity_center_scim_credential_status ${zone_id}#${credential_id}
```
