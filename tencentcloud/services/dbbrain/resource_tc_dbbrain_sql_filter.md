Provides a resource to create a dbbrain sql_filter.

Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-3"
}
variable "region" {
  default = "ap-guangzhou"
}

data "tencentcloud_mysql_instance" "mysql" {
  instance_name = "instance_name"
}

locals {
  mysql_id = data.tencentcloud_mysql_instance.mysql.instance_list.0.mysql_id
}

resource "tencentcloud_dbbrain_sql_filter" "sql_filter" {
  instance_id = local.mysql_id
  session_token {
    user = "test"
	password = "===password==="
  }
  sql_type = "SELECT"
  filter_key = "filter_key"
  max_concurrency = 10
  duration = 3600
}

```