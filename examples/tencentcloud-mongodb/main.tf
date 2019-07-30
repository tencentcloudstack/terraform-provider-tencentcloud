resource "tencentcloud_mongodb_instance" "mongodb" {
  instance_name  = "tf-mongodb-test"
  memory         = 4
  volume         = 100
  engine_version = "MONGO_36_WT"
  machine_type   = "GIO"
  available_zone = "${var.availability_zone}"
  project_id     = 0
  password       = "test1234"
}

resource "tencentcloud_mongodb_sharding_instance" "mongodb_sharding" {
  instance_name   = "tf-mongodb-sharding"
  shard_quantity  = 2
  nodes_per_shard = 3
  memory          = 4
  volume          = 100
  engine_version  = "MONGO_36_WT"
  machine_type    = "TGIO"
  available_zone  = "${var.availability_zone}"
  project_id      = 0
  password        = "test1234"
}

data "tencentcloud_mongodb_instances" "mongodb_instances" {
  instance_id = "${tencentcloud_mongodb_instance.mongodb.id}"
}
