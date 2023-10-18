---
subcategory: "Elasticsearch Service(ES)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_elasticsearch_restart_logstash_instance_operation"
sidebar_current: "docs-tencentcloud-resource-elasticsearch_restart_logstash_instance_operation"
description: |-
  Provides a resource to restart a elasticsearch logstash instance
---

# tencentcloud_elasticsearch_restart_logstash_instance_operation

Provides a resource to restart a elasticsearch logstash instance

## Example Usage

```hcl
resource "tencentcloud_elasticsearch_restart_logstash_instance_operation" "restart_logstash_instance_operation" {
  instance_id = "ls-xxxxxx"
  type        = 0
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance id.
* `type` - (Required, Int, ForceNew) Restart type, 0 full restart, 1 rolling restart.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



