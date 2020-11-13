data "tencentcloud_sqlserver_zone_config" "foo" {
}

resource "tencentcloud_vpc" "foo" {
  name       = "example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "foo" {
  name              = "example"
  availability_zone = var.availability_zone
  vpc_id            = tencentcloud_vpc.foo.id
  cidr_block        = "10.0.0.0/24"
  is_multicast      = false
}

resource "tencentcloud_sqlserver_instance" "example" {
  name              = "example"
  availability_zone = var.availability_zone
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = tencentcloud_vpc.foo.id
  subnet_id         = tencentcloud_subnet.foo.id
  engine_version    = "2008R2"
  project_id        = 0
  memory            = 2
  storage           = 10
  tags = {
      "test" = "test"
  }
}

resource "tencentcloud_sqlserver_db" "example" {
  instance_id = tencentcloud_sqlserver_instance.example.id
  name        = "example"
  charset     = "Chinese_PRC_BIN"
  remark      = "tf"
}

resource "tencentcloud_sqlserver_account" "example" {
  instance_id = tencentcloud_sqlserver_instance.example.id
  name        = "example"
  password    = "test1233"
  remark      = "tf"
}

resource "tencentcloud_sqlserver_account_db_attachment" "example" {
  instance_id  = tencentcloud_sqlserver_instance.example.id
  account_name = tencentcloud_sqlserver_account.example.name
  db_name      = tencentcloud_sqlserver_db.example.name
  privilege    = "ReadWrite"
}

resource "tencentcloud_sqlserver_readonly_instance" "example" {
  name                = "example"
  availability_zone   = var.availability_zone
  charge_type         = "POSTPAID_BY_HOUR"
  vpc_id              = tencentcloud_vpc.foo.id
  subnet_id           = tencentcloud_subnet.foo.id
  memory              = 4
  storage             = 20
  master_instance_id  = tencentcloud_sqlserver_instance.test.id
  readonly_group_type = 1
  force_upgrade       = true
  tags = {
      "test" = "test"
  }
}

resource "tencentcloud_sqlserver_publish_subscribe" "example" {
	publish_instance_id             = tencentcloud_sqlserver_instance.example.id
	subscribe_instance_id           = tencentcloud_sqlserver_instance.example_other.id
	publish_subscribe_name          = "example"
	database_tuples {
		publish_database            = tencentcloud_sqlserver_db.example.name
	}
}


data "tencentcloud_sqlserver_instances" "id_example" {
  id = tencentcloud_sqlserver_instance.example.id
}

data "tencentcloud_sqlserver_instances" "vpc_example" {
  vpc_id    = tencentcloud_vpc.foo.id
  subnet_id = tencentcloud_subnet.foo.id
}

data "tencentcloud_sqlserver_instances" "project_example" {
  project_id = 0
}

data "tencentcloud_sqlserver_dbs" "example" {
  instance_id = tencentcloud_sqlserver_db.example.instance_id
}

data "tencentcloud_sqlserver_accounts" "example" {
  instance_id = tencentcloud_sqlserver_instance.example.id
}

data "tencentcloud_sqlserver_account_db_attachments" "example" {
  instance_id  = tencentcloud_sqlserver_instance.example.id
  account_name = tencentcloud_sqlserver_account.example.name
}

data "tencentcloud_sqlserver_backups" "example" {
  instance_id = tencentcloud_sqlserver_instance.example.id
  start_time  = "2020-06-30 00:00:00"
  end_time    = "2020-07-01 00:00:00"
}

data "tencentcloud_sqlserver_readonly_groups" "example" {
  master_instance_id = tencentcloud_sqlserver_instance.example.id
}

data "tencentcloud_sqlserver_publish_subscribes" "publish_subscribes" {
	instance_id                     = tencentcloud_sqlserver_publish_subscribe.example.publish_instance_id
	pub_or_sub_instance_id          = tencentcloud_sqlserver_publish_subscribe.example.subscribe_instance_id
	publish_subscribe_name          = tencentcloud_sqlserver_publish_subscribe.example.publish_subscribe_name
}

