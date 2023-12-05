Provides a resource to create a dts migrate_job

Example Usage

```hcl
resource "tencentcloud_cynosdb_cluster" "foo" {
	available_zone               = var.availability_zone
	vpc_id                       = local.vpc_id
	subnet_id                    = local.subnet_id
	db_type                      = "MYSQL"
	db_version                   = "5.7"
	storage_limit                = 1000
	cluster_name                 = "tf-cynosdb-mysql"
	password                     = "cynos@123"
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

resource "tencentcloud_dts_migrate_service" "service" {
	src_database_type = "mysql"
	dst_database_type = "cynosdbmysql"
	src_region = "ap-guangzhou"
	dst_region = "ap-guangzhou"
	instance_class = "small"
	job_name = "tf_test_migration_service_1"
	tags {
	  tag_key = "aaa"
	  tag_value = "bbb"
	}
  }

resource "tencentcloud_dts_migrate_job" "job" {
  	service_id = tencentcloud_dts_migrate_service.service.id
	run_mode = "immediate"
	migrate_option {
		database_table {
			object_mode = "partial"
			databases {
				db_name = "tf_ci_test"
				db_mode = "partial"
				table_mode = "partial"
				tables {
					table_name = "test"
					new_table_name = "test_%s"
					table_edit_mode = "rename"
				}
			}
		}
	}
	src_info {
			region = "ap-guangzhou"
			access_type = "cdb"
			database_type = "mysql"
			node_type = "simple"
			info {
				user = "user_name"
				password = "your_pw"
				instance_id = "cdb-fitq5t9h"
			}

	}
	dst_info {
			region = "ap-guangzhou"
			access_type = "cdb"
			database_type = "cynosdbmysql"
			node_type = "simple"
			info {
				user = "user_name"
				password = "your_pw"
				instance_id = tencentcloud_cynosdb_cluster.foo.id
			}
	}
	auto_retry_time_range_minutes = 0
}

resource "tencentcloud_dts_migrate_job_start_operation" "start"{
	job_id = tencentcloud_dts_migrate_job.job.id
}
```

Import

dts migrate_job can be imported using the id, e.g.

```
terraform import tencentcloud_dts_migrate_job.migrate_job migrate_config_id
```