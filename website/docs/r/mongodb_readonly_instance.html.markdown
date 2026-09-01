---
subcategory: "TencentDB for MongoDB(mongodb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mongodb_readonly_instance"
sidebar_current: "docs-tencentcloud-resource-mongodb_readonly_instance"
description: |-
  Provide a resource to create a Mongodb readonly instance.
---

# tencentcloud_mongodb_readonly_instance

Provide a resource to create a Mongodb readonly instance.

## Example Usage

### Replset readonly instance

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

### Sharding readonly instance

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

## Argument Reference

The following arguments are supported:

* `available_zone` - (Required, String) The available zone of the Mongodb.
* `cluster_type` - (Required, String, ForceNew) Instance schema type.	- REPLSET: Replset cluster;	- SHARD: Shard cluster.
* `engine_version` - (Required, String) Refers to version information. The DescribeSpecInfo API can be called to obtain detailed information about the supported versions.
- MONGO_40_WT: version of the MongoDB 4.0 WiredTiger storage engine.
- MONGO_42_WT: version of the MongoDB 4.2 WiredTiger storage engine.
- MONGO_44_WT: version of the MongoDB 4.4 WiredTiger storage engine.
- MONGO_50_WT: version of the MongoDB 5.0 WiredTiger storage engine.
- MONGO_60_WT: version of the MongoDB 6.0 WiredTiger storage engine.
- MONGO_70_WT: version of the MongoDB 7.0 WiredTiger storage engine.
- MONGO_80_WT: version of the MongoDB 8.0 WiredTiger storage engine.
* `father_instance_id` - (Required, String, ForceNew) Indicates the main instance ID of readonly instances.
* `father_instance_region` - (Required, String, ForceNew) Indicates the region of main instance.
* `instance_name` - (Required, String) Name of the Mongodb instance.
* `machine_type` - (Required, String, ForceNew) Type of Mongodb instance, and available values include `HIO`(or `GIO` which will be deprecated, represents high IO) and `HIO10G`(or `TGIO` which will be deprecated, represents 10-gigabit high IO).
* `memory` - (Required, Int) Memory size. The minimum value is 2, and unit is GB. Memory and volume must be upgraded or degraded simultaneously.
* `volume` - (Required, Int) Disk size. The minimum value is 25, and unit is GB. Memory and volume must be upgraded or degraded simultaneously.
* `auto_renew_flag` - (Optional, Int) Auto renew flag. Valid values are `0`(NOTIFY_AND_MANUAL_RENEW), `1`(NOTIFY_AND_AUTO_RENEW) and `2`(DISABLE_NOTIFY_AND_MANUAL_RENEW). Default value is `0`. Note: only works for PREPAID instance. Only supports`0` and `1` for creation.
* `charge_type` - (Optional, String, ForceNew) The charge type of instance. Valid values are `PREPAID` and `POSTPAID_BY_HOUR`. Default value is `POSTPAID_BY_HOUR`. Note: TencentCloud International only supports `POSTPAID_BY_HOUR`. Caution that update operation on this field will delete old instances and create new one with new charge type.
* `cpu` - (Optional, Int) The CPU core count of the MongoDB instance after the configuration change. Unit: C. When this parameter is empty, the current CPU size of the instance is used by default. The supported CPU specifications can be obtained through the DescribeSpecInfo API.
* `data_encryption` - (Optional, String) Database storage encryption setting. `No_Encryption`: Storage encryption is not used. `TDE`: Enables TDE storage encryption. Note: this field does not support update, please recreate the resource if you need to change it.
* `encryption_key_source` - (Optional, String) If TDE storage encryption is selected, the key source must be specified. `auto`: Automatically generate the key. `manual`: Manually specify the key. Note: this field does not support update, please recreate the resource if you need to change it.
* `in_maintenance` - (Optional, Int) Switch time for instance configuration changes.
	- 0: When the adjustment is completed, perform the configuration task immediately. Default is 0.
	- 1: Perform reconfiguration tasks within the maintenance time window.
Note: Adjusting the number of nodes and slices does not support changes within the maintenance window.
* `key_id` - (Optional, String) Key ID. If `manual` is selected as the key resource, you must enter the specified key ID. Note: this field does not support update, please recreate the resource if you need to change it.
* `kms_region` - (Optional, String) Key ID. If `manual` is selected as the key resource, you must enter the specified key region. Note: this field does not support update, please recreate the resource if you need to change it.
* `mongos_cpu` - (Optional, Int) Number of mongos cpu.
* `mongos_memory` - (Optional, Int) Mongos memory size in GB.
* `mongos_node_num` - (Optional, Int) Number of mongos.
* `node_num` - (Optional, Int) The number of nodes in each replica set. Default value: 3.
* `nodes_per_shard` - (Optional, Int, ForceNew) Number of nodes per shard, at least 3(one master and two slaves).
* `prepaid_period` - (Optional, Int) The tenancy (time unit is month) of the prepaid instance. Valid values are 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36. NOTE: it only works when charge_type is set to `PREPAID`.
* `project_id` - (Optional, Int) ID of the project which the instance belongs.
* `security_groups` - (Optional, Set: [`String`]) ID of the security group.
* `shard_quantity` - (Optional, Int, ForceNew) Number of sharding.
* `subnet_id` - (Optional, String, ForceNew) ID of the subnet within this VPC. The value is required if `vpc_id` is set.
* `tags` - (Optional, Map) The tags of the Mongodb. Key name `project` is system reserved and can't be used.
* `vpc_id` - (Optional, String, ForceNew) ID of the VPC.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time of the Mongodb instance.
* `status` - Status of the Mongodb instance, and available values include pending initialization(expressed with 0),  processing(expressed with 1), running(expressed with 2) and expired(expressed with -2).
* `vip` - IP of the Mongodb instance.
* `vport` - IP port of the Mongodb instance.


## Import

Mongodb readonly instance can be imported using the id, e.g.

```
terraform import tencentcloud_mongodb_instance.mongodb cmgo-7ohkcdu7
```

