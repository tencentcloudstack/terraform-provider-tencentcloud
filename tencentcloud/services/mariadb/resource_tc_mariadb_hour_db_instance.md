Provides a resource to create a MariaDB hour db instance

Example Usage

Create with default init params

```hcl
resource "tencentcloud_mariadb_hour_db_instance" "example" {
  instance_name = "tf-example"
  memory        = 4
  node_count    = 2
  storage       = 100
  vpc_id        = "vpc-i5yyodl9"
  subnet_id     = "subnet-d4umunpy"
  vip           = "10.0.0.8"
  zones = [
    "ap-guangzhou-6",
    "ap-guangzhou-7",
  ]

  tags = {
    createdBy = "Terraform"
  }
}
```

Create with custom init params

```hcl
resource "tencentcloud_mariadb_hour_db_instance" "example" {
  db_version_id = "5.7"
  instance_name = "tf-example"
  memory        = 2
  node_count    = 2
  storage       = 100
  vpc_id        = "vpc-i5yyodl9"
  subnet_id     = "subnet-d4umunpy"
  vip           = "10.0.0.8"
  zones = [
    "ap-guangzhou-6",
    "ap-guangzhou-7",
  ]

  init_params {
    param = "character_set_server"
    value = "utf8"
  }

  init_params {
    param = "lower_case_table_names"
    value = "1"
  }

  tags = {
    createdBy = "Terraform"
  }
}
```

Import

MariaDB hour db instance can be imported using the id, e.g.
```
terraform import tencentcloud_mariadb_hour_db_instance.example tdsql-kjqih9nn
```