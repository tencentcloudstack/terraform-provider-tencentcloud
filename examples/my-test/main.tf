terraform {
  required_providers {
    tencentcloud = {
      source = "tencentcloudstack/tencentcloud"
    }
  }
}

provider "tencentcloud" {
  region = "ap-guangzhou"
}

variable "default_az" {
  default = "ap-guangzhou-3"
}

data "tencentcloud_vpc_subnets" "gz3" {
  availability_zone = var.default_az
  is_default        = true
}

data "tencentcloud_security_groups" "internal" {
  name = "default"
  tags = var.fixed_tags
}

data "tencentcloud_security_groups" "exclusive" {
  name = "test_preset_sg"
}

variable "fixed_tags" {
  default = {
    fixed_resource : "do_not_remove"
  }
}

locals {
  # local.sg_id
  sg_id  = data.tencentcloud_security_groups.internal.security_groups.0.security_group_id
  sg_id2 = data.tencentcloud_security_groups.exclusive.security_groups.0.security_group_id
}

locals {
  vpc_id    = data.tencentcloud_vpc_subnets.gz3.instance_list.0.vpc_id
  subnet_id = data.tencentcloud_vpc_subnets.gz3.instance_list.0.subnet_id
}

#resource "tencentcloud_sqlserver_instance" "test" {
#  name                   = "tf_sqlserver_instance"
#  availability_zone      = var.default_az
#  charge_type            = "POSTPAID_BY_HOUR"
#  vpc_id                 = local.vpc_id
#  subnet_id              = local.subnet_id
#  security_groups        = [local.sg_id]
#  project_id             = 0
#  memory                 = 2
#  storage                = 10
#  maintenance_week_set   = [1, 2, 3]
#  maintenance_start_time = "09:00"
#  maintenance_time_span  = 3
#  tags                   = {
#    "test" = "test"
#  }
#}

resource "tencentcloud_sqlserver_instance" "test" {
  name                   = "tf_sqlserver_instance_update"
  availability_zone      = var.default_az
  charge_type            = "POSTPAID_BY_HOUR"
  vpc_id                 = "vpc-1yg5ua6l"
  subnet_id              = "subnet-h7av55g8"
  vip                    = "10.0.0.198"
  security_groups        = [local.sg_id]
  memory                 = 4
  storage                = 20
  maintenance_week_set   = [2, 3, 4]
  maintenance_start_time = "08:00"
  maintenance_time_span  = 4
  tags                   = {
    "abc" = "abc"
  }
}