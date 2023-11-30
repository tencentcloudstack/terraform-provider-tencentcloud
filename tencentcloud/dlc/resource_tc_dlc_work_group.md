Provides a resource to create a dlc work_group

Example Usage

```hcl
resource "tencentcloud_dlc_work_group" "work_group" {
  work_group_name        = "tf-demo"
  work_group_description = "dlc workgroup test"
}
```

Import

dlc work_group can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_work_group.work_group work_group_id
```