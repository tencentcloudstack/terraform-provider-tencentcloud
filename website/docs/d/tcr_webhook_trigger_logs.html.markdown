---
subcategory: "Tencent Container Registry(TCR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_webhook_trigger_logs"
sidebar_current: "docs-tencentcloud-datasource-tcr_webhook_trigger_logs"
description: |-
  Use this data source to query detailed information of tencentcloud_tcr_webhook_trigger_logs
---

# tencentcloud_tcr_webhook_trigger_logs

Use this data source to query detailed information of tencentcloud_tcr_webhook_trigger_logs

## Example Usage

```hcl
data "tencentcloud_tcr_webhook_trigger_logs" "my_logs" {
  registry_id = local.tcr_id
  namespace   = var.tcr_namespace
  trigger_id  = var.trigger_id
  tags = {
    "createdBy" = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `namespace` - (Required, String) namespace.
* `registry_id` - (Required, String) instance Id.
* `trigger_id` - (Required, Int) trigger id.
* `result_output_file` - (Optional, String) Used to save results.
* `tags` - (Optional, Map) Tag description list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `logs` - log list.
  * `creation_time` - creation time.
  * `detail` - webhook trigger detail.
  * `event_type` - event type.
  * `id` - log id.
  * `notify_type` - notification type.
  * `status` - status.
  * `trigger_id` - trigger Id.
  * `update_time` - update time.


