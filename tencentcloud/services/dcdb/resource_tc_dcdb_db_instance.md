Provides a resource to create a dcdb db_instance

Example Usage

```hcl
resource "tencentcloud_dcdb_db_instance" "db_instance" {
  instance_name = "test_dcdb_db_instance"
  zones = ["ap-guangzhou-5"]
  period = 1
  shard_memory = "2"
  shard_storage = "10"
  shard_node_count = "2"
  shard_count = "2"
  vpc_id = local.vpc_id
  subnet_id = local.subnet_id
  db_version_id = "8.0"
  resource_tags {
	tag_key = "aaa"
	tag_value = "bbb"
  }
  init_params {
	 param = "character_set_server"
	 value = "utf8mb4"
  }
  init_params {
	param = "lower_case_table_names"
	value = "1"
  }
  init_params {
	param = "sync_mode"
	value = "2"
  }
  init_params {
	param = "innodb_page_size"
	value = "16384"
  }
  security_group_ids = [local.sg_id]
}
```

Import

dcdb db_instance can be imported using the id, e.g.

```
terraform import tencentcloud_dcdb_db_instance.db_instance db_instance_id
```
