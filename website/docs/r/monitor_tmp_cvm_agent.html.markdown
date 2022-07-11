---
subcategory: "Monitor"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_tmp_cvm_agent"
sidebar_current: "docs-tencentcloud-resource-monitor_tmp_cvm_agent"
description: |-
  Provides a resource to create a monitor tmpCvmAgent
---

# tencentcloud_monitor_tmp_cvm_agent

Provides a resource to create a monitor tmpCvmAgent

## Example Usage

```hcl
resource "tencentcloud_monitor_tmp_cvm_agent" "tmpCvmAgent" {
  instance_id = "prom-c89b3b3u"
  name        = "test"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance id.
* `name` - (Required, String) Agent name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

monitor tmpCvmAgent can be imported using the id, e.g.
```
$ terraform import tencentcloud_monitor_tmp_cvm_agent.tmpCvmAgent instanceId#agentName
```

