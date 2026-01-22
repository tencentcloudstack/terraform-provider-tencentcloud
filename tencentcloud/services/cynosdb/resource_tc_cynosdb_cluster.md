Provide a resource to create a CynosDB cluster.

~> **NOTE:** params `instance_count` and `instance_init_infos` only choose one. If neither parameter is set, the CynosDB cluster is created with parameter `instance_count` set to `2` by default(one RW instance + one Ro instance). If you only need to create a master instance, explicitly set the `instance_count` field to `1`, or configure the RW instance information in the `instance_init_infos` field.

Example Usage

Create a single availability zone NORMAL CynosDB cluster

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create subnet
resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
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

# create cynosdb cluster
resource "tencentcloud_cynosdb_cluster" "example" {
  available_zone               = var.availability_zone
  vpc_id                       = tencentcloud_vpc.vpc.id
  subnet_id                    = tencentcloud_subnet.subnet.id
  db_mode                      = "NORMAL"
  db_type                      = "MYSQL"
  db_version                   = "5.7"
  port                         = 3306
  cluster_name                 = "tf-example"
  password                     = "cynosDB@123"
  instance_maintain_duration   = 7200
  instance_maintain_start_time = 10800
  instance_cpu_core            = 2
  instance_memory_size         = 4
  force_delete                 = false
  instance_maintain_weekdays = [
    "Fri",
    "Mon",
    "Sat",
    "Sun",
    "Thu",
    "Wed",
    "Tue",
  ]

  param_items {
    name          = "character_set_server"
    current_value = "utf8mb4"
  }

  param_items {
    name          = "lower_case_table_names"
    current_value = "0"
  }

  rw_group_sg = [
    tencentcloud_security_group.example.id,
  ]

  ro_group_sg = [
    tencentcloud_security_group.example.id,
  ]

  instance_init_infos {
    cpu            = 2
    memory         = 4
    instance_type  = "rw"
    instance_count = 1
    device_type    = "common"
  }

  instance_init_infos {
    cpu            = 2
    memory         = 4
    instance_type  = "ro"
    instance_count = 1
    device_type    = "exclusive"
  }

  cynos_version = "2.1.14.001"

  tags = {
    createBy = "terraform"
  }
}
```

Create a multiple availability zone SERVERLESS CynosDB cluster

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-4"
}

variable "slave_zone" {
  default = "ap-guangzhou-6"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create subnet
resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
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

# create param template
resource "tencentcloud_cynosdb_param_template" "example" {
  db_mode              = "SERVERLESS"
  engine_version       = "8.0"
  template_name        = "tf-example"
  template_description = "terraform-template"

  param_list {
    current_value = "-1"
    param_name    = "optimizer_trace_offset"
  }
}

# create cynosdb cluster
resource "tencentcloud_cynosdb_cluster" "example" {
  available_zone               = var.availability_zone
  slave_zone                   = var.slave_zone
  vpc_id                       = tencentcloud_vpc.vpc.id
  subnet_id                    = tencentcloud_subnet.subnet.id
  db_mode                      = "SERVERLESS"
  db_type                      = "MYSQL"
  db_version                   = "8.0"
  port                         = 3306
  cluster_name                 = "tf-example"
  password                     = "cynosDB@123"
  instance_maintain_duration   = 7200
  instance_maintain_start_time = 10800
  min_cpu                      = 2
  max_cpu                      = 4
  param_template_id            = tencentcloud_cynosdb_param_template.example.template_id
  force_delete                 = false
  instance_maintain_weekdays   = [
    "Fri",
    "Mon",
    "Sat",
    "Sun",
    "Thu",
    "Wed",
    "Tue",
  ]

  rw_group_sg = [
    tencentcloud_security_group.example.id,
  ]

  ro_group_sg = [
    tencentcloud_security_group.example.id,
  ]

  tags = {
    createBy = "terraform"
  }
}
```

Import

CynosDB cluster can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_cluster.example cynosdbmysql-dzj5l8gz
```
