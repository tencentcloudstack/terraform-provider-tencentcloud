Provide a resource to create a Mongodb sharding instance.

Example Usage

```hcl
resource "tencentcloud_mongodb_sharding_instance" "mongodb" {
  instance_name   = "mongodb"
  shard_quantity  = 2
  nodes_per_shard = 3
  memory          = 4
  volume          = 100
  engine_version  = "MONGO_36_WT"
  machine_type    = "HIO10G"
  available_zone  = "ap-guangzhou-3"
  vpc_id          = "vpc-mz3efvbw"
  subnet_id       = "subnet-lk0svi3p"
  project_id      = 0
  password        = "password1234"
  mongos_cpu = 1
  mongos_memory =  2
  mongos_node_num = 3
}
```

Import

Mongodb sharding instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_mongodb_sharding_instance.mongodb cmgo-41s6jwy4
```