Use this resource to create PostgreSQL parameters.

Example Usage

```hcl
resource "tencentcloud_postgresql_instance" "example" {
  name              = "tf-example"
  availability_zone = "ap-guangzhou-6"
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = "vpc-i5yyodl9"         
  subnet_id         = "subnet-hhi88a58"
  db_major_version  = "17"
  engine_version    = "17.4"
  db_kernel_version = "v17.4_r1.4"
  root_user         = "root123"
  root_password     = "Root123$"
  charset           = "UTF8"
  project_id        = 0
  memory            = 4
  cpu               = 2
  storage           = 50
  tags = {
    CreateBy = "Terraform"
  }
}

resource "tencentcloud_postgresql_parameters" "example" {
  db_instance_id = tencentcloud_postgresql_instance.example.id
  param_list {
    name           = "check_function_bodies"
    expected_value = "off"
  }

  param_list {
    name           = "max_standby_archive_delay"
    expected_value = "35000"
  }
}
```

Import

PostgreSQL parameters can be imported using the id, e.g.

```
terraform import tencentcloud_postgresql_parameters.example postgres-ckwcgdf1
```
