---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcaplus_idl"
sidebar_current: "docs-tencentcloud-resource-tcaplus_idl"
description: |-
  Use this resource to create tcaplus idl file
---

# tencentcloud_tcaplus_idl

Use this resource to create tcaplus idl file

## Example Usage

```hcl
resource "tencentcloud_tcaplus_cluster" "test" {
  idl_type                 = "PROTO"
  cluster_name             = "tf_tcaplus_cluster_test"
  vpc_id                   = "vpc-7k6gzox6"
  subnet_id                = "subnet-akwgvfa3"
  password                 = "1qaA2k1wgvfa3ZZZ"
  old_password_expire_last = 3600
}

resource "tencentcloud_tcaplus_group" "group" {
  cluster_id = tencentcloud_tcaplus_cluster.test.id
  group_name = "tf_test_group_name"
}

resource "tencentcloud_tcaplus_idl" "main" {
  cluster_id    = tencentcloud_tcaplus_cluster.test.id
  group_id      = tencentcloud_tcaplus_group.group.id
  file_name     = "tf_idl_test"
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

* `cluster_id` - (Required, ForceNew) Cluster id of the idl belongs..
* `file_content` - (Required, ForceNew) Idl file content.
* `file_ext_type` - (Required, ForceNew) File ext type of this idl file. if `file_type` is PROTO  `file_ext_type` must be 'proto',if `file_type` is TDR  `file_ext_type` must be 'xml',if `file_type` is MIX  `file_ext_type` must be 'xml' or 'proto'.
* `file_name` - (Required, ForceNew) Name of this idl file.
* `file_type` - (Required, ForceNew) Type of this idl file, Valid values are PROTO,TDR,MIX.
* `group_id` - (Required, ForceNew) Group of this idl belongs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `table_infos` - Table infos in this idl.
  * `error` - Show if this table  error.
  * `index_key_set` - Index key set of this table.
  * `key_fields` - Key fields of this table.
  * `sum_key_field_size` - Key fields size of this table.
  * `sum_value_field_size` - Value fields size of this table.
  * `table_name` - Name of this table.
  * `value_fields` - Value fields of this table.


