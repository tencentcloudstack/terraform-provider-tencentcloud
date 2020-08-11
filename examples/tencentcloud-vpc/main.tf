resource "tencentcloud_vpc" "main" {
  name       = var.short_name
  cidr_block = var.vpc_cidr

  tags = {
    "test" = "test"
  }
}

data "tencentcloud_vpc_instances" "tags_instances" {
  name = tencentcloud_vpc.main.name
  tags = tencentcloud_vpc.main.tags
}

data "tencentcloud_vpc_instances" "default" {}

resource "tencentcloud_vpc_acl" "foo" {
  vpc_id  = data.tencentcloud_vpc_instances.default.instance_list.0.vpc_id
  name    = "test_acl_update"
  ingress = ["ACCEPT#192.168.1.0/24#800#TCP", "ACCEPT#192.168.1.0/24#800-900#TCP",]
  egress  = ["ACCEPT#192.168.1.0/24#800#TCP", "ACCEPT#192.168.1.0/24#800-900#TCP",]
}

resource "tencentcloud_vpc_acl_attachment" "attachment" {
  acl_id     = tencentcloud_vpc_acl.foo.id
  subnet_ids = data.tencentcloud_vpc_instances.default.instance_list[0].subnet_ids
}

data "tencentcloud_vpc_acls" "foo" {
  name = "test_acl"
}