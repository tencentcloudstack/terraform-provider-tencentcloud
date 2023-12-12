Provides a resource to create a ssm product_secret

Example Usage

Ssm secret for mysql

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "cdb"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_security_group" "security_group" {
  name        = "sg-example"
  description = "desc."
}

resource "tencentcloud_mysql_instance" "example" {
  internet_service  = 1
  engine_version    = "5.7"
  charge_type       = "POSTPAID"
  root_password     = "PassWord123"
  slave_deploy_mode = 0
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  slave_sync_mode   = 1
  instance_name     = "tf-example"
  mem_size          = 4000
  volume_size       = 200
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  intranet_port     = 3306
  security_groups   = [tencentcloud_security_group.security_group.id]

  tags = {
    createBy = "terraform"
  }

  parameters = {
    character_set_server = "utf8"
    max_connections      = "1000"
  }
}

resource "tencentcloud_kms_key" "example" {
  alias                = "tf-example-kms-key"
  description          = "example of kms key"
  key_rotation_enabled = false
  is_enabled           = true

  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_ssm_product_secret" "example" {
  secret_name      = "tf-example"
  user_name_prefix = "prefix"
  product_name     = "Mysql"
  instance_id      = tencentcloud_mysql_instance.example.id
  domains          = ["10.0.0.0"]
  privileges_list {
    privilege_name = "GlobalPrivileges"
    privileges     = ["ALTER ROUTINE"]
  }
  description         = "for ssm product test"
  kms_key_id          = tencentcloud_kms_key.example.id
  status              = "Enabled"
  enable_rotation     = true
  rotation_begin_time = "2023-08-05 20:54:33"
  rotation_frequency  = 30

  tags = {
    "createdBy" = "terraform"
  }
}
```

Ssm secret for tdsql-c-mysql
```hcl
resource "tencentcloud_ssm_product_secret" "example" {
  secret_name      = "tf-tdsql-c-example"
  user_name_prefix = "prefix"
  product_name     = "Tdsql_C_Mysql"
  instance_id      = "cynosdbmysql-xxxxxx"
  domains          = ["%"]
  privileges_list {
    privilege_name = "GlobalPrivileges"
    privileges     = [
      "ALTER",
      "CREATE",
      "DELETE",
    ]
  }
  privileges_list {
    privilege_name = "DatabasePrivileges"
    database       = "test"
    privileges     = [
      "ALTER",
      "CREATE",
      "DELETE",
      "SELECT",
    ]
  }
  description         = "test tdsql-c"
  kms_key_id          = null
  status              = "Enabled"
  enable_rotation     = false
  rotation_begin_time = "2023-08-05 20:54:33"
  rotation_frequency  = 30

  tags = {
    "createdBy" = "terraform"
  }
}
```