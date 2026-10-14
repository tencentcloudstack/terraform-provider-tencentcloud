---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_clone_db_instance"
sidebar_current: "docs-tencentcloud-resource-postgresql_clone_db_instance"
description: |-
  Provides a resource to create a postgresql clone db instance
---

# tencentcloud_postgresql_clone_db_instance

Provides a resource to create a postgresql clone db instance

## Example Usage

### Clone db instance by recovery_target_time

```hcl
resource "tencentcloud_postgresql_clone_db_instance" "example" {
  db_instance_id       = "postgres-ckwcgdf1"
  name                 = "tf-example"
  spec_code            = "pg.it.medium4"
  storage              = 100
  period               = 1
  auto_renew_flag      = 0
  vpc_id               = "vpc-i5yyodl9"
  subnet_id            = "subnet-hhi88a58"
  instance_charge_type = "POSTPAID_BY_HOUR"
  security_group_ids   = ["sg-rs32zv1r", "sg-37tigqat"]
  project_id           = 0
  recovery_target_time = "2026-07-10 01:00:06"
  deletion_protection  = true
  db_node_set {
    role = "Primary"
    zone = "ap-guangzhou-6"
  }

  db_node_set {
    role = "Standby"
    zone = "ap-guangzhou-7"
  }

  tags = {
    tagKey = "tagValue"
  }
}
```

### Clone db instance by backup_set_id

```hcl
data "tencentcloud_postgresql_base_backups" "base_backups" {
  filters {
    name   = "db-instance-id"
    values = ["postgres-evsqpyap"]
  }

  order_by      = "Size"
  order_by_type = "asc"
}

resource "tencentcloud_postgresql_clone_db_instance" "example" {
  db_instance_id       = "postgres-evsqpyap"
  name                 = "tf-example-clone"
  spec_code            = "pg.it.medium4"
  storage              = 200
  period               = 1
  auto_renew_flag      = 0
  vpc_id               = "vpc-a6zec4mf"
  subnet_id            = "subnet-b8hintyy"
  instance_charge_type = "POSTPAID_BY_HOUR"
  security_group_ids   = ["sg-8stavs03"]
  project_id           = 0
  backup_set_id        = data.tencentcloud_postgresql_base_backups.base_backups.base_backup_set.0.id
  deletion_protection  = true
  db_node_set {
    role = "Primary"
    zone = "ap-guangzhou-6"
  }

  db_node_set {
    role = "Standby"
    zone = "ap-guangzhou-6"
  }

  tags = {
    tagKey = "tagValue"
  }
}
```

### Clone db instance from CDC

