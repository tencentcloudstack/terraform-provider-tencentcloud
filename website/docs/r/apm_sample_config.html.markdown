---
subcategory: "Application Performance Management(APM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_apm_sample_config"
sidebar_current: "docs-tencentcloud-resource-apm_sample_config"
description: |-
  Provides a resource to create a APM sample config
---

# tencentcloud_apm_sample_config

Provides a resource to create a APM sample config

## Example Usage

```hcl
resource "tencentcloud_apm_instance" "example" {
  name                = "tf-example"
  description         = "desc."
  trace_duration      = 7
  span_daily_counters = 0
  tags = {
    createdBy = "Terraform"
  }
}

resource "tencentcloud_apm_sample_config" "example" {
  instance_id    = tencentcloud_apm_instance.example.id
  sample_name    = "tf-example"
  sample_rate    = 90
  service_name   = "java-order-serive"
  operation_type = 0
  tags {
    key   = "createdBy"
    value = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Business system ID.
* `sample_name` - (Required, String, ForceNew) Sampling rule name.
* `sample_rate` - (Required, Int) Sampling rate.
* `service_name` - (Required, String) Application name.
* `operation_name` - (Optional, String) API name.
* `operation_type` - (Optional, Int) 0: exact match (default); 1: prefix match; 2: suffix match.
* `tags` - (Optional, List) Sampling tags.

The `tags` object supports the following:

* `key` - (Required, String) Key value definition.
* `value` - (Required, String) Value definition.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

APM sample config can be imported using the instanceId#sampleName, e.g.

```
terraform import tencentcloud_apm_sample_config.example apm-jPr5iQL77#tf-example
```

