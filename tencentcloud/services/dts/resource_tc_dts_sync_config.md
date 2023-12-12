Provides a resource to create a dts sync_config

Example Usage

Sync mysql database to cynosdb through cdb access type

```hcl
resource "tencentcloud_cynosdb_cluster" "foo" {
	available_zone               = var.availability_zone
	vpc_id                       = local.vpc_id
	subnet_id                    = local.subnet_id
	db_type                      = "MYSQL"
	db_version                   = "5.7"
	storage_limit                = 1000
	cluster_name                 = "tf-cynosdb-mysql-sync-dst"
	password                     = "*"
	instance_maintain_duration   = 3600
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
	param_items {
	  name = "character_set_server"
	  current_value = "utf8"
	}
	param_items {
	  name = "time_zone"
	  current_value = "+09:00"
	}
	param_items {
		name = "lower_case_table_names"
		current_value = "1"
	}

	force_delete = true

	rw_group_sg = [
	  local.sg_id
	]
	ro_group_sg = [
	  local.sg_id
	]
	prarm_template_id = var.my_param_template
  }

resource "tencentcloud_dts_sync_job" "sync_job" {
	pay_mode = "PostPay"
	src_database_type = "mysql"
	src_region = "ap-guangzhou"
	dst_database_type = "cynosdbmysql"
	dst_region = "ap-guangzhou"
	tags {
	  tag_key = "aaa"
	  tag_value = "bbb"
	}
	auto_renew = 0
	instance_class = "micro"
  }

resource "tencentcloud_dts_sync_config" "sync_config" {
  job_id = tencentcloud_dts_sync_job.sync_job.job_id
  src_access_type = "cdb"
  dst_access_type = "cdb"

  job_name = "tf_test_sync_config"
  job_mode = "liteMode"
  run_mode = "Immediate"

  objects {
	mode = "Partial"
      databases {
	    db_name = "tf_ci_test"
			new_db_name = "tf_ci_test_new"
			db_mode = "Partial"
			table_mode = "All"
			tables {
				table_name = "test"
				new_table_name = "test_new"
			}
	  }
  }
  src_info {
		region        = "ap-guangzhou"
		instance_id   = "cdb-fitq5t9h"
		user          = "your_user_name"
		password      = "*"
		db_name       = "tf_ci_test"
		vpc_id        = local.vpc_id
		subnet_id     = local.subnet_id
  }
  dst_info {
		region        = "ap-guangzhou"
		instance_id   = tencentcloud_cynosdb_cluster.foo.id
		user          = "root"
		password      = "*"
		db_name       = "tf_ci_test_new"
		vpc_id        = local.vpc_id
		subnet_id     = local.subnet_id
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

resource "tencentcloud_dts_sync_config" "sync_config" {
  job_id          = data.tencentcloud_dts_sync_jobs.sync_jobs.list.0.job_id
  src_access_type = "ccn"
  dst_access_type = "cdb"

  job_mode = "liteMode"
  run_mode = "Immediate"

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

dts sync_config can be imported using the id, e.g.

```
terraform import tencentcloud_dts_sync_config.sync_config sync_config_id
```