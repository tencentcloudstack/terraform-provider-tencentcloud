---
subcategory: "Flow Logs(FL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_flow_log_config"
sidebar_current: "docs-tencentcloud-resource-vpc_flow_log_config"
description: |-
  Provides a resource to create a vpc flow_log_config
---

# tencentcloud_vpc_flow_log_config

Provides a resource to create a vpc flow_log_config

## Example Usage

```hcl
resource "tencentcloud_vpc_flow_log_config" "flow_log_config" {
  flow_log_id = "fl-geg2keoj"
  enable      = false
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Required, Bool) If enable snapshot policy.
* `flow_log_id` - (Required, String) Flow log ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

vpc flow_log_config can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_flow_log_config.flow_log_config flow_log_id
```

