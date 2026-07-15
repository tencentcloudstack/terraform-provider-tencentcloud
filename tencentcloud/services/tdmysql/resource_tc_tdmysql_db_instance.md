Provides a resource to create a TDSQL-C for MySQL (tdmysql) database instance.

Example Usage

```hcl
resource "tencentcloud_tdmysql_db_instance" "example" {
  zone             = "ap-guangzhou-3"
  vpc_id           = "vpc-xxxxxxxx"
  subnet_id        = "subnet-xxxxxxxx"
  spec_code        = "TDSQL-C-LS001"
  disk             = 200
  storage_node_num = 2
  replications     = 3
  instance_name    = "tf-example-tdmysql"
  instance_count   = 1
  pay_mode         = "0"
  storage_type     = "CLOUD_HSSD"
  instance_type    = "separate"
  init_params {
    param = "character_set_server"
    value = "utf8"
  }
  init_params {
    param = "lower_case_table_names"
    value = "1"
  }
  resource_tags {
    tag_key   = "CreatedBy"
    tag_value = "Terraform"
  }
}
```

Import

TDSQL-C for MySQL database instance can be imported using the instance_id, e.g.

```
terraform import tencentcloud_tdmysql_db_instance.example tdmysqldb-xxxxxxxx
```
