---
subcategory: "Real User Monitoring(RUM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_rum_offline_log_config_attachment"
sidebar_current: "docs-tencentcloud-resource-rum_offline_log_config_attachment"
description: |-
  Provides a resource to create a rum offline_log_config_attachment
---

# tencentcloud_rum_offline_log_config_attachment

Provides a resource to create a rum offline_log_config_attachment

## Example Usage

```hcl
resource "tencentcloud_rum_offline_log_config_attachment" "offline_log_config_attachment" {
  project_key = "ZEYrYfvaYQ30jRdmPx"
  unique_id   = "100027012454"
}
```

## Argument Reference

The following arguments are supported:

* `project_key` - (Required, String, ForceNew) Unique project key for reporting.
* `unique_id` - (Required, String, ForceNew) Unique identifier of the user to be listened on(aid or uin).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `msg` - Interface call information.


## Import

rum offline_log_config_attachment can be imported using the id, e.g.
```
$ terraform import tencentcloud_rum_offline_log_config_attachment.offline_log_config_attachment ZEYrYfvaYQ30jRdmPx#100027012454
```

