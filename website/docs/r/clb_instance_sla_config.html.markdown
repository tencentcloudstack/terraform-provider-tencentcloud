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
* `sla_type` - (Required, String) This parameter is required to create LCU-supported instances. Values:`SLA`: Super Large 4. When you have activated Super Large models, `SLA` refers to Super Large 4; `clb.c2.medium`: Standard; `clb.c3.small`: Advanced 1; `clb.c3.medium`: Advanced 1; `clb.c4.small`: Super Large 1; `clb.c4.medium`: Super Large 2; `clb.c4.large`: Super Large 3; `clb.c4.xlarge`: Super Large 4. For more details, see [Instance Specifications](https://intl.cloud.tencent.com/document/product/214/84689?from_cn_redirect=1).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

clb instance_sla_config can be imported using the id, e.g.

```
terraform import tencentcloud_clb_instance_sla_config.instance_sla_config instance_id
```

