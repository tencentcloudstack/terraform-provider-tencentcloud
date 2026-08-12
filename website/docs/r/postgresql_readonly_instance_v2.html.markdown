---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_readonly_instance_v2"
sidebar_current: "docs-tencentcloud-resource-postgresql_readonly_instance_v2"
description: |-
  Provides a resource to create a PostgreSQL readonly instance V2.
---

# tencentcloud_postgresql_readonly_instance_v2

Provides a resource to create a PostgreSQL readonly instance V2.

## Example Usage

### Create a postpaid PostgreSQL readonly instance

```hcl
resource "tencentcloud_postgresql_readonly_instance_v2" "example" {
  zone                  = "ap-guangzhou-6"
  master_db_instance_id = "postgres-ckwcgdf1"
  spec_code             = "pg.it.small2"
  storage               = 100
  vpc_id                = "vpc-lfnb8yyp"
  subnet_id             = "subnet-csay7wfo"
  instance_charge_type  = "POSTPAID_BY_HOUR"
  name                  = "tf-example"
  project_id            = 1269443
  read_only_group_id    = "pgrogrp-h62yf0re"
  deletion_protection   = false
  security_group_ids = [
    "sg-n8zf5ry9",
    "sg-rs32zv1r"
  ]

  tags = {
    createBy = "Terraform"
  }
}
```

### Create a prepaid PostgreSQL readonly instance

```hcl
resource "tencentcloud_postgresql_readonly_instance_v2" "example" {
  zone                  = "ap-guangzhou-6"
  master_db_instance_id = "postgres-ckwcgdf1"
  spec_code             = "pg.it.small2"
  storage               = 100
  vpc_id                = "vpc-lfnb8yyp"
  subnet_id             = "subnet-csay7wfo"
  instance_charge_type  = "PREPAID"
  period                = 1
  auto_renew_flag       = 1
  name                  = "tf-example"
  project_id            = 1269443
  read_only_group_id    = "pgrogrp-h62yf0re"
  deletion_protection   = true
  security_group_ids = [
    "sg-n8zf5ry9",
    "sg-8gbd3tj9",
    "sg-rs32zv1r"
  ]

  tags = {
    createBy = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `master_db_instance_id` - (Required, String, ForceNew) ID of the primary instance to which the read-only replica belongs.
* `spec_code` - (Required, String, ForceNew) Specification code, which can be obtained via DescribeClasses.
* `storage` - (Required, Int) Instance storage capacity in GB, the step is 10.
* `zone` - (Required, String, ForceNew) Availability zone ID, such as ap-guangzhou-3.
* `auto_renew_flag` - (Optional, Int, ForceNew) Auto renew flag, 0 for manual renew, 1 for auto renew. Default: 0. Only supports PREPAID.
* `auto_voucher` - (Optional, Int) Whether to use voucher automatically, 1 for yes, 0 for no. Default: 0.
* `dedicated_cluster_id` - (Optional, String, ForceNew) Dedicated cluster ID.
* `deletion_protection` - (Optional, Bool) Whether to enable deletion protection, true for enable, false for disable.
* `instance_charge_type` - (Optional, String, ForceNew) Instance billing mode. Valid values: PREPAID, POSTPAID_BY_HOUR. Default: PREPAID.
* `name` - (Optional, String) Instance name, only supports chinese/english/numbers/_/- with length less than 60.
* `need_support_ipv6` - (Optional, Int, ForceNew) Whether to support IPv6 access, 1 for yes, 0 for no. Default: 0.
* `period` - (Optional, Int, ForceNew) Purchase duration in months. PREPAID supports 1,2,3,4,5,6,7,8,9,10,11,12,24,36; POSTPAID only supports 1.
* `project_id` - (Optional, Int) Project ID. Default: 0, means default project.
* `read_only_group_id` - (Optional, String) Read-only group ID.
* `security_group_ids` - (Optional, List: [`String`]) Security group IDs bound to the instance.
* `subnet_id` - (Optional, String, ForceNew) VPC subnet ID, such as subnet-xxxxxxxx.
* `tags` - (Optional, Map) Instance tags.
* `voucher_ids` - (Optional, List: [`String`]) Voucher ID list, currently only one voucher is supported.
* `vpc_id` - (Optional, String, ForceNew) VPC ID, such as vpc-xxxxxxxx.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cpu` - DB instance CPU.
* `db_instance_id` - Instance ID managed by this resource.
* `memory` - DB instance memory.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `20m`) Used when creating the resource.
* `update` - (Defaults to `20m`) Used when updating the resource.
* `delete` - (Defaults to `20m`) Used when deleting the resource.

## Import

PostgreSQL readonly instance V2 can be imported using the id, e.g.

```
terraform import tencentcloud_postgresql_readonly_instance_v2.example pgro-gr4od5iw
```

