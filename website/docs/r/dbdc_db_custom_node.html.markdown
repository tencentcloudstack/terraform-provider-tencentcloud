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

~> **NOTE:** Both create and delete operations are asynchronous. One resource maps to exactly one node (`node_count` defaults to `1`). The delete operation isolates the node first and then destroys it, waiting for the underlying task to reach the `Succeeded` status.

## Example Usage

```hcl
resource "tencentcloud_dbdc_db_custom_node" "example" {
  zone      = "ap-shanghai-5"
  image_id  = "img-xxxxxxxx"
  vpc_id    = "vpc-xxxxxxxx"
  subnet_id = "subnet-xxxxxxxx"
  node_type = "DB.AT5.8XLARGE128"
  period    = 1
  node_name = "tf-example"

  login_settings {
    password = "Passw0rd@2024"
  }

  auto_renew = 0

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
* `auto_renew` - (Optional, Int) Auto-renew flag. Valid values: `1` (auto-renew), `0` (not auto-renew). Mutable via the renew API.
* `auto_voucher` - (Optional, Int, ForceNew) Whether to use voucher to deduct automatically. Valid values: `1` (use), `0` (not use). Default value is `0`.
* `login_settings` - (Optional, List, ForceNew) Instance login settings. You can set the login method to password, key, or keep the original image login settings. Only one method can be set.
* `node_count` - (Optional, Int, ForceNew) Number of nodes to purchase. One resource maps to exactly one node, so the default value is `1`.
* `node_name` - (Optional, String, ForceNew) Node name. Up to 128 characters.
* `period` - (Optional, Int, ForceNew) Purchase duration in months. Valid values: 1/2/3/4/5/6/7/8/9/10/11/12/24/36. Default value is `1`.
* `tags` - (Optional, Map) Node tags.
* `voucher_ids` - (Optional, List, ForceNew) Voucher ID list. Must be undeducted voucher IDs owned by the current account.

The `login_settings` object supports the following:

* `keep_image_login` - (Optional, String, ForceNew) Whether to keep the original login settings of the image. Valid values: `true`, `false`. Cannot be specified together with Password or KeyIds.
* `key_ids` - (Optional, List, ForceNew) Key pair ID list. Only a single ID is supported currently. Password and key cannot be specified at the same time.
* `password` - (Optional, String, ForceNew) Instance login password. Password complexity limits vary by operating system type.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `charge_type` - Charge type. Valid values: `PREPAID`.
* `cluster_id` - Cluster ID that the node belongs to.
* `cpu` - Node CPU size, unit: core.
* `created_time` - Node creation time.
* `data_disks` - Data disk information.
  * `disk_name` - Disk name.
  * `disk_size` - Disk size, unit: GiB.
  * `disk_type` - Disk type.
* `expire_time` - Node expiration time.
* `isolated_time` - Node isolation time.
* `lan_ip` - Intranet communication IP address of the node.
* `memory` - Node memory, unit: GiB.
* `node_id` - Node ID.
* `os_name` - Operating system name of the node.
* `ssh_endpoint` - SSH endpoint to access this node, in the format `IP:Port`.
* `status` - Node status. Valid values: `Creating`, `Running`, `Isolating`, `Isolated`, `Activating`, `Destroying`.
* `system_disk` - System disk information.
  * `disk_size` - Disk size, unit: GiB.
  * `disk_type` - Disk type.


## Timeouts

This resource provides the following [Timeouts](https://www.terraform.io/docs/configuration/blocks/resources/syntax.html#operation-timeouts) configuration options:

- `create` - (Default `30m`)
- `update` - (Default `30m`)
- `delete` - (Default `30m`)

## Import

DBDC db custom node can be imported using the id, e.g.

```
terraform import tencentcloud_dbdc_db_custom_node.example dbcn-xxxxxxxx
```

