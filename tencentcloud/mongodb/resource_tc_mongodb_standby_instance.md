Provide a resource to create a Mongodb standby instance.

Example Usage

```hcl
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

resource "tencentcloud_mongodb_standby_instance" "mongodb" {
  provider               = tencentcloud.shanghai
  instance_name          = "tf-mongodb-standby-test"
  memory                 = 4
  volume                 = 100
  available_zone         = "ap-shanghai-2"
  project_id             = 0
  father_instance_id     = tencentcloud_mongodb_instance.mongodb.id
  father_instance_region = "ap-guangzhou"

  tags = {
    test = "test"
  }
}
```

Import

Mongodb instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_mongodb_standby_instance.mongodb cmgo-41s6jwy4
```