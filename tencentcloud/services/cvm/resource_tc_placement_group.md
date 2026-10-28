Provide a resource to create a placement group.

Example Usage

```hcl
resource "tencentcloud_placement_group" "foo" {
  name     = "test"
  type     = "HOST"
  affinity = 2
  tags     = {
    createBy = "terraform"
  }
}
```

Create partition placement group

```hcl
resource "tencentcloud_placement_group" "bar" {
  name            = "test-partition"
  type            = "HOST"
  strategy        = "PARTITION"
  partition_count = 5
}
```

Import

Placement group can be imported using the id, e.g.

```
$ terraform import tencentcloud_placement_group.foo ps-ilan8vjf
```