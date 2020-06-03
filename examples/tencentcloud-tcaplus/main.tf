resource "tencentcloud_vpc" "foo" {
  name       = var.instance_name
  cidr_block = var.vpc_cidr
}

resource "tencentcloud_subnet" "subnet" {
  name              = var.instance_name
  vpc_id            = tencentcloud_vpc.foo.id
  availability_zone = var.availability_zone
  cidr_block        = var.subnet_cidr
  is_multicast      = false
}

resource "tencentcloud_tcaplus_cluster" "test_cluster" {
  idl_type                 = "PROTO"
  cluster_name             = "tf_tcaplus_g_table"
  vpc_id                   = tencentcloud_vpc.foo.id
  subnet_id                = tencentcloud_subnet.subnet.id
  password                 = "1qaA2k1wgvfa3ZZZ"
  old_password_expire_last = 3600
}

resource "tencentcloud_tcaplus_idl" "test_idl" {
  cluster_id    = tencentcloud_tcaplus_cluster.test_cluster.id
  tablegroup_id = tencentcloud_tcaplus_tablegroup.test_tablegroup.id
  file_name     = "tf_idl_test_guagua"
  file_type     = "PROTO"
  file_ext_type = "proto"
  file_content  = <<EOF
    syntax = "proto2";
    package myTcaplusTable;
    import "tcaplusservice.optionv1.proto";
    message tb_online_guagua {
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

resource "tencentcloud_tcaplus_tablegroup" "test_tablegroup" {
  cluster_id      = tencentcloud_tcaplus_cluster.test_cluster.id
  tablegroup_name = "tf_test_tablegroup"
}

resource "tencentcloud_tcaplus_table" "test_table" {
  cluster_id         = tencentcloud_tcaplus_cluster.test_cluster.id
  tablegroup_id      = tencentcloud_tcaplus_tablegroup.test_tablegroup.id
  table_name         = "tb_online_guagua"
  table_type         = "GENERIC"
  description        = "test"
  idl_id             = tencentcloud_tcaplus_idl.test_idl.id
  table_idl_type     = "PROTO"
  reserved_read_qps  = 1000
  reserved_write_qps = 20
  reserved_volume    = 1
}