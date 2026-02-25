Provides a resource to create a CynosDB backup config

Example Usage

Enable logical backup configuration and cross-region logical backup

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
  password                     = "cynosDB@123"
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
    current_value = "1"
  }

  tags = {
    createBy = "terraform"
  }
}

resource "tencentcloud_cynosdb_backup_config" "example" {
  cluster_id       = tencentcloud_cynosdb_cluster.example.id
  backup_time_beg  = 7200
  backup_time_end  = 21600
  reserve_duration = 604800
  logic_backup_config {
    logic_backup_enable        = "ON"
    logic_backup_time_beg      = 7200
    logic_backup_time_end      = 21600
    logic_cross_regions        = ["ap-shanghai"]
    logic_cross_regions_enable = "ON"
    logic_reserve_duration     = 604800
  }
}
```

Disable logical backup configuration

```hcl
resource "tencentcloud_cynosdb_backup_config" "example" {
  cluster_id       = tencentcloud_cynosdb_cluster.example.id
  backup_time_beg  = 7200
  backup_time_end  = 21600
  reserve_duration = 604800
  logic_backup_config {
    logic_backup_enable = "OFF"
  }
}
```

Import

CynosDB backup config can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_backup_config.example cynosdbmysql-bws8h88b
```