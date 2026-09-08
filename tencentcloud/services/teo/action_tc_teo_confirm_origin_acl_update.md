Provides an action to confirm TEO origin ACL update for a zone. When the origin IP ranges of TEO change, you can use this action to confirm that the latest origin IP ranges have been updated to the origin firewall, and the change notification will stop being pushed.

~> **NOTE:** Actions are supported in HashiCorp Terraform version 1.14 and later.

Example Usage

```hcl
action "tencentcloud_teo_confirm_origin_acl_update" "example" {
  config {
    zone_id = "zone-3fkff38fyw8s"
  }
}
```
