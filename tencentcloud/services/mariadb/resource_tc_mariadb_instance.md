Provides a resource to create a MariaDB instance(prepaid)

Example Usage

```hcl
resource "tencentcloud_mariadb_instance" "example" {
  zones           = ["ap-guangzhou-6", "ap-guangzhou-7"]
  instance_name   = "tf-example"
  node_count      = 2
  memory          = 8
  storage         = 500
  period          = 1
  vpc_id          = "vpc-i5yyodl9"
  subnet_id       = "subnet-hhi88a58"
  db_version_id   = "8.0"
  auto_renew_flag = 1
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

  tags = {
    createBy = "Terrafrom"
  }
}
```

Import

MariaDB instance(prepaid) can be imported using the id, e.g.

```
terraform import tencentcloud_mariadb_instance.example tdsql-4pzs5b67
```
