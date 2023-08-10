---
subcategory: "EventBridge(EB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eb_event_target"
sidebar_current: "docs-tencentcloud-resource-eb_event_target"
description: |-
  Provides a resource to create a eb event_target
---

# tencentcloud_eb_event_target

Provides a resource to create a eb event_target

## Example Usage

### Create an event target of type scf

```hcl
variable "zone" {
  default = "ap-guangzhou"
}

variable "namespace" {
  default = "default"
}

variable "function" {
  default = "keep-1676351130"
}

variable "function_version" {
  default = "$LATEST"
}

data "tencentcloud_cam_users" "foo" {
}

resource "tencentcloud_eb_event_bus" "foo" {
  event_bus_name = "tf-event_bus"
  description    = "event bus desc"
  enable_store   = false
  save_days      = 1
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_eb_event_rule" "foo" {
  event_bus_id = tencentcloud_eb_event_bus.foo.id
  rule_name    = "tf-event_rule"
  description  = "event rule desc"
  enable       = true
  event_pattern = jsonencode(
    {
      source = "apigw.cloud.tencent"
      type = [
        "connector:apigw",
      ]
    }
  )
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_eb_event_target" "scf_target" {
  event_bus_id = tencentcloud_eb_event_bus.foo.id
  rule_id      = tencentcloud_eb_event_rule.foo.rule_id
  type         = "scf"

  target_description {
    resource_description = "qcs::scf:${var.zone}:uin/${data.tencentcloud_cam_users.foo.user_list.0.uin}:namespace/${var.namespace}/function/${var.function}/${var.function_version}"

    scf_params {
      batch_event_count     = 1
      batch_timeout         = 1
      enable_batch_delivery = true
    }
  }
}
```

### Create an event target of type ckafka

```hcl
variable "ckafka" {
  default = "ckafka-qzoeaqx8"
}

resource "tencentcloud_eb_event_target" "ckafka_target" {
  event_bus_id = tencentcloud_eb_event_bus.foo.id
  rule_id      = tencentcloud_eb_event_rule.foo.rule_id
  type         = "ckafka"

  target_description {
    resource_description = "qcs::scf:${var.zone}:uin/${data.tencentcloud_cam_users.foo.user_list.0.uin}:ckafkaId/uin/${data.tencentcloud_cam_users.foo.user_list.0.uin}/${var.ckafka}"

    ckafka_target_params {
      topic_name = "dasdasd"

      retry_policy {
        max_retry_attempts = 360
        retry_interval     = 60
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `event_bus_id` - (Required, String) event bus id.
* `rule_id` - (Required, String) event rule id.
* `target_description` - (Required, List) target description.
* `type` - (Required, String) target type.

The `ckafka_target_params` object supports the following:

* `retry_policy` - (Required, List) retry strategy.
* `topic_name` - (Required, String) The ckafka topic to deliver to.

The `es_target_params` object supports the following:

* `index_prefix` - (Required, String) index prefix.
* `index_suffix_mode` - (Required, String) DTS index configuration.
* `net_mode` - (Required, String) network connection type.
* `output_mode` - (Required, String) DTS event configuration.
* `rotation_interval` - (Required, String) es log rotation granularity.
* `index_template_type` - (Optional, String) es template type.

The `retry_policy` object supports the following:

* `max_retry_attempts` - (Required, Int) Maximum number of retries.
* `retry_interval` - (Required, Int) Retry Interval Unit: Seconds.

The `scf_params` object supports the following:

* `batch_event_count` - (Optional, Int) Maximum number of events for batch delivery.
* `batch_timeout` - (Optional, Int) Maximum waiting time for bulk delivery.
* `enable_batch_delivery` - (Optional, Bool) Enable batch delivery.

The `target_description` object supports the following:

* `resource_description` - (Required, String) QCS resource six-stage format, more references [resource six-stage format](https://cloud.tencent.com/document/product/598/10606).
* `ckafka_target_params` - (Optional, List) Ckafka parameters.
* `es_target_params` - (Optional, List) ElasticSearch parameters.
* `scf_params` - (Optional, List) cloud function parameters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

eb event_target can be imported using the id, e.g.

```
terraform import tencentcloud_eb_event_target.event_target event_target_id
```

