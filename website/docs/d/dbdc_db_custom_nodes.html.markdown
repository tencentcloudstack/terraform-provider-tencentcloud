---
subcategory: "Database Dedicated Cluster(DBDC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbdc_db_custom_nodes"
sidebar_current: "docs-tencentcloud-datasource-dbdc_db_custom_nodes"
description: |-
  Use this data source to query DB Custom node list from TencentCloud DBDC product.
---

# tencentcloud_dbdc_db_custom_nodes

Use this data source to query DB Custom node list from TencentCloud DBDC product.

## Example Usage

### Query all dbdc db custom nodes

```hcl
data "tencentcloud_dbdc_db_custom_nodes" "example" {}
```

### Query dbdc db custom nodes by node_ids

```hcl
data "tencentcloud_dbdc_db_custom_nodes" "example" {
  node_ids = [
    "dbcn-abc12345",
    "dbcn-def67890"
  ]
}
```

### Query dbdc db custom nodes by filters

```hcl
data "tencentcloud_dbdc_db_custom_nodes" "example" {
  filters {
    name   = "cluster-id"
    values = ["dbcc-nmtmsew8"]
  }

  filters {
    name   = "status"
    values = ["Running"]
  }
}
```

### Query dbdc db custom nodes by tags

```hcl
data "tencentcloud_dbdc_db_custom_nodes" "example" {
  tags {
    key   = "env"
    value = "production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter conditions. Supported filter names: cluster-id, node-name (exact match), status (Creating, Running, Isolating, Isolated, Activating, Destroying), zone.
* `node_ids` - (Optional, List: [`String`]) Query by one or more Node IDs. Maximum 100 IDs per request.
* `result_output_file` - (Optional, String) Used to save results.
* `tags` - (Optional, List) Filter by tag Key and Value.

The `filters` object supports the following:

* `name` - (Required, String) Filter field name.
* `values` - (Required, List) Filter field values.

The `tags` object supports the following:

* `key` - (Required, String) Tag key.
* `value` - (Required, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `node_set` - DB Custom node list.
  * `auto_renew` - Auto-renew flag. 1=auto renew, 0=no auto renew.
  * `charge_type` - Payment type.
  * `cluster_id` - Cluster ID the node belongs to.
  * `cpu` - Node CPU size in cores.
  * `created_time` - Node creation time.
  * `data_disks` - Data disk list. Note: This field may return null, indicating that no valid value can be obtained.
    * `disk_name` - Disk name.
    * `disk_size` - Disk size in GiB.
    * `disk_type` - Disk type.
  * `expire_time` - Node expiration time.
  * `host_ip` - Host IP.
  * `image_id` - Node OS image ID.
  * `isolated_time` - Node isolation time.
  * `lan_ip` - Internal network IP address of the node.
  * `memory` - Node memory size in GiB.
  * `node_id` - Node ID.
  * `node_name` - Node name.
  * `node_type` - Node type/spec.
  * `os_name` - Node OS name.
  * `rack_id` - Rack ID.
  * `ssh_endpoint` - SSH endpoint for accessing the node, format: IP:Port.
  * `status` - Node status.
  * `subnet_id` - Subnet ID the node SSH endpoint belongs to.
  * `switch_id` - Switch ID.
  * `system_disk` - System disk info. Note: This field may return null, indicating that no valid value can be obtained.
    * `disk_size` - Disk size in GiB.
    * `disk_type` - Disk type.
  * `tags` - Node tag information. Note: This field may return null, indicating that no valid value can be obtained.
    * `key` - Tag key.
    * `value` - Tag value.
  * `vpc_id` - VPC ID the node SSH endpoint belongs to.
  * `zone` - Availability zone the node belongs to.


