---
subcategory: "TcaplusDB"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcaplus_idl"
sidebar_current: "docs-tencentcloud-resource-tcaplus_idl"
description: |-
  Use this resource to create TcaplusDB IDL file.
---

# tencentcloud_tcaplus_idl

Use this resource to create TcaplusDB IDL file.

## Example Usage

### Create a tcaplus database idl file

The file will be with a specified cluster and tablegroup.

```hcl
locals {
  vpc_id    = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  subnet_id = data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id
}

variable "availability_zone" {
  default = "ap-guangzhou-3"
}

data "tencentcloud_vpc_subnets" "vpc" {
  is_default        = true
  availability_zone = var.availability_zone
}

resource "tencentcloud_tcaplus_cluster" "example" {
  idl_type                 = "PROTO"
  cluster_name             = "tf_example_tcaplus_cluster"
  vpc_id                   = local.vpc_id
  subnet_id                = local.subnet_id
  password                 = "your_pw_123111"
  old_password_expire_last = 3600
}

resource "tencentcloud_tcaplus_tablegroup" "example" {
  cluster_id      = tencentcloud_tcaplus_cluster.example.id
  tablegroup_name = "tf_example_group_name"
}

resource "tencentcloud_tcaplus_idl" "main" {
  cluster_id    = tencentcloud_tcaplus_cluster.example.id
  tablegroup_id = tencentcloud_tcaplus_tablegroup.example.id
  file_name     = "tf_example_tcaplus_idl"
  file_type     = "PROTO"
  file_ext_type = "proto"
  file_content  = <<EOF
    syntax = "proto2";
    package myTcaplusTable;
    import "tcaplusservice.optionv1.proto";
    message tb_online {
        option(tcaplusservice.tcaplus_primary_key) = "uin,name,region";
        required int64 uin = 1;
        required string name = 2;
        required int32 region = 3;
        required int32 gamesvrid = 4;
        optional int32 logintime = 5 [default = 1];
        repeated int64 lockid = 6 [packed = true];
        optional bool is_available = 7 [default = false];
        optional pay_info pay = 8;
    }

    message pay_info {
        required int64 pay_id = 1;
        optional uint64 total_money = 2;
        optional uint64 pay_times = 3;
        optional pay_auth_info auth = 4;
        message pay_auth_info {
            required string pay_keys = 1;
            optional int64 update_time = 2;
        }
    }
    EOF
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) ID of the TcaplusDB cluster to which the table group belongs.
* `file_content` - (Required, String, ForceNew) IDL file content of the TcaplusDB table.
* `file_ext_type` - (Required, String, ForceNew) File ext type of the IDL file. If `file_type` is `PROTO`, `file_ext_type` must be 'proto'; If `file_type` is `TDR`, `file_ext_type` must be 'xml'.
* `file_name` - (Required, String, ForceNew) Name of the IDL file.
* `file_type` - (Required, String, ForceNew) Type of the IDL file. Valid values are PROTO and TDR.
* `tablegroup_id` - (Required, String, ForceNew) ID of the table group to which the IDL file belongs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `table_infos` - Table info of the IDL.
  * `error` - Error messages for creating IDL file.
  * `index_key_set` - Index key set of the TcaplusDB table.
  * `key_fields` - Primary key fields of the TcaplusDB table.
  * `sum_key_field_size` - Total size of primary key field of the TcaplusDB table.
  * `sum_value_field_size` - Total size of non-primary key fields of the TcaplusDB table.
  * `table_name` - Name of the TcaplusDB table.
  * `value_fields` - Non-primary key fields of the TcaplusDB table.


