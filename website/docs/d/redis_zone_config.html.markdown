---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_zone_config"
sidebar_current: "docs-tencentcloud-datasource-redis_zone_config"
description: |-
  Use this data source to query which instance types of Redis are available in a specific region.
---

# tencentcloud_redis_zone_config

Use this data source to query which instance types of Redis are available in a specific region.

## Example Usage

```hcl
data "tencentcloud_redis_zone_config" "redislab" {
  region             = "ap-hongkong"
  result_output_file = "/temp/mytestpath"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, ForceNew) Name of a region. If this value is not set, the current region getting from provider's configuration will be used.
* `result_output_file` - (Optional, ForceNew) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of zone. Each element contains the following attributes:
  * `mem_sizes` - The memory volume of an available instance（in MB）.  
  * `type` - Instance type. Available values: master_slave_redis, master_slave_ckv, cluster_ckv, cluster_redis and standalone_redis. 
  * `version` - Version description of an available instance. Possible values: Redis 3.2, Redis 4.0 .
  * `zone` - ID of available zone.


