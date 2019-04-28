---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_zone_config"
sidebar_current: "docs-tencentcloud-tencentcloud_mysql_zone_config"
description: |-
  Use this data source to query the available database specifications for different regions. And a maximum of 20 requests can be initiated per second for this query.
---

# tencentcloud_mysql_zone_config

Use this data source to query the available database specifications for different regions. And a maximum of 20 requests can be initiated per second for this query.

## Example Usage
```
data "tencentcloud_mysql_zone_config" "mysql" {
    region = "ap-guangzhou"
    result_output_file = "mytestpath" 
}
```
## Argument Reference

The following arguments are supported:

- `region` - (Optional) Region parameter, which is used to identify the region to which the data you want to work with belongs. 
- `result_output_file` - (Optional) Used to store results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `name` - The name of available zone which is equal to a specific datacenter.
- `is_default` - Indicates whether the current DC is the default DC for the region. Possible returned values: 0 - No; 1 - Yes.
- `is_support_disaster_recovery` - Indicates whether recovery is supported: 0 - No; 1 - Yes.  
- `is_support_vpc` - Indicates whether VPC is supported: 0 - No; 1 - Yes.
- `engine_versions` - The version number of the database engine to use. Supported versions include 5.5/5.6/5.7.
- `support_slave_sync_modes` - Data replication mode. 0 - Async replication; 1 - Semisync replication; 2 - Strongsync replication.
- `disaster_recovery_zones` - Information about available zones of recovery.
- `slave_deploy_modes` - Availability zone deployment method. Available values: 0 - Single availability zone; 1 - Multiple availability zones.
- `first_slave_zones` - Zone information about first slave instance.
- `second_slave_zones` - Zone information about second slave instance.
- `sells` - a list of supported instance types for sell.  


For supported instance types, the following information will be included:

- `mem_size` - Memory size (in MB).
- `min_volume_size` - Minimum disk size (in GB).
- `max_volume_size` - Maximum disk size (in GB).
- `volume_step` - Disk increment (in GB).
- `qps` - Queries per second.
