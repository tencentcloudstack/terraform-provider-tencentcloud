Use this resource to create postgresql readonly instance.

Example Usage

Create postgresql readonly instance

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
  root_user         = "root123"
  root_password     = "Root123$"
  charset           = "UTF8"
  project_id        = 0
  memory            = 2
  cpu               = 1
  storage           = 10

  tags = {
    test = "tf"
  }
}

# create postgresql readonly group
resource "tencentcloud_postgresql_readonly_group" "example" {
  master_db_instance_id       = tencentcloud_postgresql_instance.example.id
  name                        = "tf_ro_group"
  project_id                  = 0
  vpc_id                      = tencentcloud_vpc.vpc.id
  subnet_id                   = tencentcloud_subnet.subnet.id
  replay_lag_eliminate        = 1
  replay_latency_eliminate    = 1
  max_replay_lag              = 100
  max_replay_latency          = 512
  min_delay_eliminate_reserve = 1
}

# create security group
resource "tencentcloud_security_group" "example" {
  name        = "tf-example"
  description = "sg desc."
  project_id  = 0

  tags = {
    "example" = "test"
  }
}

# create postgresql readonly instance
resource "tencentcloud_postgresql_readonly_instance" "example" {
  read_only_group_id    = tencentcloud_postgresql_readonly_group.example.id
  master_db_instance_id = tencentcloud_postgresql_instance.example.id
  zone                  = var.availability_zone
  name                  = "example"
  auto_renew_flag       = 0
  db_version            = "10.23"
  instance_charge_type  = "POSTPAID_BY_HOUR"
  memory                = 4
  cpu                   = 2
  storage               = 250
  vpc_id                = tencentcloud_vpc.vpc.id
  subnet_id             = tencentcloud_subnet.subnet.id
  need_support_ipv6     = 0
  project_id            = 0
  security_groups_ids   = [
    tencentcloud_security_group.example.id,
  ]
}
```

Create postgresql readonly instance of CDC

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

# create postgresql readonly group
resource "tencentcloud_postgresql_readonly_group" "example" {
  master_db_instance_id       = tencentcloud_postgresql_instance.example.id
  name                        = "tf_ro_group"
  project_id                  = 0
  vpc_id                      = tencentcloud_vpc.vpc.id
  subnet_id                   = tencentcloud_subnet.subnet.id
  replay_lag_eliminate        = 1
  replay_latency_eliminate    = 1
  max_replay_lag              = 100
  max_replay_latency          = 512
  min_delay_eliminate_reserve = 1
}

# create security group
resource "tencentcloud_security_group" "example" {
  name        = "tf-example"
  description = "sg desc."
  project_id  = 0

  tags = {
    CreateBy = "terraform"
  }
}

# create postgresql readonly instance
resource "tencentcloud_postgresql_readonly_instance" "example" {
  read_only_group_id    = tencentcloud_postgresql_readonly_group.example.id
  master_db_instance_id = tencentcloud_postgresql_instance.example.id
  zone                  = var.availability_zone
  name                  = "example"
  auto_renew_flag       = 0
  db_version            = "10.23"
  instance_charge_type  = "POSTPAID_BY_HOUR"
  memory                = 4
  cpu                   = 2
  storage               = 250
  vpc_id                = tencentcloud_vpc.vpc.id
  subnet_id             = tencentcloud_subnet.subnet.id
  need_support_ipv6     = 0
  project_id            = 0
  dedicated_cluster_id  = "cluster-262n63e8"
  security_groups_ids = [
    tencentcloud_security_group.example.id,
  ]
}
```

Import

postgresql readonly instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_postgresql_readonly_instance.example pgro-gih5m0ke
```