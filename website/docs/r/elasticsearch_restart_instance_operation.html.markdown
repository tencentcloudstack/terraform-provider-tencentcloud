---
subcategory: "Elasticsearch Service(ES)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_elasticsearch_restart_instance_operation"
sidebar_current: "docs-tencentcloud-resource-elasticsearch_restart_instance_operation"
description: |-
  Provides a resource to restart a elasticsearch instance
---

# tencentcloud_elasticsearch_restart_instance_operation

Provides a resource to restart a elasticsearch instance

## Example Usage

```hcl
resource "tencentcloud_elasticsearch_restart_instance_operation" "restart_instance_operation" {
  instance_id = "es-xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance id.
* `force_restart` - (Optional, Bool, ForceNew) Force restart. Valid values:
- true: Forced restart;
- false: No forced restart;
default false.
* `restart_mode` - (Optional, Int, ForceNew) Restart mode: 0 roll restart; 1 full restart.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



