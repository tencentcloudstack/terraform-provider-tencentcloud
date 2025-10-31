Provides a resource to create a DTS sync config

Example Usage

Sync mysql database to cynosdb through cdb access type

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
      db_name     = "tf_ci_test"
      new_db_name = "tf_ci_test_new"
      db_mode     = "Partial"
      table_mode  = "All"
      tables {
        table_name     = "test"
        new_table_name = "test_new"
      }
    }
  }

  src_info {
    region      = "ap-guangzhou"
    instance_id = "cdb-fitq5t9h"
    user        = "your_user_name"
    password    = "*"
    db_name     = "tf_ci_test"
    vpc_id      = "vpc-i5yyodl9"
    subnet_id   = "subnet-hhi88a58"
  }

  dst_info {
    region      = "ap-guangzhou"
    instance_id = tencentcloud_cynosdb_cluster.example.id
    user        = "root"
    password    = "*"
    db_name     = "tf_ci_test_new"
    vpc_id      = "vpc-i5yyodl9"
    subnet_id   = "subnet-hhi88a58"
  }

  auto_retry_time_range_minutes = 0
}
```

Sync mysql database using CCN to route from ap-shanghai to ap-guangzhou

```hcl
locals {
  vpc_id_sh    = "vpc-evtcyb3g"
  subnet_id_sh = "subnet-1t83cxkp"
  src_ip       = data.tencentcloud_mysql_instance.src_mysql.instance_list.0.intranet_ip
  src_port     = data.tencentcloud_mysql_instance.src_mysql.instance_list.0.intranet_port
  ccn_id       = data.tencentcloud_ccn_instances.ccns.instance_list.0.ccn_id
  dst_mysql_id = data.tencentcloud_mysql_instance.dst_mysql.instance_list.0.mysql_id
}

variable "src_az_sh" {
  default = "ap-shanghai"
}

variable "dst_az_gz" {
  default = "ap-guangzhou"
}

data "tencentcloud_dts_sync_jobs" "sync_jobs" {
  job_name = "keep_sync_config_ccn_2_cdb"
}

data "tencentcloud_ccn_instances" "ccns" {
  name = "keep-ccn-dts-sh"
}

data "tencentcloud_mysql_instance" "src_mysql" {
  instance_name = "your_user_name_mysql_src"
}

data "tencentcloud_mysql_instance" "dst_mysql" {
  instance_name = "your_user_name_mysql_src"
}

resource "tencentcloud_dts_sync_config" "example" {
  job_id          = data.tencentcloud_dts_sync_jobs.sync_jobs.list.0.job_id
  src_access_type = "ccn"
  dst_access_type = "cdb"
  job_mode        = "liteMode"
  run_mode        = "Immediate"

  objects {
    mode = "Partial"
    databases {
      db_name     = "tf_ci_test"
      new_db_name = "tf_ci_test_new"
      db_mode     = "Partial"
      table_mode  = "All"
      tables {
        table_name     = "test"
        new_table_name = "test_new"
      }
    }
  }

  src_info { // shanghai to guangzhou via ccn
    region           = var.src_az_sh
    user             = "your_user_name"
    password         = "your_pass_word"
    ip               = local.src_ip
    port             = local.src_port
    vpc_id           = local.vpc_id_sh
    subnet_id        = local.subnet_id_sh
    ccn_id           = local.ccn_id
    database_net_env = "TencentVPC"
  }

  dst_info {
    region      = var.dst_az_gz
    instance_id = local.dst_mysql_id
    user        = "your_user_name"
    password    = "your_pass_word"
  }

  auto_retry_time_range_minutes = 0
}
````

Import

DTS sync config can be imported using the id, e.g.

```
terraform import tencentcloud_dts_sync_config.example sync-muu9ez38
```
