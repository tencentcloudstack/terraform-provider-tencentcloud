---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_origin_group_health_status"
sidebar_current: "docs-tencentcloud-datasource-teo_origin_group_health_status"
description: |-
  Use this data source to query the health status of origin groups under a TEO (EdgeOne) load balancer instance.
---

# tencentcloud_teo_origin_group_health_status

Use this data source to query the health status of origin groups under a TEO (EdgeOne) load balancer instance.

## Example Usage

### Query the health status of all origin groups under a load balancer

```hcl
data "tencentcloud_teo_origin_group_health_status" "example" {
  zone_id        = "zone-3fkff38fyw8s"
  lb_instance_id = "lb-instance-xxx"
}
```

### Query the health status of specified origin groups

```hcl
data "tencentcloud_teo_origin_group_health_status" "example" {
  zone_id          = "zone-3fkff38fyw8s"
  lb_instance_id   = "lb-instance-xxx"
  origin_group_ids = ["origin-group-xxx1", "origin-group-xxx2"]
}
```

## Argument Reference

The following arguments are supported:

* `lb_instance_id` - (Required, String) Load balancer instance ID.
* `zone_id` - (Required, String) Site ID.
* `origin_group_ids` - (Optional, List: [`String`]) Origin group ID list. When not specified, the health status of all origin groups under the load balancer is returned by default.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `origin_group_health_status_list` - Health status list of origin groups under the load balancer.
  * `check_region_health_status` - Health status of origins under each health check region.
    * `healthy` - Health status of origins under a single health check region. Valid values: Healthy, Unhealthy, Undetected. When all origins in a single health check region are healthy, the status is healthy; otherwise, it is unhealthy.
    * `origin_health_status` - Origin health status under the health check region.
      * `healthy` - Origin health status. Valid values: Healthy, Unhealthy, Undetected.
      * `origin` - Origin.
    * `region` - Health check region, ISO-3166-1 two-letter code.
  * `origin_group_id` - Origin group ID.
  * `origin_health_status` - The health status of each origin in the origin group, which is comprehensively determined based on the results of all detection regions. If more than half of the regions determine the origin as unhealthy, the corresponding status is unhealthy; otherwise, it is healthy.
    * `healthy` - Origin health status. Valid values: Healthy, Unhealthy, Undetected.
    * `origin` - Origin.


