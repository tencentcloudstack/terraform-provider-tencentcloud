---
subcategory: "Oceanus"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_oceanus_tree_jobs"
sidebar_current: "docs-tencentcloud-datasource-oceanus_tree_jobs"
description: |-
  Use this data source to query detailed information of oceanus tree_jobs
---

# tencentcloud_oceanus_tree_jobs

Use this data source to query detailed information of oceanus tree_jobs

## Example Usage

```hcl
data "tencentcloud_oceanus_tree_jobs" "example" {
  work_space_id = "space-2idq8wbr"
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.
* `work_space_id` - (Optional, String) Workspace SerialId.


