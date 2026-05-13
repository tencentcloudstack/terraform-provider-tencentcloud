Provides a resource to confirm TEO origin ACL update for a zone. When the origin IP ranges of TEO change, you can use this resource to confirm that the latest origin IP ranges have been updated to the origin firewall, and the change notification will stop being pushed.

Example Usage

```hcl
resource "tencentcloud_teo_confirm_origin_acl_update_operation" "example" {
  zone_id = "zone-12345678"
}
```
