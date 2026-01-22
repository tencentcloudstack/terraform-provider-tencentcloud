---
subcategory: "Gateway Load Balancer(GWLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gwlb_target_group"
sidebar_current: "docs-tencentcloud-resource-gwlb_target_group"
description: |-
  Provides a resource to create a gwlb gwlb_target_group
---

# tencentcloud_gwlb_target_group

Provides a resource to create a gwlb gwlb_target_group

## Example Usage

```hcl
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_gwlb_target_group" "gwlb_target_group" {
  target_group_name = "tf-test"
  vpc_id            = tencentcloud_vpc.vpc.id
  port              = 6081
  health_check {
    health_switch = true
    protocol      = "tcp"
    port          = 6081
    timeout       = 2
    interval_time = 5
    health_num    = 3
    un_health_num = 3
  }
}
```

## Argument Reference

The following arguments are supported:

* `all_dead_to_alive` - (Optional, Bool) Whether "All Dead, All Alive" is supported. It is supported by default.
* `health_check` - (Optional, List) Health check settings.
* `port` - (Optional, Int) Default port of the target group, which can be used when servers are added later. Either 'Port' or 'TargetGroupInstances.N.port' must be filled in.
* `protocol` - (Optional, String) GWLB target group protocol.
	- TENCENT_GENEVE: GENEVE standard protocol;
	- AWS_GENEVE: GENEVE compatibility protocol (a ticket is required for allowlisting).
* `schedule_algorithm` - (Optional, String) Load balancing algorithm.
	- IP_HASH_3_ELASTIC: elastic hashing.
* `target_group_name` - (Optional, String) Target group name, limited to 60 characters.
* `vpc_id` - (Optional, String) VPCID attribute of target group. If this parameter is left blank, the default VPC will be used.

The `health_check` object supports the following:

* `health_switch` - (Required, Bool) Whether to enable the health check.
* `health_num` - (Optional, Int) Health detection threshold. The default is 3 times. Value range: 2-10 times.
* `interval_time` - (Optional, Int) Detection interval time. The default is 5 seconds. Value range: 2-300 seconds.
* `port` - (Optional, Int) Health check port, which is required when the probe protocol is TCP.
* `protocol` - (Optional, String) Protocol used for health check, which supports PING and TCP and is PING by default.
	- PING: icmp;
	- TCP: tcp.
* `timeout` - (Optional, Int) Health check timeout. The default is 2 seconds. Value range: 2-30 seconds.
* `un_health_num` - (Optional, Int) Unhealth detection threshold. The default is 3 times. Value range: 2-10 times.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `target_group_instances` - Real server bound to a target group.
  * `bind_ip` - Private network IP of target group instance.
  * `port` - Port of target group instance. Only 6081 is supported.
  * `weight` - Weight of target group instance. Only 0 or 16 is supported, and non-0 is uniformly treated as 16.


## Import

gwlb gwlb_target_group can be imported using the id, e.g.

```
terraform import tencentcloud_gwlb_target_group.gwlb_target_group gwlb_target_group_id
```

