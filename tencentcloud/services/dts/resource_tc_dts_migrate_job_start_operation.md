Provides a resource to start a DTS migrate job

Example Usage

```hcl
resource "tencentcloud_mysql_instance" "example" {
  instance_name     = "tf-example"
  internet_service  = 1
  engine_version    = "5.7"
  charge_type       = "POSTPAID"
  root_password     = "PassWord@123"
  slave_deploy_mode = 0
  slave_sync_mode   = 1
  availability_zone = "ap-guangzhou-7"
  mem_size          = 128000
  volume_size       = 250
  vpc_id            = "vpc-i5yyodl9"
  subnet_id         = "subnet-hhi88a58"
  intranet_port     = 3306
  security_groups   = ["sg-e6a8xxib"]
  parameters = {
    character_set_server = "utf8"
    max_connections      = "1000"
  }
  tags = {
    createBy = "Terraform"
  }
}

resource "tencentcloud_cynosdb_cluster" "example" {
  cluster_name                 = "tf-example"
  db_mode                      = "NORMAL"
  db_type                      = "MYSQL"
  db_version                   = "5.7"
  port                         = 3306
  password                     = "Password@123"
  force_delete                 = true
  available_zone               = "ap-guangzhou-6"
  slave_zone                   = "ap-guangzhou-7"
  vpc_id                       = "vpc-i5yyodl9"
  subnet_id                    = "subnet-hhi88a58"
  instance_cpu_core            = 2
  instance_memory_size         = 4
  instance_maintain_duration   = 7200
  instance_maintain_start_time = 3600
  instance_maintain_weekdays = [
    "Fri",
    "Mon",
    "Sat",
    "Sun",
    "Thu",
    "Wed",
    "Tue",
  ]

  rw_group_sg = ["sg-e6a8xxib"]
  ro_group_sg = ["sg-e6a8xxib"]
}

resource "tencentcloud_dts_migrate_service" "example" {
  src_database_type = "mysql"
  dst_database_type = "cynosdbmysql"
  src_region        = "ap-guangzhou"
  dst_region        = "ap-guangzhou"
  instance_class    = "small"
  job_name          = "tf-example"
  tags {
    tag_key   = "createBy"
    tag_value = "Terraform"
  }
}

resource "tencentcloud_dts_migrate_job" "example" {
  service_id                    = tencentcloud_dts_migrate_service.example.id
  run_mode                      = "immediate"
  auto_retry_time_range_minutes = 0
  migrate_option {
    database_table {
      object_mode = "partial"
      databases {
        db_name    = "db_name"
        db_mode    = "partial"
        table_mode = "partial"
        tables {
          table_name      = "table_name"
          new_table_name  = "new_table_name"
          table_edit_mode = "rename"
        }
      }
    }
  }

  src_info {
    region        = "ap-guangzhou"
    access_type   = "cdb"
    database_type = "mysql"
    node_type     = "simple"
    info {
      user        = "root"
      password    = "Password@123"
      instance_id = tencentcloud_mysql_instance.example.id
    }
  }

  dst_info {
    region        = "ap-guangzhou"
    access_type   = "cdb"
    database_type = "cynosdbmysql"
    node_type     = "simple"
    info {
      user        = "root"
      password    = "Password@123"
      instance_id = tencentcloud_cynosdb_cluster.example.id
    }
  }
}

resource "tencentcloud_dts_migrate_job_start_operation" "example" {
  job_id = tencentcloud_dts_migrate_job.example.id
}
```