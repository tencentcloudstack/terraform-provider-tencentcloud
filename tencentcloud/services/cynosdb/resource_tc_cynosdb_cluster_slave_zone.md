Provides a resource to create a CynosDB cluster slave zone.

~> **NOTE:** If you use resource `tencentcloud_cynosdb_cluster_slave_zone` to configure `slave_zone` for `tencentcloud_cynosdb_cluster`, then you cannot simultaneously set the `slave_zone` field of resource `tencentcloud_cynosdb_cluster`.

Example Usage

Set a new slave zone for a cynosdb cluster

```hcl
resource "tencentcloud_cynosdb_cluster" "example" {
  available_zone               = "ap-guangzhou-6"
  vpc_id                       = "vpc-i5yyodl9"
  subnet_id                    = "subnet-hhi88a58"
  db_mode                      = "NORMAL"
  db_type                      = "MYSQL"
  db_version                   = "5.7"
  port                         = 3306
  cluster_name                 = "tf-example"
  password                     = "CynosDB@123"
  instance_maintain_duration   = 7200
  instance_maintain_start_time = 10800
  instance_cpu_core            = 2
  instance_memory_size         = 4
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

  param_items {
    name          = "character_set_server"
    current_value = "utf8mb4"
  }

  param_items {
    name          = "lower_case_table_names"
    current_value = "0"
  }

  tags = {
    createBy = "Terraform"
  }
}

resource "tencentcloud_cynosdb_cluster_slave_zone" "example" {
  cluster_id = tencentcloud_cynosdb_cluster.example.id
  slave_zone = "ap-guangzhou-7"
}
```

Import

CynosDB cluster slave zone can be imported using the clusterId#slaveZone, e.g.

```
terraform import tencentcloud_cynosdb_cluster_slave_zone.example cynosdbmysql-g76di9j5#ap-guangzhou-7
```
