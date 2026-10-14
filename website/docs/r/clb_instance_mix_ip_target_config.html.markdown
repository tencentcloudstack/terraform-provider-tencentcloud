---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_instance_mix_ip_target_config"
sidebar_current: "docs-tencentcloud-resource-clb_instance_mix_ip_target_config"
description: |-
  Provides a resource to create a clb instance_mix_ip_target_config
---

# tencentcloud_clb_instance_mix_ip_target_config

Provides a resource to create a clb instance_mix_ip_target_config

## Example Usage

```hcl
resource "tencentcloud_clb_instance_mix_ip_target_config" "instance_mix_ip_target_config" {
  load_balancer_id = "lb-5dnrkgry"
  mix_ip_target    = false
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required, String, ForceNew) ID of CLB instances to be queried.
* `mix_ip_target` - (Required, Bool) False: closed True:open.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [ModifyLoadBalancerMixIpTarget](https://cloud.tencent.com/document/api/214/70576)
* Read: [DescribeLoadBalancers](https://cloud.tencent.com/document/api/214/30685) - This API is used to query the list of CLB instances in a region.
* Update: [ModifyLoadBalancerMixIpTarget](https://cloud.tencent.com/document/api/214/70576)

## Import

clb instance_mix_ip_target_config can be imported using the id, e.g.

```
terraform import tencentcloud_clb_instance_mix_ip_target_config.instance_mix_ip_target_config instance_id
```

