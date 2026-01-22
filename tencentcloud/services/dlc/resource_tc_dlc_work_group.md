Provides a resource to create a DLC work group

Example Usage

```hcl
resource "tencentcloud_dlc_work_group" "example" {
  work_group_name        = "tf-example"
  work_group_description = "DLC workgroup demo"
}
```

Import

DLC work group can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_work_group.example 135
```