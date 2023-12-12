Provide a resource to create a CynosDB cluster.

Example Usage

```hcl
resource "tencentcloud_cynosdb_cluster" "foo" {
  available_zone               = "ap-guangzhou-4"
  vpc_id                       = "vpc-h70b6b49"
  subnet_id                    = "subnet-q6fhy1mi"
  db_type                      = "MYSQL"
  db_version                   = "5.7"
  storage_limit                = 1000
  cluster_name                 = "tf-cynosdb"
  password                     = "cynos@123"
  instance_maintain_duration   = 7200
  instance_maintain_start_time = 10800
  instance_maintain_weekdays   = [
    "Fri",
    "Mon",
    "Sat",
    "Sun",
    "Thu",
    "Wed",
    "Tue",
  ]

  instance_cpu_core    = 1
  instance_memory_size = 2

  param_item {
    name = "character_set_server"
    current_value = "utf8mb4"
  }

  prarm_template_id = "12345"

  tags = {
    test = "test"
  }

  force_delete = false

  rw_group_sg = [
    "sg-ibyjkl6r",
  ]
  ro_group_sg = [
    "sg-ibyjkl6r",
  ]
}
```

Import

CynosDB cluster can be imported using the id, e.g.

```
$ terraform import tencentcloud_cynosdb_cluster.foo cynosdbmysql-dzj5l8gz
```