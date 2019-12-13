---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcaplus_table"
sidebar_current: "docs-tencentcloud-resource-tcaplus_table"
description: |-
  Use this resource to create tcaplus table
---

# tencentcloud_tcaplus_table

Use this resource to create tcaplus table

## Example Usage

```hcl
resource "tencentcloud_tcaplus_application" "test" {
  idl_type                 = "PROTO"
  app_name                 = "tf_tcaplus_app_test"
  vpc_id                   = "vpc-7k6gzox6"
  subnet_id                = "subnet-akwgvfa3"
  password                 = "1qaA2k1wgvfa3ZZZ"
  old_password_expire_last = 3600
}

resource "tencentcloud_tcaplus_zone" "zone" {
  app_id    = tencentcloud_tcaplus_application.test.id
  zone_name = "tf_test_zone_name"
}

resource "tencentcloud_tcaplus_idl" "main" {
  app_id        = tencentcloud_tcaplus_application.test.id
  zone_id       = tencentcloud_tcaplus_zone.zone.id
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
  app_id             = tencentcloud_tcaplus_application.test.id
  zone_id            = tencentcloud_tcaplus_zone.zone.id
  table_name         = "tb_online"
  table_type         = "GENERIC"
  description        = "test"
  idl_id             = tencentcloud_tcaplus_idl.main.id
  table_idl_type     = "PROTO"
  reserved_read_qps  = 1000
  reserved_write_qps = 20
  reserved_volume    = 1
}
```

## Argument Reference

The following arguments are supported:

* `app_id` - (Required, ForceNew) Application of this table belongs.
* `idl_id` - (Required) Idl file for this table.
* `reserved_read_qps` - (Required, ForceNew) Table reserved read QPS.
* `reserved_volume` - (Required, ForceNew) Table reserved capacity(GB).
* `reserved_write_qps` - (Required, ForceNew) Table reserved write QPS.
* `table_idl_type` - (Required) Type of this table idl, Valid values are PROTO,TDR.
* `table_name` - (Required, ForceNew) Name of this table.
* `table_type` - (Required, ForceNew) Type of this table, Valid values are GENERIC,LIST.
* `zone_id` - (Required, ForceNew) Zone of this table belongs.
* `description` - (Optional) Description of this table.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - Create time of the tcapplus table.
* `error` - Show if this table  create error.
* `status` - Status of this table.
* `table_size` - Size of this table.


