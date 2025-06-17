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
  db_node_set {
    role = "Primary"
    zone = "ap-guangzhou-6"
  }

  db_node_set {
    role = "Standby"
    zone = "ap-guangzhou-6"
  }

  tag_list {
    tag_key   = "createBy"
    tag_value = "Terraform"
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
  db_node_set {
    role = "Primary"
    zone = "ap-guangzhou-6"
  }

  db_node_set {
    role = "Standby"
    zone = "ap-guangzhou-6"
  }

  tag_list {
    tag_key   = "createBy"
    tag_value = "Terraform"
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

  tag_list {
    tag_key   = "createBy"
    tag_value = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `auto_renew_flag` - (Required, Int, ForceNew) Renewal Flag:

- `0`: manual renewal
`1`: auto-renewal

Default value: 0.
* `db_instance_id` - (Required, String, ForceNew) ID of the original instance to be cloned.
* `period` - (Required, Int, ForceNew) Purchase duration, in months.
- Prepaid: Supports `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, and `36`.
- Pay-as-you-go: Only supports `1`.
* `spec_code` - (Required, String, ForceNew) Purchasable code, which can be obtained from the `SpecCode` field in the return value of the [DescribeClasses](https://intl.cloud.tencent.com/document/api/409/89019?from_cn_redirect=1) API.
* `storage` - (Required, Int, ForceNew) Instance storage capacity in GB.
* `subnet_id` - (Required, String, ForceNew) VPC subnet ID in the format of `subnet-xxxxxxxx`, which can be obtained in the console or from the `unSubnetId` field in the return value of the [DescribeSubnets](https://intl.cloud.tencent.com/document/api/215/15784?from_cn_redirect=1) API.
* `vpc_id` - (Required, String, ForceNew) VPC ID in the format of `vpc-xxxxxxx`, which can be obtained in the console or from the `unVpcId` field in the return value of the [DescribeVpcEx](https://intl.cloud.tencent.com/document/api/215/1372?from_cn_redirect=1) API.
* `activity_id` - (Optional, Int, ForceNew) Campaign ID.
* `backup_set_id` - (Optional, String, ForceNew) Basic backup set ID.
* `db_node_set` - (Optional, List, ForceNew) Deployment information of the instance node, which will display the information of each AZ when the instance node is deployed across multiple AZs.
The information of AZ can be obtained from the `Zone` field in the return value of the [DescribeZones](https://intl.cloud.tencent.com/document/api/409/16769?from_cn_redirect=1) API.
* `instance_charge_type` - (Optional, String, ForceNew) Instance billing type, which currently supports:

- PREPAID: Prepaid, i.e., monthly subscription
- POSTPAID_BY_HOUR: Pay-as-you-go, i.e., pay by consumption

Default value: PREPAID.
* `name` - (Optional, String, ForceNew) Name of the newly purchased instance, which can contain up to 60 letters, digits, or symbols (-_). If this parameter is not specified, "Unnamed" will be displayed by default.
* `project_id` - (Optional, Int, ForceNew) Project ID.
* `recovery_target_time` - (Optional, String, ForceNew) Restoration point in time.
* `security_group_ids` - (Optional, Set: [`String`], ForceNew) Security group of the instance, which can be obtained from the `sgld` field in the return value of the [DescribeSecurityGroups](https://intl.cloud.tencent.com/document/api/215/15808?from_cn_redirect=1) API. If this parameter is not specified, the default security group will be bound.
* `sync_mode` - (Optional, String, ForceNew) Primary-standby sync mode, which supports:
Semi-sync: Semi-sync
Async: Asynchronous
Default value for the primary instance: Semi-sync
Default value for the read-only instance: Async.
* `tag_list` - (Optional, List, ForceNew) The information of tags to be bound with the instance, which is left empty by default. This parameter can be obtained from the `Tags` field in the return value of the [DescribeTags](https://intl.cloud.tencent.com/document/api/651/35316?from_cn_redirect=1) API.

The `db_node_set` object supports the following:

* `role` - (Required, String, ForceNew) Node type. Valid values:
`Primary`;
`Standby`.
* `zone` - (Required, String, ForceNew) AZ where the node resides, such as ap-guangzhou-1.
* `dedicated_cluster_id` - (Optional, String, ForceNew) Dedicated cluster ID.

The `tag_list` object supports the following:

* `tag_key` - (Required, String, ForceNew) Tag key.
* `tag_value` - (Required, String, ForceNew) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



