Provide a resource to create a Readonly mongodb instance.

Example Usage

Replset readonly instance

```hcl
resource "tencentcloud_mongodb_readonly_instance" "mongodb" {
  instance_name          = "tf-mongodb-readonly-test"
  memory                 = 4
  volume                 = 100
  engine_version         = "MONGO_44_WT"
  machine_type           = "HIO10G"
  available_zone         = "ap-guangzhou-3"
  project_id             = 0
  father_instance_id     = "cmgo-xxxxxx"
  father_instance_region = "ap-guangzhou"
  vpc_id                 = "vpc-xxxxxx"
  subnet_id              = "subnet-xxxxxx"
  security_groups        = ["sg-xxxxxx"]
  cluster_type           = "REPLSET"
}
```

Shard readonly instance

```hcl
resource "tencentcloud_mongodb_readonly_instance" "sharding_mongodb" {
  instance_name          = "tf-mongodb-readonly-shard"
  memory                 = 4
  volume                 = 100
  engine_version         = "MONGO_44_WT"
  machine_type           = "HIO10G"
  available_zone         = "ap-guangzhou-3"
  project_id             = 0
  father_instance_id     = "cmgo-xxxxxx"
  father_instance_region = "ap-guangzhou"
  vpc_id                 = "vpc-xxxxxx"
  subnet_id              = "subnet-xxxxxx"
  security_groups        = ["sg-xxxxxx"]
  cluster_type           = "SHARD"
  mongos_cpu             = 1
  mongos_memory          = 2
  mongos_node_num        = 3
}
```

Import

Mongodb instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_mongodb_instance.mongodb cmgo-xxxxxx
```