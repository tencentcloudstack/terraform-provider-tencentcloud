---
subcategory: "Elasticsearch Service(ES)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_elasticsearch_restart_kibana_operation"
sidebar_current: "docs-tencentcloud-resource-elasticsearch_restart_kibana_operation"
description: |-
  Provides a resource to restart a elasticsearch kibana
---

# tencentcloud_elasticsearch_restart_kibana_operation

Provides a resource to restart a elasticsearch kibana

## Example Usage

```hcl
resource "tencentcloud_elasticsearch_restart_kibana_operation" "restart_kibana_operation" {
  instance_id = "es-xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



