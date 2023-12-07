Provides a resource to create a dlc user

Example Usage

```hcl
resource "tencentcloud_dlc_user" "user" {
  user_id          = "100027012454"
  user_type        = "COMMON"
  user_alias       = "terraform-test"
  user_description = "for terraform test"
}
```

Import

dlc user can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_user.user user_id
```