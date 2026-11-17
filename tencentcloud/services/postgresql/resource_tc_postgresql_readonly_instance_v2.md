Provides a resource to create a PostgreSQL readonly instance V2.

Example Usage

Create a postpaid PostgreSQL readonly instance

```hcl
resource "tencentcloud_postgresql_readonly_instance_v2" "example" {
  zone                  = "ap-guangzhou-6"
  master_db_instance_id = "postgres-ckwcgdf1"
  spec_code             = "pg.it.small2"
  storage               = 100
  vpc_id                = "vpc-lfnb8yyp"
  subnet_id             = "subnet-csay7wfo"
  instance_charge_type  = "POSTPAID_BY_HOUR"
  name                  = "tf-example"
  project_id            = 1269443
  read_only_group_id    = "pgrogrp-h62yf0re"
  deletion_protection   = false
  security_group_ids = [
    "sg-n8zf5ry9",
    "sg-rs32zv1r"
  ]

  tags = {
    createBy = "Terraform"
  }
}
```

Create a prepaid PostgreSQL readonly instance

```hcl
resource "tencentcloud_postgresql_readonly_instance_v2" "example" {
  zone                  = "ap-guangzhou-6"
  master_db_instance_id = "postgres-ckwcgdf1"
  spec_code             = "pg.it.small2"
  storage               = 100
  vpc_id                = "vpc-lfnb8yyp"
  subnet_id             = "subnet-csay7wfo"
  instance_charge_type  = "PREPAID"
  period                = 1
  auto_renew_flag       = 1
  name                  = "tf-example"
  project_id            = 1269443
  read_only_group_id    = "pgrogrp-h62yf0re"
  deletion_protection   = true
  security_group_ids = [
    "sg-n8zf5ry9",
    "sg-8gbd3tj9",
    "sg-rs32zv1r"
  ]

  tags = {
    createBy = "Terraform"
  }
}
```

Import

PostgreSQL readonly instance V2 can be imported using the id, e.g.

```
terraform import tencentcloud_postgresql_readonly_instance_v2.example pgro-gr4od5iw
```
