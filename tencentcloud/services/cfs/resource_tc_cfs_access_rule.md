Provides a resource to create a CFS access rule.

Example Usage

```hcl
resource "tencentcloud_cfs_access_rule" "foo" {
  access_group_id = "pgroup-7nx89k7l"
  auth_client_ip  = "10.10.1.0/24"
  priority        = 1
  rw_permission   = "RO"
  user_permission = "root_squash"
}
```