---
subcategory: "Elasticsearch Service(ES)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_elasticsearch_diagnose_instance"
sidebar_current: "docs-tencentcloud-resource-elasticsearch_diagnose_instance"
description: |-
  Provides a resource to create a elasticsearch diagnose instance
---

# tencentcloud_elasticsearch_diagnose_instance

Provides a resource to create a elasticsearch diagnose instance

## Example Usage

```hcl
resource "tencentcloud_elasticsearch_diagnose_instance" "diagnose_instance" {
  instance_id      = "es-xxxxxx"
  diagnose_jobs    = ["cluster_health"]
  diagnose_indices = "*"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance id.
* `diagnose_indices` - (Optional, String, ForceNew) Indexes that need to be diagnosed. Wildcards are supported.
* `diagnose_jobs` - (Optional, Set: [`String`], ForceNew) Diagnostic items that need to be triggered.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



