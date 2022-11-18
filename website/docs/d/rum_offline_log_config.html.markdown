---
subcategory: "Real User Monitoring(RUM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_rum_offline_log_config"
sidebar_current: "docs-tencentcloud-datasource-rum_offline_log_config"
description: |-
  Use this data source to query detailed information of rum offlineLogConfig
---

# tencentcloud_rum_offline_log_config

Use this data source to query detailed information of rum offlineLogConfig

## Example Usage

```hcl
data "tencentcloud_rum_offline_log_config" "offlineLogConfig" {
  project_key = "ZEYrYfvaYQ30jRdmPx"
}
```

## Argument Reference

The following arguments are supported:

* `project_key` - (Required, String) Unique project key for reporting.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `msg` - API call information.
* `unique_id_set` - Unique identifier of the user to be listened on(aid or uin).


