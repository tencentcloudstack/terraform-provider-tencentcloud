---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_readonly_instance_v2"
sidebar_current: "docs-tencentcloud-resource-postgresql_readonly_instance_v2"
description: |-
  Provides a resource to create a PostgreSQL readonly instance.
---

# tencentcloud_postgresql_readonly_instance_v2

Provides a resource to create a PostgreSQL readonly instance.

## Example Usage

### Create a postgresql readonly instance

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

# create postgresql primary instance
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

# create postgresql readonly instance v2
resource "tencentcloud_postgresql_readonly_instance_v2" "example" {
  zone                  = var.availability_zone
  master_db_instance_id = tencentcloud_postgresql_instance.example.id
  spec_code             = "pgro.c2.large.ha"
  storage               = 250
  instance_count        = 1
  period                = 1
  instance_charge_type  = "POSTPAID_BY_HOUR"
  vpc_id                = tencentcloud_vpc.vpc.id
  subnet_id             = tencentcloud_subnet.subnet.id
  read_only_group_id    = tencentcloud_postgresql_readonly_group.example.id
  name                  = "tf_ro_instance_v2"
  auto_renew_flag       = 0
  need_support_ipv6     = 0
  project_id            = 0

  timeouts {
    create = "60m"
    delete = "60m"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_count` - (Required, Int, ForceNew) Number of instances to purchase, value range: [1-6]. Only the first instance ID is managed by this resource.
* `master_db_instance_id` - (Required, String, ForceNew) ID of the primary instance to which the read-only replica belongs.
* `period` - (Required, Int, ForceNew) Purchase duration in months. PREPAID supports 1,2,3,4,5,6,7,8,9,10,11,12,24,36; POSTPAID only supports 1.
* `spec_code` - (Required, String, ForceNew) Specification code, which can be obtained via DescribeClasses.
* `storage` - (Required, Int, ForceNew) Instance storage capacity in GB, the step is 10.
* `zone` - (Required, String, ForceNew) Availability zone ID, such as ap-guangzhou-3.
* `activity_id` - (Optional, Int) Activity ID.
* `auto_renew_flag` - (Optional, Int) Auto renew flag, 0 for manual renew, 1 for auto renew. Default: 0. Only supports PREPAID.
* `auto_voucher` - (Optional, Int) Whether to use voucher automatically, 1 for yes, 0 for no. Default: 0.
* `db_version` - (Optional, String) PostgreSQL kernel version, no longer needed, it will keep the same as the primary instance.
* `dedicated_cluster_id` - (Optional, String) Dedicated cluster ID.
* `deletion_protection` - (Optional, Bool) Whether to enable deletion protection, true for enable, false for disable.
* `instance_charge_type` - (Optional, String) Instance billing mode. Valid values: PREPAID, POSTPAID_BY_HOUR. Default: PREPAID.
* `name` - (Optional, String) Instance name, only supports chinese/english/numbers/_/- with length less than 60.
* `need_support_ipv6` - (Optional, Int) Whether to support IPv6 access, 1 for yes, 0 for no. Default: 0.
* `project_id` - (Optional, Int) Project ID. Default: 0, means default project.
* `read_only_group_id` - (Optional, String) Read-only group ID.
* `security_group_ids` - (Optional, List: [`String`]) Security group IDs bound to the instance.
* `subnet_id` - (Optional, String) VPC subnet ID, such as subnet-xxxxxxxx.
* `tag_list` - (Optional, List) Instance tag info (legacy, single tag). It is recommended to use the new field `tags`.
* `tags` - (Optional, Map) Instance tags.
* `voucher_ids` - (Optional, List: [`String`]) Voucher ID list, currently only one voucher is supported.
* `vpc_id` - (Optional, String) VPC ID, such as vpc-xxxxxxxx.

The `tag_list` object supports the following:

* `tag_key` - (Required, String) Tag key.
* `tag_value` - (Optional, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `bill_id` - Frozen flow number.
* `billing_parameters` - Billing parameters for order placement, only returned when the input parameter BillingParameters has a value.
* `db_instance_id_set` - Created instance ID set, only returned in POSTPAID scenario.
* `db_instance_id` - Instance ID managed by this resource.
* `deal_names` - Order number list, each instance corresponds to one order.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `1h0m`) Used when creating the resource.
* `delete` - (Defaults to `1h0m`) Used when deleting the resource.

## Import

postgresql readonly instance can be imported using the instance id, e.g.

```
$ terraform import tencentcloud_postgresql_readonly_instance_v2.example pgro-gih5m0ke
```

