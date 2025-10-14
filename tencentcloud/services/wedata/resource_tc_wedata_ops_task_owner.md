Provides a resource to create a wedata ops task owner

Example Usage

```hcl
resource "tencentcloud_wedata_ops_task_owner" "wedata_ops_task_owner" {
    owner_uin  = "100029411056;100042282926"
    project_id = "2430455587205529600"
    task_id    = "20251009144419600"
}
```

Import

wedata ops task owner can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_ops_task_owner.wedata_ops_task_owner projectId#askId
```