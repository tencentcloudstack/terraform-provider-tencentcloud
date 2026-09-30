---
subcategory: "CdwDoris"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cdwdoris_instance_nodes"
sidebar_current: "docs-tencentcloud-datasource-cdwdoris_instance_nodes"
description: |-
  Use this data source to query detailed information of CDWDoris cluster instance nodes.
---

# tencentcloud_cdwdoris_instance_nodes

Use this data source to query detailed information of CDWDoris cluster instance nodes.

## Example Usage

### Query cdwdoris instance nodes by instance id

```hcl
data "tencentcloud_cdwdoris_instance_nodes" "example" {
  instance_id = "cdwdoris-rhbflamd"
}
```

### Query cdwdoris instance nodes with optional filters

```hcl
data "tencentcloud_cdwdoris_instance_nodes" "example" {
  instance_id    = "cdwdoris-rhbflamd"
  node_role      = "data"
  display_policy = "All"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Cluster instance ID, `cdw-xxxx` string type.
* `display_policy` - (Optional, String) Display policy, `All` to display all.
* `node_role` - (Optional, String) Cluster role type, default is `data` data node.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_nodes_list` - List of cluster node information.
  * `compute_group_id` - Compute group ID. Note: This field may return null, indicating that no valid values can be obtained.
  * `core` - Number of CPU cores. Note: This field may return null, indicating that no valid values can be obtained.
  * `create_time` - Creation time. Note: This field may return null, indicating that no valid values can be obtained.
  * `disk_size` - Disk size. Note: This field may return null, indicating that no valid values can be obtained.
  * `disk_type` - Disk type. Note: This field may return null, indicating that no valid values can be obtained.
  * `fe_role` - FE node role. Note: This field may return null, indicating that no valid values can be obtained.
  * `ip` - IP address. Note: This field may return null, indicating that no valid values can be obtained.
  * `memory` - Memory size. Note: This field may return null, indicating that no valid values can be obtained.
  * `rip` - Rip. Note: This field may return null, indicating that no valid values can be obtained.
  * `role` - Name of the clickhouse cluster to which it belongs. Note: This field may return null, indicating that no valid values can be obtained.
  * `spec` - Model, such as S1. Note: This field may return null, indicating that no valid values can be obtained.
  * `status` - Status. Note: This field may return null, indicating that no valid values can be obtained.
  * `uuid` - UUID. Note: This field may return null, indicating that no valid values can be obtained.
  * `virtual_zone` - Virtual availability zone. Note: This field may return null, indicating that no valid values can be obtained.
  * `zone` - Availability zone. Note: This field may return null, indicating that no valid values can be obtained.
* `node_roles` - List of cluster-supported node role types.


