---
subcategory: "Elasticsearch Service(ES)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_elasticsearch_logstash_pipeline"
sidebar_current: "docs-tencentcloud-resource-elasticsearch_logstash_pipeline"
description: |-
  Provides a resource to create a elasticsearch logstash pipeline
---

# tencentcloud_elasticsearch_logstash_pipeline

Provides a resource to create a elasticsearch logstash pipeline

## Example Usage

```hcl
resource "tencentcloud_elasticsearch_logstash_pipeline" "logstash_pipeline" {
  instance_id = "ls-xxxxxx"
  pipeline {
    pipeline_id              = "logstash-pipeline-test"
    pipeline_desc            = ""
    config                   = <<EOF
input{

}
filter{

}
output{

}
EOF
    queue_type               = "memory"
    queue_check_point_writes = 0
    queue_max_bytes          = ""
    batch_delay              = 50
    batch_size               = 125
    workers                  = 1
  }
  op_type = 2
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Logstash instance id.
* `op_type` - (Required, Int) Operation type. 1: save only; 2: save and deploy.
* `pipeline` - (Required, List) Pipeline information.

The `pipeline` object supports the following:

* `batch_delay` - (Required, Int) Pipeline batch processing delay.
* `batch_size` - (Required, Int) Pipe batch size.
* `config` - (Required, String) Pipeline configuration content.
* `pipeline_desc` - (Required, String) Pipeline description information.
* `pipeline_id` - (Required, String) Pipeline id.
* `queue_check_point_writes` - (Required, Int) Number of pipeline buffer queue checkpoint writes.
* `queue_max_bytes` - (Required, String) Pipeline buffer queue size.
* `queue_type` - (Required, String) Pipeline buffer queue type.
* `workers` - (Required, Int) Number of Worker of pipe.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

elasticsearch logstash_pipeline can be imported using the id, e.g.

```
terraform import tencentcloud_elasticsearch_logstash_pipeline.logstash_pipeline ${instance_id}#${pipeline_id}
```

