Provides a resource to create a cls machine group.

Example Usage

```hcl
resource "tencentcloud_cls_machine_group" "group" {
  group_name        = "group"
  service_logging   = true
  tags              = {
    "test" = "test1"
  }
  update_end_time   = "19:05:40"
  update_start_time = "17:05:40"

  machine_group_type {
    type   = "ip"
    values = [
      "192.168.1.1",
      "192.168.1.2",
    ]
  }
}
```

Import

cls machine group can be imported using the id, e.g.

```
$ terraform import tencentcloud_cls_machine_group.group caf168e7-32cd-4ac6-bf89-1950a760e09c
```