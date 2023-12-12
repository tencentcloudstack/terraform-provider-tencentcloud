Use this resource to create postgresql instance.

-> **Note:** To update the charge type, please update the `charge_type` and specify the `period` for the charging period. It only supports updating from `POSTPAID_BY_HOUR` to `PREPAID`, and the `period` field only valid in that upgrading case.

Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-1"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "guagua_vpc_instance_test"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "guagua_vpc_subnet_test"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create postgresql
resource "tencentcloud_postgresql_instance" "foo" {
  name              = "example"
  availability_zone = var.availability_zone
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  engine_version    = "10.4"
  root_user         = "root123"
  root_password     = "Root123$"
  charset           = "UTF8"
  project_id        = 0
  memory            = 2
  storage           = 10

  tags = {
    test = "tf"
  }
}
```

Create a multi available zone bucket

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-6"
}

variable "standby_availability_zone" {
  default = "ap-guangzhou-7"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "guagua_vpc_instance_test"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "guagua_vpc_subnet_test"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create postgresql
resource "tencentcloud_postgresql_instance" "foo" {
  name              = "example"
  availability_zone = var.availability_zone
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  engine_version    = "10.4"
  root_user         = "root123"
  root_password     = "Root123$"
  charset           = "UTF8"
  project_id        = 0
  memory            = 2
  storage           = 10

  db_node_set {
    role = "Primary"
    zone = var.availability_zone
  }
  db_node_set {
    zone = var.standby_availability_zone
  }

  tags = {
    test = "tf"
  }
}
```

create pgsql with kms key
```
resource "tencentcloud_postgresql_instance" "pg" {
  name              = "tf_postsql_instance"
  availability_zone = "ap-guangzhou-6"
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = "vpc-86v957zb"
  subnet_id         = "subnet-enm92y0m"
  engine_version    = "11.12"
  #  db_major_vesion   = "11"
  db_kernel_version = "v11.12_r1.3"
  need_support_tde  = 1
  kms_key_id        = "788c606a-c7b7-11ec-82d1-5254001e5c4e"
  kms_region        = "ap-guangzhou"
  root_password     = "xxxxxxxxxx"
  charset           = "LATIN1"
  project_id        = 0
  memory            = 4
  storage           = 100

  backup_plan {
    min_backup_start_time        = "00:10:11"
    max_backup_start_time        = "01:10:11"
    base_backup_retention_period = 7
    backup_period                = ["tuesday", "wednesday"]
  }

  tags = {
    tf = "test"
  }
}
```

upgrade kernel version
```
resource "tencentcloud_postgresql_instance" "test" {
  name = "tf_postsql_instance_update"
  availability_zone = data.tencentcloud_availability_zones_by_product.zone.zones[5].name
  charge_type	    = "POSTPAID_BY_HOUR"
  vpc_id  	  		= local.vpc_id
  subnet_id 		= local.subnet_id
  engine_version	= "13.3"
  root_password	    = "*"
  charset 			= "LATIN1"
  project_id 		= 0
  public_access_switch = false
  security_groups   = [local.sg_id]
  memory 			= 4
  storage 			= 250
  backup_plan {
	min_backup_start_time 		 = "01:10:11"
	max_backup_start_time		 = "02:10:11"
	base_backup_retention_period = 5
	backup_period 			     = ["monday", "thursday", "sunday"]
  }

  db_kernel_version = "v13.3_r1.4" # eg:from v13.3_r1.1 to v13.3_r1.4

  tags = {
	tf = "teest"
  }
}
```

Import

postgresql instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_postgresql_instance.foo postgres-cda1iex1
```