Use this data source to query Publish Subscribe resources for the specific SQL Server instance.

Example Usage

```hcl
data "tencentcloud_sqlserver_publish_subscribes" "example" {
  instance_id = tencentcloud_sqlserver_publish_subscribe.example.publish_instance_id
}

data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_security_group" "security_group" {
  name        = "sg-example"
  description = "desc."
}

resource "tencentcloud_sqlserver_general_cloud_instance" "example_pub" {
  name                 = "tf-example-pub"
  zone                 = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  memory               = 4
  storage              = 100
  cpu                  = 2
  machine_type         = "CLOUD_HSSD"
  instance_charge_type = "POSTPAID"
  project_id           = 0
  subnet_id            = tencentcloud_subnet.subnet.id
  vpc_id               = tencentcloud_vpc.vpc.id
  db_version           = "2008R2"
  security_group_list  = [tencentcloud_security_group.security_group.id]
  weekly               = [1, 2, 3, 5, 6, 7]
  start_time           = "00:00"
  span                 = 6
  resource_tags {
    tag_key   = "test"
    tag_value = "test"
  }
  collation = "Chinese_PRC_CI_AS"
  time_zone = "China Standard Time"
}

resource "tencentcloud_sqlserver_general_cloud_instance" "example_sub" {
  name                 = "tf-example-sub"
  zone                 = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  memory               = 4
  storage              = 100
  cpu                  = 2
  machine_type         = "CLOUD_HSSD"
  instance_charge_type = "POSTPAID"
  project_id           = 0
  subnet_id            = tencentcloud_subnet.subnet.id
  vpc_id               = tencentcloud_vpc.vpc.id
  db_version           = "2008R2"
  security_group_list  = [tencentcloud_security_group.security_group.id]
  weekly               = [1, 2, 3, 5, 6, 7]
  start_time           = "00:00"
  span                 = 6
  resource_tags {
    tag_key   = "test"
    tag_value = "test"
  }
  collation = "Chinese_PRC_CI_AS"
  time_zone = "China Standard Time"
}

resource "tencentcloud_sqlserver_db" "example_pub" {
  instance_id = tencentcloud_sqlserver_general_cloud_instance.example_pub.id
  name        = "tf_example_db_pub"
  charset     = "Chinese_PRC_BIN"
  remark      = "test-remark"
}

resource "tencentcloud_sqlserver_db" "example_sub" {
  instance_id = tencentcloud_sqlserver_general_cloud_instance.example_sub.id
  name        = "tf_example_db_sub"
  charset     = "Chinese_PRC_BIN"
  remark      = "test-remark"
}

resource "tencentcloud_sqlserver_publish_subscribe" "example" {
  publish_instance_id    = tencentcloud_sqlserver_general_cloud_instance.example_pub.id
  subscribe_instance_id  = tencentcloud_sqlserver_general_cloud_instance.example_sub.id
  publish_subscribe_name = "example"
  delete_subscribe_db    = false
  database_tuples {
    publish_database   = tencentcloud_sqlserver_db.example_pub.name
    subscribe_database = tencentcloud_sqlserver_db.example_sub.name
  }
}
```