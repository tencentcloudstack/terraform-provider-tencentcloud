Provides a resource to create a CLB target group.

Example Usage

If type is v1

```hcl
resource "tencentcloud_clb_target_group" "example" {
  target_group_name = "tf-example"
  vpc_id            = "vpc-jy6pwoy2"
  port              = 8090
  type              = "v1"

  tags {
    tag_key   = "tagKey"
    tag_value = "tagValue"
  }
}
```

If type is v2

```hcl
resource "tencentcloud_clb_target_group" "example" {
  target_group_name = "tf-example"
  vpc_id            = "vpc-jy6pwoy2"
  port              = 8090
  type              = "v2"
  protocol          = "TCP"
  weight            = 60

  tags {
    tag_key   = "tagKey"
    tag_value = "tagValue"
  }
}
```

Or full_listen_switch is true

```hcl
resource "tencentcloud_clb_target_group" "example" {
  target_group_name  = "tf-example"
  vpc_id             = "vpc-jy6pwoy2"
  type               = "v2"
  protocol           = "TCP"
  weight             = 60
  full_listen_switch = true

  tags {
    tag_key   = "tagKey"
    tag_value = "tagValue"
  }
}
```

Import

CLB target group can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_target_group.example lbtg-3k3io0i0
```
