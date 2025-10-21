---
subcategory: "Cloud Kafka(ckafka)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ckafka_renew_instance"
sidebar_current: "docs-tencentcloud-resource-ckafka_renew_instance"
description: |-
  Provides a resource to create a ckafka renew_instance
---

# tencentcloud_ckafka_renew_instance

Provides a resource to create a ckafka renew_instance

## Example Usage

```hcl
resource "tencentcloud_ckafka_renew_instance" "renew_ckafka_instance" {
  instance_id = "InstanceId"
  time_span   = 1
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) instance id.
* `time_span` - (Optional, Int, ForceNew) Renewal duration, the default is 1, and the unit is month.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



