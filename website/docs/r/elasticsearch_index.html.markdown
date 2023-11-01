---
subcategory: "Elasticsearch Service(ES)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_elasticsearch_index"
sidebar_current: "docs-tencentcloud-resource-elasticsearch_index"
description: |-
  Provides a resource to create a elasticsearch index
---

# tencentcloud_elasticsearch_index

Provides a resource to create a elasticsearch index

## Example Usage

```hcl
resource "tencentcloud_elasticsearch_index" "index" {
  instance_id     = "es-xxxxxx"
  index_type      = "normal"
  index_name      = "test-es-index"
  index_meta_json = "{\"mappings\":{},\"settings\":{\"index.number_of_replicas\":1,\"index.number_of_shards\":1,\"index.refresh_interval\":\"30s\"}}"
}
```

## Argument Reference

The following arguments are supported:

* `index_name` - (Required, String) index name to create.
* `index_type` - (Required, String) type of the index to be created. auto: autonomous index. normal: indicates a common index.
* `instance_id` - (Required, String) es instance id.
* `index_meta_json` - (Optional, String) Create index metadata JSON, such as mappings, settings.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

elasticsearch index can be imported using the id, e.g.

```
terraform import tencentcloud_elasticsearch_index.index index_id
```

