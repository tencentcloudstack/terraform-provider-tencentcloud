resource "tencentcloud_route_table" "rtb_test" {
  name   = "rtb-test"
  vpc_id = "${tencentcloud_vpc.vpc_test.id}"
}

resource "tencentcloud_subnet" "subnet_test" {
  name              = "subnet-test"
  cidr_block        = "10.0.1.0/24"
  availability_zone = "ap-guangzhou-6"
  vpc_id            = "${tencentcloud_vpc.vpc_test.id}"
  route_table_id    = "${tencentcloud_route_table.rtb_test.id}"
}

resource "tencentcloud_vpc" "vpc_test" {
  name       = "vpc-test"
  cidr_block = "10.0.0.0/16"
}

# create ckafka instance
resource "tencentcloud_ckafka_instance" "foo" {
  instance_name      = "tf-test"
  zone_id            = 100006
  period             = 1
  vpc_id             = "${tencentcloud_vpc.vpc_test.id}"
  subnet_id          = "${tencentcloud_subnet.subnet_test.id}"
  msg_retention_time = 1300
  renew_flag         = 0
  kafka_version      = "1.1.1"
  disk_size          = 500
  disk_type          = "CLOUD_BASIC"


  config {
    auto_create_topic_enable   = true
    default_num_partitions     = 3
    default_replication_factor = 3
  }

  dynamic_retention_config {
    enable = 1
  }
}