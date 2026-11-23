---
subcategory: "TDSQL Boundless(tdmysql)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmysql_specs"
sidebar_current: "docs-tencentcloud-data_source-tdmysql_specs"
description: |-
  Use this data source to query TDSQL Boundless (tdmysql) salable component specs.
---

# tencentcloud_tdmysql_specs

Use this data source to query TDSQL Boundless (tdmysql) salable component specs.

## Example Usage

```hcl
data "tencentcloud_tdmysql_specs" "example" {
  full_replications     = 2
  is_exclusive_instance = 1
  zones                 = ["ap-guangzhou-1"]
  instance_id           = "tdsql-o2mtest"
}

output "first_hybrid_spec_code" {
  value = data.tencentcloud_tdmysql_specs.example.hybrid_node_specs[0].spec_code
}
```

## Argument Reference

The following arguments are supported:

* `full_replications` - (Optional, Int) Number of replicas of hybrid instances. Valid range: 1-3, default is 3.
* `instance_id` - (Optional, String) ID of the instance.
* `is_exclusive_instance` - (Optional, Int) Whether the instance is exclusive (dedicated).
* `zones` - (Optional, Set) List of availability zone combinations.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `disk_type_category` - Disk type category. Valid values: CLOUD_DISK, LOCAL_DISK.
  * `instance_mode` - List of supported billing modes of the spec.
  * `max_ccu` - Range of maximum CCU values.
  * `max_node_num` - Maximum number of nodes.
  * `min_ccu` - Minimum CCU value.
  * `min_node_num` - Minimum number of nodes.
  * `spec_code` - Spec code.
  * `storage_node_cpu` - Number of CPU cores of the storage node.
  * `storage_node_default_disk` - Default disk size of the storage node, for display purposes.
  * `storage_node_max_disk` - Maximum disk size of the storage node.
  * `storage_node_max_num` - Maximum number of storage nodes.
  * `storage_node_mem` - Memory size of the storage node.
  * `storage_node_min_disk` - Minimum disk size of the storage node.
  * `storage_node_min_num` - Minimum number of storage nodes.
  * `storage_type` - Disk type. Valid values: CLOUD_HSSD (enhanced SSD), CLOUD_TCS (local SSD).
* `hybrid_node_specs` - List of salable specs of hybrid (peer) nodes.
* `serverless_ccu_spec` - List of salable specs of serverless nodes.
* `serverless_node_num_spec` - Serverless node number configuration. Contains at most one element because the API returns a single object.


