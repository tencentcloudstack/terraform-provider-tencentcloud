---
subcategory: "CdwDoris"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cdwdoris_instance"
sidebar_current: "docs-tencentcloud-resource-cdwdoris_instance"
description: |-
  Provides a resource to create a CDWDoris instance
---

# tencentcloud_cdwdoris_instance

Provides a resource to create a CDWDoris instance

## Example Usage

### Create a POSTPAID instance

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

### Create a PREPAID instance

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

## Argument Reference

The following arguments are supported:

* `be_spec` - (Required, List) BE specifications.
* `charge_properties` - (Required, List) Payment type.
* `doris_user_pwd` - (Required, String) Database password.
* `fe_spec` - (Required, List) FE specifications.
* `ha_flag` - (Required, Bool) Whether it is highly available.
* `instance_name` - (Required, String) Instance name.
* `product_version` - (Required, String) Product version number.
* `user_subnet_id` - (Required, String) User subnet ID.
* `user_vpc_id` - (Required, String) User VPCID.
* `workload_group_status` - (Required, String) Whether to enable resource group. `open` - enable, `close` - disable.
* `zone` - (Required, String) Availability zone.
* `case_sensitive` - (Optional, Int) Whether the table name is case sensitive, 0 refers to sensitive, 1 refers to insensitive, compared in lowercase; 2 refers to insensitive, and the table name is changed to lowercase for storage.
* `enable_multi_zones` - (Optional, Bool) Whether to enable multi-availability zone.
* `ha_type` - (Optional, Int) High availability type: 0 indicates non-high availability (only one FE, FeSpec.CreateInstanceSpec.Count=1), 1 indicates read high availability (at least 3 FEs must be deployed, FeSpec.CreateInstanceSpec.Count>=3, and it must be an odd number), 2 indicates read and write high availability (at least 5 FEs must be deployed, FeSpec.CreateInstanceSpec.Count>=5, and it must be an odd number).
* `security_group_ids` - (Optional, List: [`String`]) Security Group Id list.
* `tags` - (Optional, List) Tag list.
* `user_multi_zone_infos` - (Optional, List) After the Multi-AZ is enabled, all user's Availability Zones and Subnets information are shown.

The `be_spec` object supports the following:

* `count` - (Required, Int) Quantities.
* `disk_size` - (Required, Int) Cloud disk size.
* `spec_name` - (Required, String) Specification name.

The `charge_properties` object supports the following:

* `charge_type` - (Optional, String) Billing type: `PREPAID` for prepayment, and `POSTPAID_BY_HOUR` for postpayment. Note: This field may return null, indicating that no valid values can be obtained.
* `renew_flag` - (Optional, Int) Whether to automatically renew. 1 means automatic renewal is enabled. Note: This field may return null, indicating that no valid values can be obtained.
* `time_span` - (Optional, Int) Billing duration Note: This field may return null, indicating that no valid values can be obtained.
* `time_unit` - (Optional, String) Billing time unit, and `m` means month, etc. Note: This field may return null, indicating that no valid values can be obtained.

The `fe_spec` object supports the following:

* `count` - (Required, Int) Quantities.
* `disk_size` - (Required, Int) Cloud disk size.
* `spec_name` - (Required, String) Specification name.

The `tags` object supports the following:

* `tag_key` - (Required, String) Tag key.
* `tag_value` - (Required, String) Tag value.

The `user_multi_zone_infos` object supports the following:

* `subnet_id` - (Optional, String) Subnet ID Note: This field may return null, indicating that no valid values can be obtained.
* `subnet_ip_num` - (Optional, Int) The number of available IP addresses in the current subnet Note: This field may return null, indicating that no valid values can be obtained.
* `zone` - (Optional, String) Availability zone Note: This field may return null, indicating that no valid values can be obtained.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



