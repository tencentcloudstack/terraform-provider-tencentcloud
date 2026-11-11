Provides a resource to create a TDSQL-C for MySQL (tdmysql) database instance.

Example Usage

```hcl
resource "tencentcloud_tdmysql_db_instance" "example" {
  zone               = "ap-guangzhou-6"
  vpc_id             = "vpc-i5yyodl9"
  subnet_id          = "subnet-hhi88a58"
  spec_code          = "8c16g"
  disk               = 400
  storage_node_num   = 3
  replications       = 3
  create_version     = "21.6.4.3"
  instance_name      = "tf-example"
  storage_node_cpu   = 8
  storage_node_mem   = 16
  pay_mode           = "0"
  vport              = 3306
  zones              = ["ap-guangzhou-6"]
  instance_type      = "hybrid"
  storage_type       = "CLOUD_HSSD"
  instance_mode      = "enhanced"
  sql_mode           = "MySQL"
  security_group_ids = ["sg-8gbd3tj9", "sg-2g6p85pr", "sg-hvnj11z7"]
  password           = "Password@2026"
  reserved_space_gb  = 5
  auto_renew_flag    = 0
  enable_ssl         = true

  init_params {
    param = "character_set_server"
    value = "utf8mb4"
  }

  init_params {
    param = "lower_case_table_names"
    value = "0"
  }

  init_params {
    param = "collation_server"
    value = "utf8mb4_0900_ai_ci"
  }

  init_params {
    param = "tdstore_enter_readonly_threshold_percent"
    value = "99.00"
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
terraform import tencentcloud_tdmysql_db_instance.example tdsql3-gpyfw59om6
```
