---
subcategory: "TDSQL-C for MySQL(tdmysql)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmysql_db_instance"
sidebar_current: "docs-tencentcloud-resource-tdmysql_db_instance"
description: |-
  Provides a resource to create a TDSQL-C for MySQL (tdmysql) database instance.
---

# tencentcloud_tdmysql_db_instance

Provides a resource to create a TDSQL-C for MySQL (tdmysql) database instance.

## Example Usage

```hcl
resource "tencentcloud_tdmysql_db_instance" "example" {
  zone               = "ap-guangzhou-6"
  vpc_id             = "vpc-i5yyodl9"
  subnet_id          = "subnet-hhi88a58"
  spec_code          = "8c16g"
  disk               = 400
  storage_node_num   = 3
  replications       = 3
  create_version     = "21.6.4.3"
  instance_name      = "tf-example"
  storage_node_cpu   = 8
  storage_node_mem   = 16
  pay_mode           = "0"
  vport              = 3306
  zones              = ["ap-guangzhou-6"]
  instance_type      = "hybrid"
  storage_type       = "CLOUD_HSSD"
  instance_mode      = "enhanced"
  sql_mode           = "MySQL"
  security_group_ids = ["sg-8gbd3tj9", "sg-2g6p85pr", "sg-hvnj11z7"]
  password           = "Password@2026"
  reserved_space_gb  = 5
  auto_renew_flag    = 0
  enable_ssl         = true

  init_params {
    param = "character_set_server"
    value = "utf8mb4"
  }

  init_params {
    param = "lower_case_table_names"
    value = "0"
  }

  init_params {
    param = "collation_server"
    value = "utf8mb4_0900_ai_ci"
  }

  init_params {
    param = "tdstore_enter_readonly_threshold_percent"
    value = "99.00"
  }

  resource_tags {
    tag_key   = "CreatedBy"
    tag_value = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `disk` - (Required, Int) Storage node disk capacity, unit GB.
* `instance_name` - (Required, String) Instance name, length 1-60.
* `replications` - (Required, Int, ForceNew) Storage node replica number, max 5, must be odd.
* `spec_code` - (Required, String) Specification code.
* `storage_node_num` - (Required, Int) Storage node number.
* `subnet_id` - (Required, String) Subnet ID.
* `vpc_id` - (Required, String) VPC ID.
* `zone` - (Required, String, ForceNew) Instance zone.
* `auto_renew_flag` - (Optional, Int) Auto renew flag, 1 indicates enabling auto-renewal; 0 indicates disabling auto-renewal.
* `auto_scale_config` - (Optional, List, ForceNew) Auto scaling config for svls instance.
* `auto_voucher` - (Optional, Bool) Whether to use voucher.
* `az_mode` - (Optional, Int) AZ mode, 1: single AZ, 2: multi AZ non-master, 3: multi AZ master.
* `create_version` - (Optional, String, ForceNew) Create version, defaults to latest.
* `enable_ssl` - (Optional, Bool) Whether to enable SSL.
* `encryption_enable` - (Optional, Int, ForceNew) Transparent encryption, 0: disable, 1: enable.
* `full_replications` - (Optional, Int) Full replica number.
* `init_params` - (Optional, Set) Init instance params.
* `instance_mode` - (Optional, String, ForceNew) Instance mode.
* `instance_type` - (Optional, String, ForceNew) Instance architecture type, separate or hybrid.
* `mc_num` - (Optional, Int, ForceNew) Control node number.
* `password` - (Optional, String) dbaadmin password.
* `pay_mode` - (Optional, String, ForceNew) Pay mode, 0: postpaid, 1: prepaid.
* `reserved_space_gb` - (Optional, Float64) Data reserved space of each storage node, unit GB.
* `resource_tags` - (Optional, List, ForceNew) Resource tag list.
* `security_group_ids` - (Optional, List: [`String`]) Security group ID list.
* `sql_mode` - (Optional, String, ForceNew) Compatible mode, MySQL or HBase.
* `storage_node_cpu` - (Optional, Int) Storage node CPU cores.
* `storage_node_mem` - (Optional, Int) Storage node memory.
* `storage_type` - (Optional, String) Disk type, CLOUD_HSSD or CLOUD_TCS.
* `template_id` - (Optional, String, ForceNew) Parameter template ID.
* `time_span` - (Optional, Int, ForceNew) Time span.
* `time_unit` - (Optional, String, ForceNew) Time unit, m: month.
* `user_name` - (Optional, String, ForceNew) Root username, defaults to dbaadmin.
* `voucher_ids` - (Optional, List: [`String`]) Voucher ID list.
* `vport` - (Optional, Int) Custom port.
* `zones` - (Optional, List: [`String`]) Multi AZ zone list.

The `auto_scale_config` object supports the following:

* `range_max` - (Required, Float64, ForceNew) CCU max value.
* `range_min` - (Required, Float64, ForceNew) CCU min value.

The `init_params` object supports the following:

* `param` - (Required, String) Param key.
* `value` - (Required, String) Param value.

The `resource_tags` object supports the following:

* `tag_key` - (Required, String, ForceNew) Tag key.
* `tag_value` - (Optional, String, ForceNew) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `instance_id` - Instance ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `20m`) Used when creating the resource.
* `update` - (Defaults to `20m`) Used when updating the resource.
* `delete` - (Defaults to `20m`) Used when deleting the resource.
## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateDBInstances](https://www.tencentcloud.com/search/?wd=CreateDBInstances&type=document&lang=en)
  - [DescribeFlow](https://www.tencentcloud.com/search/?wd=DescribeFlow&type=document&lang=en)
  - [ModifyAutoRenewFlag](https://www.tencentcloud.com/search/?wd=ModifyAutoRenewFlag&type=document&lang=en)
  - [ModifyInstanceDataReservedSpace](https://www.tencentcloud.com/search/?wd=ModifyInstanceDataReservedSpace&type=document&lang=en)
  - [ModifyInstanceSSLStatus](https://www.tencentcloud.com/search/?wd=ModifyInstanceSSLStatus&type=document&lang=en)
* Read:
  - [DescribeDBInstanceDetail](https://www.tencentcloud.com/search/?wd=DescribeDBInstanceDetail&type=document&lang=en)
  - [DescribeDBSecurityGroups](https://www.tencentcloud.com/search/?wd=DescribeDBSecurityGroups&type=document&lang=en)
  - [DescribeInstanceDataReservedSpace](https://www.tencentcloud.com/search/?wd=DescribeInstanceDataReservedSpace&type=document&lang=en)
  - [DescribeInstanceSSLStatus](https://www.tencentcloud.com/search/?wd=DescribeInstanceSSLStatus&type=document&lang=en)
* Update:
  - [DescribeDBInstanceDetail](https://www.tencentcloud.com/search/?wd=DescribeDBInstanceDetail&type=document&lang=en)
  - [DescribeFlow](https://www.tencentcloud.com/search/?wd=DescribeFlow&type=document&lang=en)
  - [ExpandInstance](https://www.tencentcloud.com/search/?wd=ExpandInstance&type=document&lang=en)
  - [ModifyAutoRenewFlag](https://www.tencentcloud.com/search/?wd=ModifyAutoRenewFlag&type=document&lang=en)
  - [ModifyDBInstance](https://www.tencentcloud.com/search/?wd=ModifyDBInstance&type=document&lang=en)
  - [ModifyDBInstanceSecurityGroups](https://www.tencentcloud.com/search/?wd=ModifyDBInstanceSecurityGroups&type=document&lang=en)
  - [ModifyDBInstanceVPort](https://www.tencentcloud.com/search/?wd=ModifyDBInstanceVPort&type=document&lang=en)
  - [ModifyDBParameters](https://www.tencentcloud.com/search/?wd=ModifyDBParameters&type=document&lang=en)
  - [ModifyInstanceDataReservedSpace](https://www.tencentcloud.com/search/?wd=ModifyInstanceDataReservedSpace&type=document&lang=en)
  - [ModifyInstanceName](https://www.tencentcloud.com/search/?wd=ModifyInstanceName&type=document&lang=en)
  - [ModifyInstanceNetwork](https://www.tencentcloud.com/search/?wd=ModifyInstanceNetwork&type=document&lang=en)
  - [ModifyInstanceSSLStatus](https://www.tencentcloud.com/search/?wd=ModifyInstanceSSLStatus&type=document&lang=en)
  - [ResetUsersPassword](https://www.tencentcloud.com/search/?wd=ResetUsersPassword&type=document&lang=en)
* Delete:
  - [DescribeDBInstanceDetail](https://www.tencentcloud.com/search/?wd=DescribeDBInstanceDetail&type=document&lang=en)
  - [DestroyInstances](https://www.tencentcloud.com/search/?wd=DestroyInstances&type=document&lang=en)
  - [IsolateDBInstance](https://www.tencentcloud.com/search/?wd=IsolateDBInstance&type=document&lang=en)

## Import

TDSQL-C for MySQL database instance can be imported using the id, e.g.

```
terraform import tencentcloud_tdmysql_db_instance.example tdsql3-gpyfw59om6
```

