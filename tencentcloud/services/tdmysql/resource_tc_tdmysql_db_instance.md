Provides a resource to create a TDSQL-C MySQL(tdmysql) db instance.

Example Usage

```hcl
resource "tencentcloud_tdmysql_db_instance" "example" {
  zone              = "ap-guangzhou-3"
  vpc_id            = "vpc-xxxxxxxx"
  subnet_id         = "subnet-xxxxxxxx"
  spec_code         = "spec-code"
  disk              = 100
  storage_node_num  = 2
  replications      = 3
  instance_count    = 1
  instance_name     = "tf-tdmysql-example"
  pay_mode          = "0"
  instance_type     = "separate"
  storage_type      = "CLOUD_HSSD"
  time_unit         = "m"
  time_span         = 1
  init_params {
    param = "character_set_server"
    value = "utf8"
  }
  init_params {
    param = "lower_case_table_names"
    value = "0"
  }
  resource_tags {
    tag_key   = "createdBy"
    tag_value = "terraform"
  }
}
```

Import

TDSQL-C MySQL(tdmysql) db instance can be imported using the instance id, e.g.

```
terraform import tencentcloud_tdmysql_db_instance.example tdsqlshard-test1234
```