```hcl
resource "tencentcloud_postgresql_clone_db_instance" "example" {
  db_instance_id       = "postgres-evsqpyap"
  name                 = "tf-example-clone"
  spec_code            = "pg.it.medium4"
  storage              = 200
  period               = 1
  auto_renew_flag      = 0
  vpc_id               = "vpc-a6zec4mf"
  subnet_id            = "subnet-b8hintyy"
  instance_charge_type = "POSTPAID_BY_HOUR"
  security_group_ids   = ["sg-8stavs03"]
  project_id           = 0
  recovery_target_time = "2024-10-12 18:17:00"
  deletion_protection  = true
  db_node_set {
    role                 = "Primary"
    zone                 = "ap-guangzhou-6"
    dedicated_cluster_id = "cluster-262n63e8"
  }

  db_node_set {
    role                 = "Standby"
    zone                 = "ap-guangzhou-6"
    dedicated_cluster_id = "cluster-262n63e8"
  }

  tags = {
    tagKey = "tagValue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `auto_renew_flag` - (Required, Int) Renewal Flag:

- `0`: manual renewal
`1`: auto-renewal

Default value: 0.
* `db_instance_id` - (Required, String, ForceNew) ID of the original instance to be cloned.
* `period` - (Required, Int) Purchase duration, in months.
- Prepaid: Supports `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, and `36`.
- Pay-as-you-go: Only supports `1`.
* `spec_code` - (Required, String, ForceNew) Purchasable code, which can be obtained from the `SpecCode` field in the return value of the [DescribeClasses](https://intl.cloud.tencent.com/document/api/409/89019?from_cn_redirect=1) API.
* `storage` - (Required, Int, ForceNew) Instance storage capacity in GB.
* `subnet_id` - (Required, String) VPC subnet ID in the format of `subnet-xxxxxxxx`, which can be obtained in the console or from the `unSubnetId` field in the return value of the [DescribeSubnets](https://intl.cloud.tencent.com/document/api/215/15784?from_cn_redirect=1) API.
* `vpc_id` - (Required, String) VPC ID in the format of `vpc-xxxxxxx`, which can be obtained in the console or from the `unVpcId` field in the return value of the [DescribeVpcEx](https://intl.cloud.tencent.com/document/api/215/1372?from_cn_redirect=1) API.
* `activity_id` - (Optional, Int, ForceNew) Campaign ID.
* `backup_set_id` - (Optional, String, ForceNew) Basic backup set ID.
* `db_node_set` - (Optional, List) Deployment information of the instance node, which will display the information of each AZ when the instance node is deployed across multiple AZs.
The information of AZ can be obtained from the `Zone` field in the return value of the [DescribeZones](https://intl.cloud.tencent.com/document/api/409/16769?from_cn_redirect=1) API.
* `deletion_protection` - (Optional, Bool) Whether deletion protection is enabled for the instance: `true` deletion protection enabled; `false` deletion protection disabled.
* `instance_charge_type` - (Optional, String) Instance billing type, which currently supports:

- PREPAID: Prepaid, i.e., monthly subscription
- POSTPAID_BY_HOUR: Pay-as-you-go, i.e., pay by consumption

Default value: PREPAID.
* `name` - (Optional, String) Name of the newly purchased instance, which can contain up to 60 letters, digits, or symbols (-_). If this parameter is not specified, "Unnamed" will be displayed by default.
* `project_id` - (Optional, Int) Project ID.
* `recovery_target_time` - (Optional, String, ForceNew) Restoration point in time.
* `security_group_ids` - (Optional, Set: [`String`]) Security group of the instance, which can be obtained from the `sgld` field in the return value of the [DescribeSecurityGroups](https://intl.cloud.tencent.com/document/api/215/15808?from_cn_redirect=1) API. If this parameter is not specified, the default security group will be bound.
* `sync_mode` - (Optional, String, ForceNew) Primary-standby sync mode, which supports:
Semi-sync: Semi-sync
Async: Asynchronous
Default value for the primary instance: Semi-sync
Default value for the read-only instance: Async.
* `tag_list` - (Optional, List, ForceNew, **Deprecated**) It has been deprecated from version 1.83.10. Use `tags` instead. The information of tags to be bound with the instance, which is left empty by default. This parameter can be obtained from the `Tags` field in the return value of the [DescribeTags](https://intl.cloud.tencent.com/document/api/651/35316?from_cn_redirect=1) API.
* `tags` - (Optional, Map) The available tags within this postgresql.

The `db_node_set` object supports the following:

* `role` - (Required, String) Node type. Valid values:
`Primary`;
`Standby`.
* `zone` - (Required, String) AZ where the node resides, such as ap-guangzhou-1.
* `dedicated_cluster_id` - (Optional, String) Dedicated cluster ID.

The `tag_list` object supports the following:

* `tag_key` - (Required, String, ForceNew) Tag key.
* `tag_value` - (Required, String, ForceNew) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `availability_zone` - Availability zone.
* `new_db_instance_id` - ID of the cloned instance.
* `private_access_ip` - IP for private access.
* `private_access_port` - Port for private access.
* `public_access_host` - Host for public access.
* `public_access_port` - Port for public access.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `1h0m`) Used when creating the resource.
* `update` - (Defaults to `1h0m`) Used when updating the resource.
* `delete` - (Defaults to `1h0m`) Used when deleting the resource.
## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CloneDBInstance](https://cloud.tencent.com/document/api/409/68071) - This API is used to clone an instance by specifying a backup set or a point in time.
* Create: [DescribeDBInstanceAttribute](https://cloud.tencent.com/document/api/409/16772) - This API is used to query the details of one instance.
* Create: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - This API is used to modify all tags associated with a resource.
* Read: [DescribeDBInstanceAttribute](https://cloud.tencent.com/document/api/409/16772) - This API is used to query the details of one instance.
* Read: [DescribeDBInstanceSecurityGroups](https://cloud.tencent.com/document/api/409/84064) - This API is used to query the security group of an instance.
* Read: [DescribeResourceTagsByResourceIds](https://cloud.tencent.com/document/api/651/35319) - This API is used to query the tag key-value pairs associated with an existing resource.
* Update: [CreateDBInstanceNetworkAccess](https://cloud.tencent.com/document/api/409/70581) - This API is used to create instance networks. Each instance allows a maximum of 2 network configurations and a minimum of 1.
* Update: [DeleteDBInstanceNetworkAccess](https://cloud.tencent.com/document/api/409/70579) - This API is used to delete a network from an instance while keeping at least one network.
* Update: [DescribeDBInstanceAttribute](https://cloud.tencent.com/document/api/409/16772) - This API is used to query the details of one instance.
* Update: [ModifyDBInstanceChargeType](https://cloud.tencent.com/document/api/409/90458) - This API is used to switch the instance billing mode from pay-as-you-go to yearly/monthly subscription.
* Update: [ModifyDBInstanceDeletionProtection](https://cloud.tencent.com/document/api/409/125336) - This interface (DeletionProtection) is used to enable or disable instance destruction protection.
* Update: [ModifyDBInstanceDeployment](https://cloud.tencent.com/document/api/409/68066) - This API is used to modify the AZs where the nodes of a source instance reside.
* Update: [ModifyDBInstanceName](https://cloud.tencent.com/document/api/409/18101) - This API is used to rename a TencentDB for PostgreSQL instance.
* Update: [ModifyDBInstanceSecurityGroups](https://cloud.tencent.com/document/api/409/84063) - This API is used to modify the security group of an instance.
* Update: [ModifyDBInstancesProject](https://cloud.tencent.com/document/api/409/18100) - This API is used to modify the project of an instance.
* Update: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - This API is used to modify all tags associated with a resource.
* Update: [RenewInstance](https://cloud.tencent.com/document/api/409/18098) - This API is used to renew an instance.
* Update: [SetAutoRenewFlag](https://cloud.tencent.com/document/api/409/18096) - This API is used to set auto-renewal.
* Delete: [DescribeDBInstanceAttribute](https://cloud.tencent.com/document/api/409/16772) - This API is used to query the details of one instance.
* Delete: [DestroyDBInstance](https://cloud.tencent.com/document/api/409/43352) - This API is used to terminate an isolated instance by specifying the `DBInstanceId` parameter. The data of a terminated instance will be deleted and cannot be recovered. Be cautious with this API calling.
* Delete: [IsolateDBInstances](https://cloud.tencent.com/document/api/409/54790) - This API is used to isolate an instance.

