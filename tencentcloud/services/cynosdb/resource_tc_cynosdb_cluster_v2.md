Provide a resource to create a CynosDB cluster.

~> **NOTE:** Compared to Resource `tencentcloud_cynosdb_cluster`, Resource `tencentcloud_cynosdb_cluster_v2` places greater emphasis on optimizing security group configurations for read-only groups and read-only instances, making them more precise and efficient. `rw_group_sg` represents the read-write instance security group, `ro_group_sg` represents the read-only group security group, and `single_ro_group_sg` represents the read-only instance security group. notably, to configure `ro_group_sg`, ``open_ro_group` must be set `true` first. If you need to configure `ro_group_sg` or `single_ro_group_sg` security group, please use Resource `tencentcloud_cynosdb_cluster_v2`.

~> **NOTE:** params `instance_count` and `instance_init_infos` only choose one. If neither parameter is set, the CynosDB cluster is created with parameter `instance_count` set to `2` by default(one RW instance + one Ro instance). If you only need to create a master instance, explicitly set the `instance_count` field to `1`, or configure the RW instance information in the `instance_init_infos` field.

Example Usage

```hcl
resource "tencentcloud_cynosdb_cluster_v2" "example" {
  available_zone               = "ap-guangzhou-6"
  vpc_id                       = "vpc-i5yyodl9"
  subnet_id                    = "subnet-hhi88a58"
  db_mode                      = "NORMAL"
  db_type                      = "MYSQL"
  db_version                   = "5.7"
  port                         = 3306
  cluster_name                 = "tf-example"
  password                     = "cynosDB@123"
  instance_maintain_duration   = 7200
  instance_maintain_start_time = 10800
  instance_cpu_core            = 2
  instance_memory_size         = 4
  open_ro_group                = true
  force_delete                 = true
  instance_maintain_weekdays = [
    "Fri",
    "Mon",
    "Sat",
    "Sun",
    "Thu",
    "Wed",
    "Tue",
  ]

  rw_group_sg        = ["sg-7pnojfur", "sg-37tigqat"]
  ro_group_sg        = ["sg-7pnojfur", "sg-37tigqat", "sg-08cqf7d5"]
  single_ro_group_sg = ["sg-7pnojfur", "sg-37tigqat", "sg-08cqf7d5", "sg-l1txcqtj"]

  param_items {
    name          = "character_set_server"
    current_value = "utf8mb4"
  }

  param_items {
    name          = "lower_case_table_names"
    current_value = "0"
  }

  tags = {
    createBy = "terraform"
  }
}
```

Import

CynosDB cluster can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_cluster_v2.example cynosdbmysql-im25yazt
```
