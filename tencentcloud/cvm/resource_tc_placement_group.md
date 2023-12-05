Provide a resource to create a placement group.

Example Usage

```hcl
resource "tencentcloud_placement_group" "foo" {
  name = "test"
  type = "HOST"
}
```

Import

Placement group can be imported using the id, e.g.

```
$ terraform import tencentcloud_placement_group.foo ps-ilan8vjf
```