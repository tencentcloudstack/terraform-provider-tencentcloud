---
subcategory: "Database Dedicated Cluster(DBDC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbdc_db_custom_node"
sidebar_current: "docs-tencentcloud-resource-dbdc_db_custom_node"
description: |-
  Provides a resource to create a DBDC db custom node.
---

# tencentcloud_dbdc_db_custom_node

Provides a resource to create a DBDC db custom node.

~> **NOTE:** Not all `zone` support configuration `system_disk` and `data_disks`. You can check the supported zone list with `tencentcloud_dbdc_db_custom_node_types`.

## Example Usage

### Create a PREPAID DBDC db custom node

```hcl
resource "tencentcloud_dbdc_db_custom_node" "example" {
  zone        = "ap-shanghai-5"
  image_id    = "img-rm13akp3"
  vpc_id      = "vpc-cseo7req"
  subnet_id   = "subnet-huka6qhj"
  node_type   = "DB.AT5.8XLARGE128"
  period      = 1
  auto_renew  = 1
  charge_type = "PREPAID"
  node_name   = "tf-example"

  login_settings {
    password = "Password@2026"
  }

  tags = {
    createBy = "Terraform"
  }
}
```

### Create a POSTPAID DBDC db custom node

```hcl
resource "tencentcloud_dbdc_db_custom_node" "example" {
  zone         = "ap-bangkok-2"
  image_id     = "img-7rqxtnh9"
  vpc_id       = "vpc-doprnsrq"
  subnet_id    = "subnet-kzp25bun"
  node_type    = "DB.SA5.8XLARGE128"
  node_name    = "tf-example"
  host_name    = "hostName"
  network_mode = "cross_tenant_eni"
  charge_type  = "POSTPAID"

  login_settings {
    password = "Password@2026"
  }

  system_disk {
    disk_size = 100
    disk_type = "CLOUD_HSSD"
  }

  data_disks {
    disk_size = 200
    disk_type = "CLOUD_HSSD"
  }

  security_group_ids = [
    "sg-avup6l04",
  ]

  tags = {
    createBy = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `image_id` - (Required, String, ForceNew) Image ID, format `img-xxxxxxx`. Must be an image owned by the DB Custom product under the current account.
* `node_type` - (Required, String, ForceNew) Node spec, e.g. `DB.AT5.8XLARGE128`, `DB.AT5.16XLARGE256`, `DB.AT5.32XLARGE512`, `DB.AT5.64XLARGE1152`, `DB.AT5.128XLARGE2304`.
* `subnet_id` - (Required, String, ForceNew) Subnet ID used to establish the SSH connection for the node. Must belong to the VPC and match the availability zone.
* `vpc_id` - (Required, String, ForceNew) VPC ID used to establish the SSH connection for the node. Must be owned by the current account and cannot be cross-region.
* `zone` - (Required, String, ForceNew) Availability zone supported by the product, e.g. `ap-shanghai-5`, `ap-shanghai-8`, `ap-nanjing-3`.
* `auto_renew` - (Optional, Int) Auto-renew flag. Valid values: `1` (auto-renew), `2` (not auto-renew). Mutable via the renew API.
* `auto_voucher` - (Optional, Int) Whether to use voucher to deduct automatically. Valid values: `1` (use), `0` (not use). Default value is `0`.
* `charge_type` - (Optional, String, ForceNew) Charge type. Valid values: `PREPAID` (subscription, default), `POSTPAID` (pay-as-you-go).
* `data_disks` - (Optional, List, ForceNew) Data disk configuration. Only cloud-disk node types (e.g. `DB.SA5`) support setting this; local-disk types (e.g. `DB.AT5`) do not. Refreshed from the `DescribeDBCustomNodes` API response. Note: `disk_name` is read-only and ignored as a create input.
* `host_name` - (Optional, String, ForceNew) Hostname of the node. Dots (`.`) and hyphens (`-`) cannot be the first/last character or be used consecutively; underscores (`_`) are not allowed. Windows: 2-15 chars (letters, digits, `-`, no `.`); Linux/others: 2-60 chars (supports multiple dot-separated segments). Write-only: not returned by `DescribeDBCustomNodes`, so the configured value is preserved in state.
* `login_settings` - (Optional, List, ForceNew) Instance login settings. You can set the login method to password, key, or keep the original image login settings. Only one method can be set.
* `network_mode` - (Optional, String, ForceNew) Node network mode. Valid values: `privatelink` (four-layer SSH connectivity), `cross_tenant_eni` (three-layer dual-NIC access). Default is `privatelink`. Refreshed from the `DescribeDBCustomNodes` API response.
* `node_name` - (Optional, String, ForceNew) Node name. Up to 128 characters.
* `period` - (Optional, Int) Purchase duration in months. Valid values: 1/2/3/4/5/6/7/8/9/10/11/12/24/36. Default value is `1`.
* `security_group_ids` - (Optional, Set: [`String`]) Set of security group IDs to bind to the node. Treated as an unordered set; HCL element order has no semantic meaning. Mutable via the `ModifyDBCustomNodeSecurityGroups` API; refreshed from the `DescribeDBCustomNodeSecurityGroups` API.
* `system_disk` - (Optional, List, ForceNew) System disk configuration. Only cloud-disk node types (e.g. `DB.SA5`) support setting this; local-disk types (e.g. `DB.AT5`) do not. Refreshed from the `DescribeDBCustomNodes` API response.
* `tags` - (Optional, Map) Node tags.
* `voucher_ids` - (Optional, List: [`String`]) Voucher ID list. Must be undeducted voucher IDs owned by the current account.

The `data_disks` object supports the following:

* `disk_size` - (Optional, Int, ForceNew) Disk size, unit: GiB.
* `disk_type` - (Optional, String, ForceNew) Disk type. Valid values: `CLOUD_HSSD` (enhanced cloud disk), `LOCAL_NVME` (local disk).

The `login_settings` object supports the following:

* `keep_image_login` - (Optional, String, ForceNew) Whether to keep the original login settings of the image. Valid values: `true`, `false`. Cannot be specified together with Password or KeyIds.
* `key_ids` - (Optional, List, ForceNew) Key pair ID list. Only a single ID is supported currently. Password and key cannot be specified at the same time.
* `password` - (Optional, String, ForceNew) Instance login password. Password complexity limits vary by operating system type.

The `system_disk` object supports the following:

* `disk_size` - (Optional, Int, ForceNew) Disk size, unit: GiB.
* `disk_type` - (Optional, String, ForceNew) Disk type. Valid values: `CLOUD_HSSD` (enhanced cloud disk).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cluster_id` - Cluster ID that the node belongs to.
* `cpu` - Node CPU size, unit: core.
* `created_time` - Node creation time.
* `eni_ip` - Node access IP address when the `NetworkModeCrossTenantENI` network mode is selected. Refreshed from the `DescribeDBCustomNodes` API response.
* `expire_time` - Node expiration time.
* `isolated_time` - Node isolation time.
* `lan_ip` - Intranet communication IP address of the node.
* `memory` - Node memory, unit: GiB.
* `node_id` - Node ID.
* `os_name` - Operating system name of the node.
* `ssh_endpoint` - SSH endpoint to access this node, in the format `IP:Port`.
* `status` - Node status. Valid values: `Creating`, `Running`, `Isolating`, `Isolated`, `Activating`, `Destroying`.

The `data_disks` object exports the following:

* `disk_name` - Disk name. Read-only; ignored as a create input (per API).

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `1h0m`) Used when creating the resource.
* `update` - (Defaults to `1h0m`) Used when updating the resource.
* `delete` - (Defaults to `1h0m`) Used when deleting the resource.

## Import

DBDC db custom node can be imported using the id, e.g.

```
terraform import tencentcloud_dbdc_db_custom_node.example dbcn-ttiyh58n
```

