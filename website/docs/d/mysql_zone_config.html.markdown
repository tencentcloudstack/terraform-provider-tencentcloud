---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_zone_config"
sidebar_current: "docs-tencentcloud-datasource-mysql_zone_config"
description: |-
  Use this data source to query the available database specifications for different regions. And a maximum of 20 requests can be initiated per second for this query.
---

# tencentcloud_mysql_zone_config

Use this data source to query the available database specifications for different regions. And a maximum of 20 requests can be initiated per second for this query.

## Example Usage

```hcl
data "tencentcloud_mysql_zone_config" "mysql" {
  region             = "ap-guangzhou"
  result_output_file = "mytestpath"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Region parameter, which is used to identify the region to which the data you want to work with belongs.
* `result_output_file` - (Optional, String) Used to store results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of zone config. Each element contains the following attributes:
  * `disaster_recovery_zones` - Information about available zones of recovery.
  * `engine_versions` - The version number of the database engine to use. Supported versions include `5.5`/`5.6`/`5.7`.
  * `first_slave_zones` - Zone information about first slave instance.
  * `is_default` - Indicates whether the current DC is the default DC for the region. Possible returned values: `0` - no; `1` - yes.
  * `is_support_disaster_recovery` - Indicates whether recovery is supported: `0` - No; `1` - Yes.
  * `is_support_vpc` - Indicates whether VPC is supported: `0` - No; `1` - Yes.
  * `name` - The name of available zone which is equal to a specific datacenter.
  * `remote_ro_zones` - Zone information about remote ro instance.
  * `second_slave_zones` - Zone information about second slave instance.
  * `sells` - A list of supported instance types for sell:
    * `max_volume_size` - Maximum disk size (in GB).
    * `mem_size` - Memory size (in MB).
    * `min_volume_size` - Minimum disk size (in GB).
    * `qps` - Queries per second.
    * `volume_step` - Disk increment (in GB).
  * `slave_deploy_modes` - Availability zone deployment method. Available values: `0` - Single availability zone; `1` - Multiple availability zones.
  * `support_slave_sync_modes` - Data replication mode. `0` - Async replication; `1` - Semisync replication; `2` - Strongsync replication.


