---
subcategory: "TcaplusDB"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcaplus_table"
sidebar_current: "docs-tencentcloud-resource-tcaplus_table"
description: |-
  Use this resource to create TcaplusDB table.
---

# tencentcloud_tcaplus_table

Use this resource to create TcaplusDB table.

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

resource "tencentcloud_tcaplus_tablegroup" "tablegroup" {
  cluster_id      = tencentcloud_tcaplus_cluster.test.id
  tablegroup_name = "tf_test_group_name"
}

resource "tencentcloud_tcaplus_idl" "main" {
  cluster_id    = tencentcloud_tcaplus_cluster.test.id
  tablegroup_id = tencentcloud_tcaplus_tablegroup.tablegroup.id
  file_name     = "tf_idl_test_2"
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

resource "tencentcloud_tcaplus_table" "table" {
  cluster_id        = tencentcloud_tcaplus_cluster.test.id
  tablegroup_id     = tencentcloud_tcaplus_tablegroup.tablegroup.id
  table_name        = "tb_online"
  table_type        = "GENERIC"
  description       = "test"
  idl_id            = tencentcloud_tcaplus_idl.main.id
  table_idl_type    = "PROTO"
  reserved_read_cu  = 1000
  reserved_write_cu = 20
  reserved_volume   = 1
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, ForceNew) ID of the TcaplusDB cluster to which the table belongs.
* `idl_id` - (Required) ID of the IDL File.
* `reserved_read_cu` - (Required, ForceNew) Reserved read capacity units of the TcaplusDB table.
* `reserved_volume` - (Required, ForceNew) Reserved storage capacity of the TcaplusDB table (unit: GB).
* `reserved_write_cu` - (Required, ForceNew) Reserved write capacity units of the TcaplusDB table.
* `table_idl_type` - (Required) IDL type of the TcaplusDB table. Valid values are PROTO and TDR.
* `table_name` - (Required, ForceNew) Name of the TcaplusDB table.
* `table_type` - (Required, ForceNew) Type of the TcaplusDB table. Valid values are GENERIC and LIST.
* `tablegroup_id` - (Required, ForceNew) ID of the table group to which the table belongs.
* `description` - (Optional) Description of the TcaplusDB table.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time of the TcaplusDB table.
* `error` - Error messages for creating TcaplusDB table.
* `status` - Status of the TcaplusDB table.
* `table_size` - Size of the TcaplusDB table.


