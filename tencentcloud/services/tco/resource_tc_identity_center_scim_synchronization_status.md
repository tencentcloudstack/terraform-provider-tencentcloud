Provides a resource to manage identity center scim synchronization status

Example Usage

```hcl
resource "tencentcloud_identity_center_scim_synchronization_status" "identity_center_scim_synchronization_status" {
  zone_id = "z-xxxxxx"
  scim_synchronization_status = "Enabled"
}
```

Import

organization identity_center_scim_synchronization_status can be imported using the id, e.g.

```
terraform import tencentcloud_identity_center_scim_synchronization_status.identity_center_scim_synchronization_status ${zone_id}
```
