data "tencentcloud_vpc_subnets" "gz3" {
  availability_zone = "ap-guangzhou-3"
  is_default = true
}

locals {
  vpc_id = data.tencentcloud_vpc_subnets.gz3.instance_list.0.vpc_id
  subnet_id = data.tencentcloud_vpc_subnets.gz3.instance_list.0.subnet_id
}

resource "tencentcloud_tcaplus_cluster" "foo" {
  idl_type                 = "TDR"
  cluster_name             = "keep_tdr_tcaplus_cluster"
  vpc_id                   = local.vpc_id
  subnet_id                = local.subnet_id
  password                 = "1qaA2k1wgvfa3ZZZ"
  old_password_expire_last = 3600
}

resource "tencentcloud_tcaplus_tablegroup" "tablegroup" {
  cluster_id      = tencentcloud_tcaplus_cluster.foo.id
  tablegroup_name = "keep_tdr_table_group"
}