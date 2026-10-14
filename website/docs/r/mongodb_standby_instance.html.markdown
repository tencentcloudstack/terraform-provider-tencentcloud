---
subcategory: "TencentDB for MongoDB(mongodb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mongodb_standby_instance"
sidebar_current: "docs-tencentcloud-resource-mongodb_standby_instance"
description: |-
  Provide a resource to create a Mongodb standby instance.
---

# tencentcloud_mongodb_standby_instance

Provide a resource to create a Mongodb standby instance.

## Example Usage

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

resource "tencentcloud_mongodb_standby_instance" "example" {
  instance_name          = "tf-example"
  cpu                    = 2
  memory                 = 4
  volume                 = 100
  available_zone         = "ap-guangzhou-7"
  vpc_id                 = "vpc-i5yyodl9"
  subnet_id              = "subnet-d4umunpy"
  project_id             = 0
  father_instance_id     = tencentcloud_mongodb_instance.example.id
  father_instance_region = "ap-guangzhou"
  data_encryption        = "TDE"
  encryption_key_source  = "manual"
  key_id                 = "KMS-MONGODB"
  kms_region             = "ap-guangzhou"

  tags = {
    createBy = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `available_zone` - (Required, String, ForceNew) The available zone of the Mongodb standby instance. NOTE: must not be same with father instance's.
* `father_instance_id` - (Required, String, ForceNew) Indicates the main instance ID of standby instances.
* `father_instance_region` - (Required, String, ForceNew) Indicates the region of main instance.
* `instance_name` - (Required, String) Name of the Mongodb instance.
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
* `prepaid_period` - (Optional, Int) The tenancy (time unit is month) of the prepaid instance. Valid values are 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36. NOTE: it only works when charge_type is set to `PREPAID`.
* `project_id` - (Optional, Int) ID of the project which the instance belongs.
* `security_groups` - (Optional, Set: [`String`]) ID of the security group.
* `subnet_id` - (Optional, String, ForceNew) ID of the subnet within this VPC. The value is required if `vpc_id` is set.
* `tags` - (Optional, Map) The tags of the Mongodb. Key name `project` is system reserved and can't be used.
* `vpc_id` - (Optional, String, ForceNew) ID of the VPC.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time of the Mongodb instance.
* `engine_version` - Version of the standby Mongodb instance and must be same as the version of main instance.
* `machine_type` - Type of standby Mongodb instance and must be same as the type of main instance.
* `status` - Status of the Mongodb instance, and available values include pending initialization(expressed with 0),  processing(expressed with 1), running(expressed with 2) and expired(expressed with -2).
* `vip` - IP of the Mongodb instance.
* `vport` - IP port of the Mongodb instance.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateDBInstance](https://cloud.tencent.com/document/api/240/38571) - This API is used to create a yearly/monthly subscription TencentDB for MongoDB instance. The [DescribeSpecInfo](https://www.tencentcloud.com/document/product/240/34701) API can be called to query and obtain the supported sales specifications.
* Create: [CreateDBInstanceHour](https://cloud.tencent.com/document/api/240/38570) - This API is used to create a pay-as-you-go TencentDB for MongoDB instance.
* Create: [DescribeDBInstances](https://cloud.tencent.com/document/api/240/38568) - This API is used to query the list of TencentDB for MongoDB instances. It supports filtering primary instances, disaster recovery instances, and read-only instances by project ID, instance ID, instance status, and other conditions.
* Create: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - This API is used to modify all tags associated with a resource.
* Create: [RenameInstance](https://cloud.tencent.com/document/api/240/38563) - This API is used to rename a TencentDB instance.
* Read: [DescribeDBInstances](https://cloud.tencent.com/document/api/240/38568) - This API is used to query the list of TencentDB for MongoDB instances. It supports filtering primary instances, disaster recovery instances, and read-only instances by project ID, instance ID, instance status, and other conditions.
* Read: [DescribeSecurityGroup](https://cloud.tencent.com/document/api/240/55675) - This API is used to query security groups bound to an instance.
* Read: [DescribeTransparentDataEncryptionStatus](https://cloud.tencent.com/document/api/240/107168)
* Update: [AssignProject](https://cloud.tencent.com/document/api/240/38572) - This API is used to specify the project of a TencentDB for MongoDB instance.
* Update: [DescribeDBInstanceDeal](https://cloud.tencent.com/document/api/240/43667) - This API is used to get order details of purchase, renewal, and specification adjustment of a MongoDB instance.
* Update: [ModifyDBInstanceSecurityGroup](https://cloud.tencent.com/document/api/240/70577) - This API is used to modify security groups bound to an instance.
* Update: [ModifyDBInstanceSpec](https://cloud.tencent.com/document/api/240/38565) - This API is used to adjust the TencentDB for MongoDB instance configuration. The [DescribeSpecInfo](https://www.tencentcloud.com/document/product/240/38567?from_cn_redirect=1) API can be called to query and obtain the supported sales specifications.
* Update: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - This API is used to modify all tags associated with a resource.
* Update: [RenameInstance](https://cloud.tencent.com/document/api/240/38563) - This API is used to rename a TencentDB instance.
* Update: [RenewDBInstances](https://cloud.tencent.com/document/api/240/43595) - This API is used to renew a monthly subscription TencentDB instance. Only monthly subscription instances are supported, while pay-as-you-go instances do not need to be renewed.
* Delete: [DescribeDBInstances](https://cloud.tencent.com/document/api/240/38568) - This API is used to query the list of TencentDB for MongoDB instances. It supports filtering primary instances, disaster recovery instances, and read-only instances by project ID, instance ID, instance status, and other conditions.
* Delete: [IsolateDBInstance](https://cloud.tencent.com/document/api/240/38566) - This API is used to isolate a pay-as-you-go TencentDB for MongoDB instance. After isolation, the instance is retained in the recycle bin, and data cannot be written into it. After a certain period of isolation, the instance is deleted permanently. For the retention time in the recycle bin, see the pay-as-you-go service terms. The deleted pay-as-you-go instance cannot be recovered. Proceed with caution.
* Delete: [OfflineIsolatedDBInstance](https://cloud.tencent.com/document/api/240/38564) - This API is used to deactivate isolated TencentDB instances immediately. The instances must be in isolated status.

## Import

Mongodb standby instance can be imported using the id, e.g.

```
terraform import tencentcloud_mongodb_standby_instance.mongodb cmgo-41s6jwy4
```

