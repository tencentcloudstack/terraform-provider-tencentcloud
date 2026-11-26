---
subcategory: "TDSQL Boundless(tdmysql)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmysql_specs"
sidebar_current: "docs-tencentcloud-data_source-tdmysql_specs"
description: |-
  Provides a data source to query the sellable component specs of TDSQL Boundless (tdmysql).
---

# tencentcloud_tdmysql_specs

Provides a data source to query the sellable component specs of TDSQL Boundless (tdmysql).

## Example Usage

### Query all specs

```hcl
data "tencentcloud_tdmysql_specs" "example" {}
```

### Query specs by filters

```hcl
data "tencentcloud_tdmysql_specs" "example" {
  full_replications     = 3
  is_exclusive_instance = 0
  zones                 = ["ap-guangzhou-6", "ap-guangzhou-7"]
  instance_id           = "tdsql3-gpyfw59om6"
}
```

## Argument Reference

The following arguments are supported:

* `full_replications` - (Optional, Int) Number of replicas of the all-round (hybrid) node, currently supported range: 1-3, default is 3.
* `instance_id` - (Optional, String) Instance ID.
* `is_exclusive_instance` - (Optional, Int) Exclusive instance.
* `zones` - (Optional, List) List of availability zone combinations.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `disk_type_category` - Disk type category: CLOUD_DISK for cloud disk, LOCAL_DISK for local disk.
  * `instance_mode` - List of billing modes supported by the spec.
  * `max_ccu` - Range of maximum CCU values.
  * `max_node_num` - Maximum number of nodes.
  * `min_ccu` - Minimum CCU value.
  * `min_node_num` - Minimum number of nodes.
  * `spec_code` - Spec code.
  * `storage_node_cpu` - Number of CPU cores of the storage node.
  * `storage_node_default_disk` - Default disk size of the storage node, used for frontend display.
  * `storage_node_max_disk` - Upper limit of the storage node disk size.
  * `storage_node_max_num` - Maximum number of storage nodes.
  * `storage_node_mem` - Memory size of the storage node.
  * `storage_node_min_disk` - Lower limit of the storage node disk size.
  * `storage_node_min_num` - Minimum number of storage nodes.
  * `storage_type` - Disk type, CLOUD_HSSD enhanced SSD, CLOUD_TCS local SSD disk.
* `hybrid_node_specs` - List of hybrid node sellable specs.
* `id` - ID of the data source, a random token.
* `serverless_ccu_spec` - List of serverless node sellable CCU specs.
* `serverless_node_num_spec` - Serverless node number configuration.


