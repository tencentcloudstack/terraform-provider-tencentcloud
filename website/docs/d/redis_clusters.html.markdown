---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_clusters"
sidebar_current: "docs-tencentcloud-datasource-redis_clusters"
description: |-
  Use this data source to query detailed information of redis clusters
---

# tencentcloud_redis_clusters

Use this data source to query detailed information of redis clusters

## Example Usage

### Query all instance

```hcl
data "tencentcloud_redis_clusters" "clusters" {}
```

### Also Support the following query conditions

```hcl
data "tencentcloud_redis_clusters" "clusters" {
  dedicated_cluster_id = "cluster-0astoh6a"
  redis_cluster_ids    = ["crs-cdc-9nyfki8h"]
  cluster_name         = "crs-cdc-9nyfki8h"
  project_ids          = [0, 1]
  status               = [0, 1, 2]
  auto_renew_flag      = [0, 1, 2]
}
```

## Argument Reference

The following arguments are supported:

* `auto_renew_flag` - (Optional, Set: [`Int`]) Renewal mode: 0- default state (manual renewal); 1- Automatic renewal; 2- Clearly stating that automatic renewal is not allowed.
* `cluster_name` - (Optional, String) Cluster name.
* `dedicated_cluster_id` - (Optional, String) Dedicated cluster Id.
* `project_ids` - (Optional, Set: [`Int`]) Project Ids.
* `redis_cluster_ids` - (Optional, Set: [`String`]) Redis Cluster Ids.
* `result_output_file` - (Optional, String) Used to save results.
* `status` - (Optional, Set: [`Int`]) Cluster status: 1- In process, 2- Running, 3- Isolated.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `resources` - .
  * `app_id` - User's Appid.
  * `auto_renew_flag` - Renewal mode: 0- default state (manual renewal); 1- Automatic renewal; 2- Clearly stating that automatic renewal is not allowed.
  * `base_bundles` - Basic Control Resource Package.
    * `available_memory` - Saleable memory, unit: GB.
    * `count` - Resource bundle count.
    * `resource_bundle_name` - Resource bundle name.
  * `cluster_name` - Cluster name.
  * `dedicated_cluster_id` - Dedicated cluster Id.
  * `end_time` - Instance expiration time.
  * `pay_mode` - Billing mode, 1-annual and monthly package, 0-quantity based billing.
  * `project_id` - Project Id.
  * `redis_cluster_id` - Redis Cluster Id.
  * `region_id` - Region Id.
  * `resource_bundles` - List of Resource Packages.
    * `available_memory` - Saleable memory, unit: GB.
    * `count` - Resource bundle count.
    * `resource_bundle_name` - Resource bundle name.
  * `start_time` - Instance create time.
  * `status` - Cluster status: 1- In process, 2- Running, 3- Isolated.
  * `zone_id` - zone Id.


