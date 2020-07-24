provider "tencentcloud" {
  region = "ap-guangzhou"
}

provider "tencentcloud" {
  alias  = "shanghai"
  region = "ap-shanghai"
}

resource "tencentcloud_mongodb_instance" "mongodb" {
  instance_name  = "tf-mongodb-test"
  memory         = 4
  volume         = 100
  engine_version = "MONGO_40_WT"
  machine_type   = "HIO10G"
  available_zone = var.availability_zone
  project_id     = 0
  password       = "test1234"

  tags = {
    test = "test"
  }
}

resource "tencentcloud_mongodb_sharding_instance" "mongodb_sharding" {
  instance_name   = "tf-mongodb-sharding"
  shard_quantity  = 2
  nodes_per_shard = 3
  memory          = 4
  volume          = 100
  engine_version  = "MONGO_36_WT"
  machine_type    = "TGIO"
  available_zone  = var.availability_zone
  project_id      = 0
  password        = "test1234"

  tags = {
    test = "test"
  }
}

data "tencentcloud_mongodb_instances" "mongodb_instances" {
  instance_id = tencentcloud_mongodb_instance.mongodb.id

  tags = {
    test = "test"
  }
}

resource "tencentcloud_mongodb_standby_instance" "mongodb" {
  provider               = tencentcloud.shanghai
  instance_name          = "tf-mongodb-standby-test"
  memory                 = 4
  volume                 = 100
  available_zone         = "ap-shanghai-2"
  project_id             = 0
  father_instance_id     = tencentcloud_mongodb_instance.mongodb.id
  father_instance_region = "ap-guangzhou"
  charge_type            = "PREPAID"
  prepaid_period         = 1

  tags = {
    test = "test"
  }
}