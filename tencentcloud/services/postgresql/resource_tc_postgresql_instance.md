Use this resource to create postgresql instance.

-> **Note:** To update the charge type, please update the `charge_type` and specify the `period` for the charging period. It only supports updating from `POSTPAID_BY_HOUR` to `PREPAID`, and the `period` field only valid in that upgrading case.

-> **Note:** If no values are set for the two parameters: `db_major_version` and `engine_version`, then `engine_version` is set to `10.4` by default. Suggest using parameter `db_major_version` to create an instance

Example Usage

Create a postgresql instance

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create postgresql
resource "tencentcloud_postgresql_instance" "example" {
  name              = "example"
  availability_zone = var.availability_zone
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  db_major_version  = "10"
  engine_version    = "10.23"
  root_user         = "root123"
  root_password     = "Root123$"
  charset           = "UTF8"
  project_id        = 0
  cpu               = 1
  memory            = 2
  storage           = 10

  tags = {
    test = "tf"
  }
}
```

Create a postgresql instance with delete protection

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create postgresql
resource "tencentcloud_postgresql_instance" "example" {
  name              = "example"
  availability_zone = var.availability_zone
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  db_major_version  = "10"
  engine_version    = "10.23"
  root_user         = "root123"
  root_password     = "Root123$"
  charset           = "UTF8"
  project_id        = 0
  cpu               = 1
  memory            = 2
  storage           = 10
  delete_protection = true

  tags = {
    test = "tf"
  }
}
```

Create a multi available zone postgresql instance

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-6"
}

variable "standby_availability_zone" {
  default = "ap-guangzhou-7"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create postgresql
resource "tencentcloud_postgresql_instance" "example" {
  name              = "example"
  availability_zone = var.availability_zone
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  db_major_version  = "10"
  root_user         = "root123"
  root_password     = "Root123$"
  charset           = "UTF8"
  project_id        = 0
  memory            = 2
  cpu               = 1
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

Create a multi available zone postgresql instance of CDC

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-4"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create postgresql
resource "tencentcloud_postgresql_instance" "example" {
  name              = "tf-example"
  availability_zone = var.availability_zone
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  db_major_version  = "10"
  root_user         = "root123"
  root_password     = "Root123$"
  charset           = "UTF8"
  project_id        = 0
  memory            = 2
  cpu               = 1
  storage           = 10

  db_node_set {
    role                 = "Primary"
    zone                 = var.availability_zone
    dedicated_cluster_id = "cluster-262n63e8"
  }

  db_node_set {
    zone                 = var.availability_zone
    dedicated_cluster_id = "cluster-262n63e8"
  }

  tags = {
    CreateBy = "terraform"
  }
}
```

Create pgsql with kms key

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-6"
}

resource "tencentcloud_postgresql_instance" "example" {
  name              = "tf_postsql_instance"
  availability_zone = var.availability_zone
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = "vpc-86v957zb"
  subnet_id         = "subnet-enm92y0m"
  db_major_version  = "11"
  engine_version    = "11.12"
  db_kernel_version = "v11.12_r1.3"
  need_support_tde  = 1
  kms_key_id        = "788c606a-c7b7-11ec-82d1-5254001e5c4e"
  kms_region        = "ap-guangzhou"
  root_password     = "Root123$"
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

Upgrade kernel version

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-6"
}

resource "tencentcloud_postgresql_instance" "example" {
  name                 = "tf_postsql_instance_update_kernel"
  availability_zone    = var.availability_zone
  charge_type          = "POSTPAID_BY_HOUR"
  vpc_id               = "vpc-86v957zb"
  subnet_id            = "subnet-enm92y0m"
  engine_version       = "13.3"
  root_password        = "Root123$"
  charset              = "LATIN1"
  project_id           = 0
  public_access_switch = false
  security_groups      = ["sg-cm7fbbf3"]
  memory               = 4
  storage              = 250
  
  backup_plan {
    min_backup_start_time        = "01:10:11"
    max_backup_start_time        = "02:10:11"
    base_backup_retention_period = 5
    backup_period                = ["monday", "thursday", "sunday"]
  }

  db_kernel_version = "v13.3_r1.4" # eg:from v13.3_r1.1 to v13.3_r1.4

  tags = {
    tf = "test"
  }
}
```

Import

postgresql instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_postgresql_instance.example postgres-cda1iex1
```