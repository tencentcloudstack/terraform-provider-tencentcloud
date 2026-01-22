---
subcategory: "CdwDoris"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cdwdoris_workload_group"
sidebar_current: "docs-tencentcloud-resource-cdwdoris_workload_group"
description: |-
  Provides a resource to create a CDWDoris workload group
---

# tencentcloud_cdwdoris_workload_group

Provides a resource to create a CDWDoris workload group

~> **NOTE:** To use this resource, The `workload_group_status` field of `tencentcloud_cdwdoris_instance` needs to be set to true.

## Example Usage

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

# create instance
resource "tencentcloud_cdwdoris_instance" "example" {
  zone                  = var.availability_zone
  user_vpc_id           = tencentcloud_vpc.vpc.id
  user_subnet_id        = tencentcloud_subnet.subnet.id
  product_version       = "2.1"
  instance_name         = "tf-example"
  doris_user_pwd        = "Password@test"
  ha_flag               = false
  case_sensitive        = 0
  enable_multi_zones    = false
  workload_group_status = "open"

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

# create workload group
resource "tencentcloud_cdwdoris_workload_group" "example" {
  instance_id = tencentcloud_cdwdoris_instance.example.id
  workload_group {
    workload_group_name       = "example"
    cpu_share                 = 1024
    memory_limit              = 20
    enable_memory_over_commit = true
    cpu_hard_limit            = "30%"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance id.
* `workload_group` - (Optional, List) Resource group configuration.

The `workload_group` object supports the following:

* `cpu_hard_limit` - (Optional, String) Cpu hard limit. Note: This field may return null, indicating that no valid value can be obtained.
* `cpu_share` - (Optional, Int) CPU weight. Note: This field may return null, indicating that no valid value can be obtained.
* `enable_memory_over_commit` - (Optional, Bool) Whether to allow over-allocation. Note: This field may return null, indicating that no valid value can be obtained.
* `memory_limit` - (Optional, Int) Memory limit, the sum of the memory limit values of all resource groups should be less than or equal to 100. Note: This field may return null, indicating that no valid value can be obtained.
* `workload_group_name` - (Optional, String, ForceNew) Workload group name. Note: This field may return null, indicating that no valid value can be obtained.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cdwdoris workload group can be imported using the id, e.g.

```
terraform import tencentcloud_cdwdoris_workload_group.example cdwdoris-rhbflamd#tf-example
```

