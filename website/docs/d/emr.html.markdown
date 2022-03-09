---
subcategory: "EMR"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_emr"
sidebar_current: "docs-tencentcloud-datasource-emr"
description: |-
  Provides an available EMR for the user.
---

# tencentcloud_emr

Provides an available EMR for the user.

The EMR data source fetch proper EMR from user's EMR pool.

## Example Usage

```hcl
data "tencentcloud_emr" "my_emr" {
  display_strategy = "clusterList"
  instance_ids     = ["emr-rnzqrleq"]
}
```

## Argument Reference

The following arguments are supported:

* `display_strategy` - (Required) Display strategy(e.g.:clusterList, monitorManage).
* `instance_ids` - (Optional) fetch all instances with same prefix(e.g.:emr-xxxxxx).
* `project_id` - (Optional) Fetch all instances which owner same project. Default 0 meaning use default project id.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `clusters` - A list of clusters will be exported and its every element contains the following attributes:
  * `add_time` - Add time of instance.
  * `charge_type` - Charge type of instance.
  * `cluster_id` - Cluster id of instance.
  * `cluster_name` - Cluster name of instance.
  * `ftitle` - Title of instance.
  * `id` - Id of instance.
  * `master_ip` - Master ip of instance.
  * `project_id` - Project id of instance.
  * `region_id` - Region id of instance.
  * `status` - Status of instance.
  * `zone_id` - Zone id of instance.
  * `zone` - Zone of instance.


