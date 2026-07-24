Provide a resource to create a Mongodb sharding instance.

~> **NOTE:** The `add_node_list` and `remove_node_list` arguments are used to submit node change actions. When updating the resource, only newly added items in these lists will be sent to the API. If an existing item is removed from the Terraform configuration, Terraform only updates the local state and does not submit a repeated add or remove request. To add or remove another read-only node, append a new block instead of modifying an existing one. After the change is completed, obsolete action records can be removed from the configuration, and this cleanup does not trigger a new node operation when the remaining list is a subset of the previous list. In general, it is recommended to keep these action records in the configuration and avoid cleanup unless necessary.

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

Add a read-only node

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

  add_node_list {
    role = "READONLY"
    zone = "ap-guangzhou-6"
  }
}
```

Remove a read-only node

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

  remove_node_list {
    role      = "READONLY"
    node_name = "cmgo-xxxx_0-node-readonly0"
    zone      = "ap-guangzhou-6"
  }
}
```

Import

Mongodb sharding instance can be imported using the id, e.g.

```
terraform import tencentcloud_mongodb_sharding_instance.example cmgo-41s6jwy4
```