---
subcategory: "ClickHouse(CDWCH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clickhouse_instance"
sidebar_current: "docs-tencentcloud-resource-clickhouse_instance"
description: |-
  Provides a resource to create a Clickhouse instance.
---

# tencentcloud_clickhouse_instance

Provides a resource to create a Clickhouse instance.

## Example Usage

### Create POSTPAID instance

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-6"
}

data "tencentcloud_clickhouse_spec" "spec" {
  zone       = var.availability_zone
  pay_mode   = "POSTPAID_BY_HOUR"
  is_elastic = false
}

locals {
  data_spec              = [for i in data.tencentcloud_clickhouse_spec.spec.data_spec : i if lookup(i, "cpu") == 4 && lookup(i, "mem") == 16]
  data_spec_name_4c16m   = local.data_spec.0.name
  common_spec            = [for i in data.tencentcloud_clickhouse_spec.spec.common_spec : i if lookup(i, "cpu") == 4 && lookup(i, "mem") == 16]
  common_spec_name_4c16m = local.common_spec.0.name
}

resource "tencentcloud_vpc" "vpc" {
  name       = "cdwch-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "cdwch-subnet"
  cidr_block        = "10.0.0.0/16"
  availability_zone = var.availability_zone
  is_multicast      = false
}

resource "tencentcloud_clickhouse_instance" "example" {
  instance_name       = "tf-example"
  charge_type         = "POSTPAID_BY_HOUR"
  zone                = var.availability_zone
  ha_flag             = true
  vpc_id              = tencentcloud_vpc.vpc.id
  subnet_id           = tencentcloud_subnet.subnet.id
  product_version     = "21.8.12.29"
  ck_default_user_pwd = "Password@123"
  data_spec {
    spec_name = local.data_spec_name_4c16m
    count     = 2
    disk_size = 300
  }

  common_spec {
    spec_name = local.common_spec_name_4c16m
    count     = 3
    disk_size = 300
  }
}
```

### Create PREPAID instance

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-6"
}

data "tencentcloud_clickhouse_spec" "spec" {
  zone       = var.availability_zone
  pay_mode   = "POSTPAID_BY_HOUR"
  is_elastic = false
}

locals {
  data_spec              = [for i in data.tencentcloud_clickhouse_spec.spec.data_spec : i if lookup(i, "cpu") == 4 && lookup(i, "mem") == 16]
  data_spec_name_4c16m   = local.data_spec.0.name
  common_spec            = [for i in data.tencentcloud_clickhouse_spec.spec.common_spec : i if lookup(i, "cpu") == 4 && lookup(i, "mem") == 16]
  common_spec_name_4c16m = local.common_spec.0.name
}

resource "tencentcloud_vpc" "vpc" {
  name       = "cdwch-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "cdwch-subnet"
  cidr_block        = "10.0.0.0/16"
  availability_zone = var.availability_zone
  is_multicast      = false
}

resource "tencentcloud_clickhouse_instance" "example" {
  instance_name       = "tf-example"
  charge_type         = "PREPAID"
  renew_flag          = 1
  time_span           = 1
  zone                = var.availability_zone
  ha_flag             = true
  vpc_id              = tencentcloud_vpc.vpc.id
  subnet_id           = tencentcloud_subnet.subnet.id
  product_version     = "21.8.12.29"
  ck_default_user_pwd = "Password@123"
  data_spec {
    spec_name = local.data_spec_name_4c16m
    count     = 2
    disk_size = 300
  }

  common_spec {
    spec_name = local.common_spec_name_4c16m
    count     = 3
    disk_size = 300
  }
}
```

## Argument Reference

The following arguments are supported:

* `charge_type` - (Required, String) Billing type: `PREPAID` prepaid, `POSTPAID_BY_HOUR` postpaid.
* `data_spec` - (Required, List) Data spec.
* `ha_flag` - (Required, Bool) Whether it is highly available.
* `instance_name` - (Required, String) Instance name.
* `product_version` - (Required, String) Product version.
* `subnet_id` - (Required, String) Subnet.
* `vpc_id` - (Required, String) Private network.
* `zone` - (Required, String) Availability zone.
* `ck_default_user_pwd` - (Optional, String) The password for the default account to log in to the instance. 8-16 characters, including at least three of the following: uppercase letters, lowercase letters, numbers, and special characters `!@#%^*`. The first character cannot be a special character.
* `cls_log_set_id` - (Optional, String) CLS log set id.
* `common_spec` - (Optional, List) ZK node.
* `cos_bucket_name` - (Optional, String) COS bucket name.
* `ha_zk` - (Optional, Bool) Whether ZK is highly available.
* `mount_disk_type` - (Optional, Int) Whether it is mounted on a bare disk.
* `renew_flag` - (Optional, Int) PREPAID needs to be passed. Whether to renew automatically. 1 means auto renewal is enabled.
* `secondary_zone_info` - (Optional, List) Secondary zone info.
* `tags` - (Optional, Map) Tag description list.
* `time_span` - (Optional, Int) Prepaid needs to be delivered, billing time length, how many months.

The `common_spec` object supports the following:

* `count` - (Required, Int) Node count. NOTE: Only support value 3.
* `disk_size` - (Required, Int) Disk size.
* `spec_name` - (Required, String) Spec name.

The `data_spec` object supports the following:

* `count` - (Required, Int) Data spec count.
* `disk_size` - (Required, Int) Disk size.
* `spec_name` - (Required, String) Spec name.

The `secondary_zone_info` object supports the following:

* `secondary_subnet` - (Optional, String) Secondary subnet.
* `secondary_zone` - (Optional, String) Secondary zone.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `access_info` - access address info.
* `expire_time` - Expire time.


## Import

Clickhouse instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_clickhouse_instance.example cdwch-4l6mm8p7
```

