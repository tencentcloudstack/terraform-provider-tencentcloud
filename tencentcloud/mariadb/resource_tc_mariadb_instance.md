Provides a resource to create a mariadb instance

Example Usage

```hcl
resource "tencentcloud_mariadb_instance" "instance" {
  zones = ["ap-guangzhou-3",]
  node_count = 2
  memory = 8
  storage = 10
  period = 1
  # auto_voucher =
  # voucher_ids =
  vpc_id = "vpc-ii1jfbhl"
  subnet_id = "subnet-3ku415by"
  # project_id = ""
  db_version_id = "8.0"
  instance_name = "terraform-test"
  # security_group_ids = ""
  auto_renew_flag = 1
  ipv6_flag = 0
  tags = {
    "createby" = "terrafrom-2"
  }
  init_params {
    param = "character_set_server"
    value = "utf8mb4"
  }
  init_params {
    param = "lower_case_table_names"
    value = "0"
  }
  init_params {
    param = "innodb_page_size"
    value = "16384"
  }
  init_params {
    param = "sync_mode"
    value = "1"
  }
  dcn_region = ""
  dcn_instance_id = ""
}
```
Import

mariadb tencentcloud_mariadb_instance can be imported using the id, e.g.
```
$ terraform import tencentcloud_mariadb_instance.instance tdsql-4pzs5b67
```