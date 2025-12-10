Provides a resource to create a CDWDoris instance

Example Usage

Create a POSTPAID instance

```hcl
# availability zone
variable "availability_zone" {
  default = "ap-guangzhou-6"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "172.16.0.0/16"
}

# create subnet
resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "172.16.0.0/24"
  is_multicast      = false
}

# create security group
resource "tencentcloud_security_group" "example" {
  name        = "tf-example"
  description = "security group desc."

  tags = {
    "createBy" = "Terraform"
  }
}

# create POSTPAID instance
resource "tencentcloud_cdwdoris_instance" "example" {
  zone                  = var.availability_zone
  user_vpc_id           = tencentcloud_vpc.vpc.id
  user_subnet_id        = tencentcloud_subnet.subnet.id
  product_version       = "2.1"
  instance_name         = "tf-example"
  doris_user_pwd        = "Password@test"
  ha_flag               = true
  ha_type               = 1
  case_sensitive        = 0
  enable_multi_zones    = false
  workload_group_status = "open"

  security_group_ids = [
    tencentcloud_security_group.example.id
  ]

  charge_properties {
    charge_type = "POSTPAID_BY_HOUR"
  }

  fe_spec {
    spec_name = "S_4_16_P"
    count     = 3
    disk_size = 200
  }

  be_spec {
    spec_name = "S_4_16_P"
    count     = 3
    disk_size = 200
  }

  tags {
    tag_key   = "createBy"
    tag_value = "Terraform"
  }
}
```

Create a PREPAID instance

```hcl
# availability zone
variable "availability_zone" {
  default = "ap-guangzhou-6"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "172.16.0.0/16"
}

# create subnet
resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "172.16.0.0/24"
  is_multicast      = false
}

# create security group
resource "tencentcloud_security_group" "example" {
  name        = "tf-example"
  description = "security group desc."

  tags = {
    "createBy" = "Terraform"
  }
}

# create PREPAID instance
resource "tencentcloud_cdwdoris_instance" "example" {
  zone                  = var.availability_zone
  user_vpc_id           = tencentcloud_vpc.vpc.id
  user_subnet_id        = tencentcloud_subnet.subnet.id
  product_version       = "2.1"
  instance_name         = "tf-example"
  doris_user_pwd        = "Password@test"
  ha_flag               = true
  ha_type               = 1
  case_sensitive        = 0
  enable_multi_zones    = false
  workload_group_status = "close"

  security_group_ids = [
    tencentcloud_security_group.example.id
  ]

  charge_properties {
    charge_type = "PREPAID"
    time_span   = 1
    time_unit   = "m"
  }

  fe_spec {
    spec_name = "S_4_16_P"
    count     = 3
    disk_size = 200
  }

  be_spec {
    spec_name = "S_4_16_P"
    count     = 3
    disk_size = 200
  }

  tags {
    tag_key   = "createBy"
    tag_value = "Terraform"
  }
}
```
