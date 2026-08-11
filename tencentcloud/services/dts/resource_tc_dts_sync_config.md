Provides a resource to create a DTS sync config

Example Usage

Sync MySQL database to CynosDB through cdb access type

```hcl
resource "tencentcloud_mysql_instance" "example" {
  instance_name     = "tf-example"
  internet_service  = 1
  engine_version    = "5.7"
  charge_type       = "POSTPAID"
  root_password     = "Mysql@2026"
  slave_deploy_mode = 0
  slave_sync_mode   = 0
  device_type       = "CLOUD_NATIVE_CLUSTER"
  availability_zone = "ap-guangzhou-6"
  cpu               = 2
  mem_size          = 4000
  volume_size       = 200
  vpc_id            = "vpc-i5yyodl9"
  subnet_id         = "subnet-hhi88a58"
  intranet_port     = 3306
  security_groups   = ["sg-4rd5741x"]
  parameters = {
    character_set_server = "utf8"
    max_connections      = "1000"
  }

  tags = {
    createBy = "terraform"
  }

  cluster_topology {
    read_write_node {
      zone = "ap-guangzhou-6"
    }
  }

  timeouts {
    create = "30m"
    delete = "30m"
  }
}

resource "tencentcloud_cynosdb_cluster" "example" {
  available_zone               = "ap-guangzhou-6"
  vpc_id                       = "vpc-i5yyodl9"
  subnet_id                    = "subnet-hhi88a58"
  db_mode                      = "NORMAL"
  db_type                      = "MYSQL"
  db_version                   = "5.7"
  port                         = 3306
  cluster_name                 = "tf-example"
  password                     = "CynosDB@2026"
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
    createBy = "terraform"
  }
}

resource "tencentcloud_dts_sync_job" "example" {
  pay_mode          = "PostPay"
  src_database_type = "mysql"
  src_region        = "ap-guangzhou"
  dst_database_type = "cynosdbmysql"
  dst_region        = "ap-guangzhou"
  auto_renew        = 0
  instance_class    = "micro"
  tags {
    tag_key   = "key"
    tag_value = "value"
  }
}

resource "tencentcloud_dts_sync_config" "example" {
  job_id          = tencentcloud_dts_sync_job.example.job_id
  src_access_type = "cdb"
  dst_access_type = "cdb"
  job_name        = "tf_example"
  job_mode        = "liteMode"
  run_mode        = "Immediate"

  objects {
    mode = "Partial"
    databases {
      db_name        = "testDB"
      db_mode        = "Partial"
      table_mode     = "Partial"
      view_mode      = "Partial"
      procedure_mode = "Partial"
      function_mode  = "Partial"
      tables {
        table_name  = "testTable"
        column_mode = "Partial"
        columns {
          column_name = "id"
        }

        tmp_tables = [
          "_testTable_new",
          "_testTable_old",
          "_testTable_ghc",
          "_testTable_gho",
          "_testTable_del",
        ]

        table_edit_mode = "pt"
      }
    }

    advanced_objects = [
      "procedure",
      "function",
    ]
  }

  src_info {
    region      = "ap-guangzhou"
    instance_id = tencentcloud_mysql_instance.example.id
    user        = "root"
    password    = tencentcloud_mysql_instance.example.root_password
    db_name     = "testDB"
    vpc_id      = "vpc-i5yyodl9"
    subnet_id   = "subnet-hhi88a58"
  }

  dst_info {
    region      = "ap-guangzhou"
    instance_id = tencentcloud_cynosdb_cluster.example.id
    user        = "root"
    password    = tencentcloud_cynosdb_cluster.example.password
    db_name     = "testDB"
    vpc_id      = "vpc-i5yyodl9"
    subnet_id   = "subnet-hhi88a58"
  }

  auto_retry_time_range_minutes = 5
}
```

Import

DTS sync config can be imported using the id, e.g.

```
terraform import tencentcloud_dts_sync_config.example sync-muu9ez38
```
