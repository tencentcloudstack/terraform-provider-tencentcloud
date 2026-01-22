Use this resource to create postgresql parameter.

Example Usage

```hcl
variable "default_az" {
  default = "ap-guangzhou-3"
}

data "tencentcloud_vpc_subnets" "gz3" {
  availability_zone = var.default_az
  is_default = true
}

locals {
  vpc_id = data.tencentcloud_vpc_subnets.gz3.instance_list.0.vpc_id
  subnet_id = data.tencentcloud_vpc_subnets.gz3.instance_list.0.subnet_id
}

data "tencentcloud_availability_zones_by_product" "zone" {
  product = "postgres"
}

resource "tencentcloud_postgresql_instance" "test" {
  name 				= "tf_postsql_postpaid"
  availability_zone = var.default_az
  charge_type 		= "POSTPAID_BY_HOUR"
  period            = 1
  vpc_id  	  		= local.vpc_id
  subnet_id 		= local.subnet_id
  engine_version	= "13.3"
  root_password	    = "t1qaA2k1wgvfa3?ZZZ"
  security_groups   = ["sg-5275dorp"]
  charset			= "LATIN1"
  project_id 		= 0
  memory 			= 2
  storage 			= 20
}
resource "tencentcloud_postgresql_parameters" "postgresql_parameters" {
  db_instance_id = tencentcloud_postgresql_instance.test.id
  param_list {
    expected_value = "off"
    name           = "check_function_bodies"
  }
}
```

Import

postgresql parameters can be imported, e.g.

```
$ terraform import tencentcloud_postgresql_parameters.example pgrogrp-lckioi2a
```
