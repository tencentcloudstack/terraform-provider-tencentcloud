Provides a resource to create a dasb user_group_members

Example Usage

```hcl
resource "tencentcloud_dasb_user_group_members" "example" {
  user_group_id = 3
  member_id_set = [1, 2, 3]
}
```

Import

dasb user_group_members can be imported using the id, e.g.

```
terraform import tencentcloud_dasb_user_group_members.example 3#1,2,3
```