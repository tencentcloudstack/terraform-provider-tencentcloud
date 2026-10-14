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


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DescribeInstances](https://cloud.tencent.com/document/api/845/30631) - Query all instances meeting conditions for the user in that region
* Create: [RestartKibana](https://cloud.tencent.com/document/api/845/50224)

