---
subcategory: "Elasticsearch Service(ES)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_elasticsearch_start_logstash_pipeline_operation"
sidebar_current: "docs-tencentcloud-resource-elasticsearch_start_logstash_pipeline_operation"
description: |-
  Provides a resource to start a elasticsearch logstash pipeline
---

# tencentcloud_elasticsearch_start_logstash_pipeline_operation

Provides a resource to start a elasticsearch logstash pipeline

## Example Usage

```hcl
resource "tencentcloud_elasticsearch_start_logstash_pipeline_operation" "start_logstash_pipeline_operation" {
  instance_id = "ls-xxxxxx"
  pipeline_id = "xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance id.
* `pipeline_id` - (Required, String, ForceNew) Pipeline id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



