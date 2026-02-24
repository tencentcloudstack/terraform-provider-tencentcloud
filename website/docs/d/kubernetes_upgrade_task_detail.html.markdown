---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_upgrade_task_detail"
sidebar_current: "docs-tencentcloud-datasource-kubernetes_upgrade_task_detail"
description: |-
  Use this data source to query detailed information of TKE kubernetes upgrade task detail
---

# tencentcloud_kubernetes_upgrade_task_detail

Use this data source to query detailed information of TKE kubernetes upgrade task detail

## Example Usage

```hcl
data "tencentcloud_kubernetes_upgrade_task_detail" "example" {
  task_id = 21
}
```

## Argument Reference

The following arguments are supported:

* `task_id` - (Required, Int) Upgrade task ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `upgrade_plans` - Upgrade plans.
  * `cluster_id` - Cluster ID.
  * `cluster_name` - Cluster name.
  * `id` - Upgrade plan ID.
  * `planed_start_at` - Planned start time.
  * `reason` - Reason.
  * `region` - Cluster region.
  * `status` - Upgrade status.
  * `upgrade_end_at` - Upgrade end time.
  * `upgrade_start_at` - Upgrade start time.


