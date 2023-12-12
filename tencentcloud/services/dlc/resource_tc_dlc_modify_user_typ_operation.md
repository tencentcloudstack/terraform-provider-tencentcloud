Provides a resource to create a dlc modify_user_typ_operation

Example Usage

```hcl
resource "tencentcloud_dlc_modify_user_typ_operation" "modify_user_typ_operation" {
  user_id = "127382378"
  user_type = "ADMIN"
}
```

Import

dlc modify_user_typ_operation can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_modify_user_typ_operation.modify_user_typ_operation modify_user_typ_operation_id
```