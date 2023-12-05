Provide a resource to create a CynosDB readonly instance.

Example Usage

```hcl
resource "tencentcloud_cynosdb_readonly_instance" "foo" {
  cluster_id           = cynosdbmysql-dzj5l8gz
  instance_name        = "tf-cynosdb-readonly-instance"
  force_delete         = true
  instance_cpu_core    = 2
  instance_memory_size = 4

  instance_maintain_duration   = 7200
  instance_maintain_start_time = 21600
  instance_maintain_weekdays   = [
    "Fri",
    "Mon",
    "Sat",
    "Sun",
    "Thu",
    "Wed",
    "Tue",
  ]
}
```

Import

CynosDB readonly instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_cynosdb_readonly_instance.foo cynosdbmysql-ins-dhwynib6
```