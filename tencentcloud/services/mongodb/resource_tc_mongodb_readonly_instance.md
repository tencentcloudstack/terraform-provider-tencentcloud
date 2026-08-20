Provide a resource to create a Mongodb readonly instance.

Example Usage

Replset readonly instance

```hcl
resource "tencentcloud_mongodb_instance" "example" {
  instance_name         = "tf-example"
  cpu                   = 2
  memory                = 4
  volume                = 100
  engine_version        = "MONGO_80_WT"
  machine_type          = "GE.LD.T1"
  available_zone        = "ap-guangzhou-6"
  vpc_id                = "vpc-i5yyodl9"
  subnet_id             = "subnet-hhi88a58"
  project_id            = 0
  password              = "Password@2026"
  data_encryption       = "TDE"
  encryption_key_source = "manual"
  key_id                = "KMS-MONGODB"
  kms_region            = "ap-guangzhou"
}

resource "tencentcloud_mongodb_readonly_instance" "example" {
  instance_name          = "tf-mongodb"
  cpu                    = 2
  memory                 = 4
  volume                 = 100
  engine_version         = "MONGO_80_WT"
  machine_type           = "GE.LD.T1"
  available_zone         = "ap-guangzhou-7"
  project_id             = 0
  father_instance_id     = tencentcloud_mongodb_instance.example.id
  father_instance_region = "ap-guangzhou"
  vpc_id                 = "vpc-i5yyodl9"
  subnet_id              = "subnet-d4umunpy"
  security_groups        = ["sg-n8zf5ry9"]
  cluster_type           = "REPLSET"
  data_encryption        = "No_Encryption"
}
```

Sharding readonly instance

```hcl
resource "tencentcloud_mongodb_sharding_instance" "example" {
  instance_name         = "tf-example"
  shard_quantity        = 2
  nodes_per_shard       = 3
  cpu                   = 2
  memory                = 4
  volume                = 100
  engine_version        = "MONGO_80_WT"
  machine_type          = "GE.LD.T1"
  available_zone        = "ap-guangzhou-6"
  vpc_id                = "vpc-i5yyodl9"
  subnet_id             = "subnet-hhi88a58"
  project_id            = 0
  password              = "Password@2026"
  mongos_cpu            = 1
  mongos_memory         = 2
  mongos_node_num       = 3
  data_encryption       = "TDE"
  encryption_key_source = "auto"
}

resource "tencentcloud_mongodb_readonly_instance" "example" {
  instance_name          = "tf-mongodb"
  cpu                    = 2
  memory                 = 4
  volume                 = 100
  engine_version         = "MONGO_80_WT"
  machine_type           = "GE.LD.T1"
  available_zone         = "ap-guangzhou-7"
  project_id             = 0
  father_instance_id     = tencentcloud_mongodb_sharding_instance.example.id
  father_instance_region = "ap-guangzhou"
  vpc_id                 = "vpc-i5yyodl9"
  subnet_id              = "subnet-d4umunpy"
  security_groups        = ["sg-n8zf5ry9"]
  cluster_type           = "SHARD"
  mongos_cpu             = 1
  mongos_memory          = 2
  mongos_node_num        = 3
  data_encryption        = "TDE"
  encryption_key_source  = "manual"
  key_id                 = "KMS-MONGODB"
  kms_region             = "ap-guangzhou"
}
```

Import

Mongodb readonly instance can be imported using the id, e.g.

```
terraform import tencentcloud_mongodb_instance.mongodb cmgo-7ohkcdu7
```
