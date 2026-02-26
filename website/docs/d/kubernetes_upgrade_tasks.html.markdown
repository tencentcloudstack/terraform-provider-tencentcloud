---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_upgrade_tasks"
sidebar_current: "docs-tencentcloud-datasource-kubernetes_upgrade_tasks"
description: |-
  Use this data source to query detailed information of TKE kubernetes upgrade tasks
---

# tencentcloud_kubernetes_upgrade_tasks

Use this data source to query detailed information of TKE kubernetes upgrade tasks

## Example Usage

```hcl
data "tencentcloud_kubernetes_upgrade_tasks" "example" {}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `upgrade_tasks` - Upgrade tasks.
  * `component` - Component name.
  * `created_at` - Creation time.
  * `id` - Task ID.
  * `name` - Task name.
  * `planed_start_at` - Planned start time.
  * `related_resources` - Related resources.
  * `upgrade_impact` - Upgrade impact.


