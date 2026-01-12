Provides a resource to create a CLB target group.

Example Usage

Create V1 target group

```hcl
resource "tencentcloud_clb_target_group" "test" {
  target_group_name = "test"
  port              = 33
  type              = "v1"
}
```

Create V2 target group

```hcl
resource "tencentcloud_clb_target_group" "test" {
  target_group_name = "test"
  port              = 33
  type              = "v2"
  protocol          = "TCP"
}
```

Import

CLB target group can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_target_group.test lbtg-3k3io0i0
```