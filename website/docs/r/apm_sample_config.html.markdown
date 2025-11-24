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
resource "tencentcloud_apm_sample_config" "example" {
  instance_id          = ""
  sample_name          = ""
  sample_rate          = 10
  service_name         = ""
  operation_name       = ""
  operation_type       = ""
  sample_config_status = 1
  tags {
    key   = "createdBy"
    value = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Business system ID.
* `sample_name` - (Required, String) Sampling rule name.
* `sample_rate` - (Required, Int) Sampling rate.
* `service_name` - (Required, String, ForceNew) Application name.
* `operation_name` - (Optional, String) API name.
* `operation_type` - (Optional, Int) 0: exact match (default); 1: prefix match; 2: suffix match.
* `sample_config_status` - (Optional, Int) Sample config status. 0: disabled; 1: enabled.
* `tags` - (Optional, List) Sampling tags.

The `tags` object supports the following:

* `key` - (Required, String) Key value definition.
* `value` - (Required, String) Value definition.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `sample_id` - Sample config ID.


## Import

APM sample config can be imported using the instanceId#sampleName, e.g.

```
terraform import tencentcloud_apm_sample_config.example apm-1o8yMC47u#tf-example
```

