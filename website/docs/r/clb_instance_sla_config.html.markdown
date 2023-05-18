---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_instance_sla_config"
sidebar_current: "docs-tencentcloud-resource-clb_instance_sla_config"
description: |-
  Provides a resource to create a clb instance_sla_config
---

# tencentcloud_clb_instance_sla_config

Provides a resource to create a clb instance_sla_config

## Example Usage

```hcl
resource "tencentcloud_clb_instance_sla_config" "instance_sla_config" {
  load_balancer_id = "lb-5dnrkgry"
  sla_type         = "SLA"
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required, String) ID of the CLB instance.
* `sla_type` - (Required, String) To upgrade to LCU-supported CLB instances. It must be SLA.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

clb instance_sla_config can be imported using the id, e.g.

```
terraform import tencentcloud_clb_instance_sla_config.instance_sla_config instance_id
```

