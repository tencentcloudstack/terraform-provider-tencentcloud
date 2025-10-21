---
subcategory: "CDC"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cdc_dedicated_cluster_hosts"
sidebar_current: "docs-tencentcloud-datasource-cdc_dedicated_cluster_hosts"
description: |-
  Use this data source to query detailed information of CDC dedicated cluster hosts
---

# tencentcloud_cdc_dedicated_cluster_hosts

Use this data source to query detailed information of CDC dedicated cluster hosts

## Example Usage

```hcl
data "tencentcloud_cdc_dedicated_cluster_hosts" "hosts" {
  dedicated_cluster_id = "cluster-262n63e8"
}
```

## Argument Reference

The following arguments are supported:

* `dedicated_cluster_id` - (Required, String) Dedicated Cluster ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `host_info_set` - Dedicated Cluster Host Info.
  * `cpu_available` - Dedicated Cluster Host CPU Available Count.
  * `cpu_total` - Dedicated Cluster Host CPU Total Count.
  * `expire_time` - Dedicated Cluster Host Expire Time.
  * `host_id` - Dedicated Cluster Host ID.
  * `host_ip` - Dedicated Cluster Host Ip (Deprecated).
  * `host_status` - Dedicated Cluster Host Status.
  * `host_type` - Dedicated Cluster Host Type.
  * `mem_available` - Dedicated Cluster Host Memory Available Count (GB).
  * `mem_total` - Dedicated Cluster Host Memory Total Count (GB).
  * `run_time` - Dedicated Cluster Host Run Time.
  * `service_type` - Dedicated Cluster Service Type.


