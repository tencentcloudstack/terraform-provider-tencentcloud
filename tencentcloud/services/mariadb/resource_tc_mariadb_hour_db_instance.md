Provides a resource to create a MariaDB hour_db_instance

Example Usage

```hcl
resource "tencentcloud_mariadb_hour_db_instance" "basic" {
  db_version_id = "10.0"
  instance_name = "db-test-del"
  memory        = 2
  node_count    = 2
  storage       = 10
  subnet_id     = "subnet-jdi5xn22"
  vpc_id        = "vpc-k1t8ickr"
  vip           = "10.0.0.197"
  zones         = [
    "ap-guangzhou-6",
    "ap-guangzhou-7",
  ]
  tags          = {
    createdBy   = "terraform"
  }
}
```

Create with custom init_params

```hcl
resource "tencentcloud_mariadb_hour_db_instance" "with_init_params" {
  db_version_id = "10.0"
  instance_name = "db-test-init"
  memory        = 2
  node_count    = 2
  storage       = 10
  subnet_id     = "subnet-jdi5xn22"
  vpc_id        = "vpc-k1t8ickr"
  zones         = [
    "ap-guangzhou-6",
    "ap-guangzhou-7",
  ]

  init_params {
    param = "character_set_server"
    value = "utf8mb4"
  }
  init_params {
    param = "lower_case_table_names"
    value = "1"
  }
  init_params {
    param = "innodb_page_size"
    value = "16384"
  }
  init_params {
    param = "sync_mode"
    value = "2"
  }
}
```

Import

mariadb hour_db_instance can be imported using the id, e.g.
```
$ terraform import tencentcloud_mariadb_hour_db_instance.hour_db_instance tdsql-kjqih9nn
```