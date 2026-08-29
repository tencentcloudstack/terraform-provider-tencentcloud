Provides a resource to create a TDSQL-C for MySQL (tdmysql) database instance.

Example Usage

```hcl
resource "tencentcloud_tdmysql_db_instance" "example" {
  zone               = "ap-guangzhou-6"
  vpc_id             = "vpc-i5yyodl9"
  subnet_id          = "subnet-hhi88a58"
  spec_code          = "4c8g"
  disk               = 200
  storage_node_num   = 3
  replications       = 3
  create_version     = "21.2.7.1"
  instance_name      = "tf-example"
  storage_node_cpu   = 4
  storage_node_mem   = 8
  pay_mode           = "0"
  vport              = 3306
  zones              = ["ap-guangzhou-6"]
  instance_type      = "hybrid"
  storage_type       = "CLOUD_HSSD"
  instance_mode      = "enhanced"
  sql_mode           = "MySQL"
  security_group_ids = ["sg-8gbd3tj9", "sg-2g6p85pr", "sg-hvnj11z7"]
  password           = "Password@2026"

  init_params {
    param = "character_set_server"
    value = "utf8mb4"
  }

  init_params {
    param = "lower_case_table_names"
    value = "0"
  }

  resource_tags {
    tag_key   = "CreatedBy"
    tag_value = "Terraform"
  }
}
```

Import

TDSQL-C for MySQL database instance can be imported using the id, e.g.

```
terraform import tencentcloud_tdmysql_db_instance.example tdsql3-f7e5dc9c
```
