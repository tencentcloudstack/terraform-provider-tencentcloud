Provide a resource to create a Mongodb sharding instance.

Example Usage

```hcl
resource "tencentcloud_mongodb_sharding_instance" "example" {
  instance_name   = "tf-example"
  shard_quantity  = 2
  nodes_per_shard = 3
  memory          = 4
  volume          = 100
  engine_version  = "MONGO_40_WT"
  machine_type    = "HIO10G"
  available_zone  = "ap-guangzhou-6"
  vpc_id          = "vpc-i5yyodl9"
  subnet_id       = "subnet-hhi88a58"
  project_id      = 0
  password        = "Password@123"
  mongos_cpu      = 1
  mongos_memory   = 2
  mongos_node_num = 3
}
```

Import

Mongodb sharding instance can be imported using the id, e.g.

```
terraform import tencentcloud_mongodb_sharding_instance.example cmgo-41s6jwy4
```