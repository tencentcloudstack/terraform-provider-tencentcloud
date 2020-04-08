---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_as_scaling_group"
sidebar_current: "docs-tencentcloud-resource-kubernetes_as_scaling_group"
description: |-
  Provide a resource to create an auto scaling group for kubernetes cluster.
---

# tencentcloud_kubernetes_as_scaling_group

Provide a resource to create an auto scaling group for kubernetes cluster.

## Example Usage

```hcl
resource "tencentcloud_kubernetes_as_scaling_group" "test" {

  cluster_id = "cls-kb32pbv4"

  auto_scaling_group {
    scaling_group_name   = "tf-guagua-as-group"
    max_size             = "5"
    min_size             = "0"
    vpc_id               = "vpc-dk8zmwuf"
    subnet_ids           = ["subnet-pqfek0t8"]
    project_id           = 0
    default_cooldown     = 400
    desired_capacity     = "0"
    termination_policies = ["NEWEST_INSTANCE"]
    retry_policy         = "INCREMENTAL_INTERVALS"

    tags = {
      "test" = "test"
    }

  }

  auto_scaling_config {
    configuration_name = "tf-guagua-as-config"
    instance_type      = "SN3ne.8XLARGE64"
    project_id         = 0
    system_disk_type   = "CLOUD_PREMIUM"
    system_disk_size   = "50"

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 10
    public_ip_assigned         = true
    password                   = "test123#"
    enhanced_security_service  = false
    enhanced_monitor_service   = false

    instance_tags = {
      tag = "as"
    }

  }
}
```

## Argument Reference

The following arguments are supported:

* `auto_scaling_config` - (Required, ForceNew) Auto scaling config parameters.
* `auto_scaling_group` - (Required, ForceNew) Auto scaling group parameters.
* `cluster_id` - (Required, ForceNew) ID of the cluster.

The `auto_scaling_config` object supports the following:

* `configuration_name` - (Required, ForceNew) Name of a launch configuration.
* `instance_type` - (Required, ForceNew) Specified types of CVM instance.
* `data_disk` - (Optional, ForceNew) Configurations of data disk.
* `enhanced_monitor_service` - (Optional, ForceNew) To specify whether to enable cloud monitor service. Default is TRUE.
* `enhanced_security_service` - (Optional, ForceNew) To specify whether to enable cloud security service. Default is TRUE.
* `instance_tags` - (Optional, ForceNew) A list of tags used to associate different resources.
* `internet_charge_type` - (Optional, ForceNew) Charge types for network traffic. Available values include `BANDWIDTH_PREPAID`, `TRAFFIC_POSTPAID_BY_HOUR`, `TRAFFIC_POSTPAID_BY_HOUR` and `BANDWIDTH_PACKAGE`.
* `internet_max_bandwidth_out` - (Optional, ForceNew) Max bandwidth of Internet access in Mbps. Default is 0.
* `key_ids` - (Optional, ForceNew) ID list of keys.
* `password` - (Optional, ForceNew) Password to access.
* `project_id` - (Optional, ForceNew) Specifys to which project the configuration belongs.
* `public_ip_assigned` - (Optional, ForceNew) Specify whether to assign an Internet IP address.
* `security_group_ids` - (Optional, ForceNew) Security groups to which a CVM instance belongs.
* `system_disk_size` - (Optional, ForceNew) Volume of system disk in GB. Default is 50.
* `system_disk_type` - (Optional, ForceNew) Type of a CVM disk, and available values include CLOUD_PREMIUM and CLOUD_SSD. Default is CLOUD_PREMIUM.

The `auto_scaling_group` object supports the following:

* `max_size` - (Required, ForceNew) Maximum number of CVM instances (0~2000).
* `min_size` - (Required, ForceNew) Minimum number of CVM instances (0~2000).
* `scaling_group_name` - (Required, ForceNew) Name of a scaling group.
* `vpc_id` - (Required, ForceNew) ID of VPC network.
* `default_cooldown` - (Optional, ForceNew) Default cooldown time in second, and default value is 300.
* `desired_capacity` - (Optional, ForceNew) Desired volume of CVM instances, which is between max_size and min_size.
* `forward_balancer_ids` - (Optional, ForceNew) List of application load balancers, which can't be specified with load_balancer_ids together.
* `load_balancer_ids` - (Optional, ForceNew) ID list of traditional load balancers.
* `project_id` - (Optional, ForceNew) Specifys to which project the scaling group belongs.
* `retry_policy` - (Optional, ForceNew) Available values for retry policies include IMMEDIATE_RETRY and INCREMENTAL_INTERVALS.
* `subnet_ids` - (Optional, ForceNew) ID list of subnet, and for VPC it is required.
* `tags` - (Optional, ForceNew) Tags of a scaling group.
* `termination_policies` - (Optional, ForceNew) Available values for termination policies include OLDEST_INSTANCE and NEWEST_INSTANCE.
* `zones` - (Optional, ForceNew) List of available zones, for Basic network it is required.

The `data_disk` object supports the following:

* `disk_size` - (Optional, ForceNew) Volume of disk in GB. Default is 0.
* `disk_type` - (Optional, ForceNew) Types of disk, available values: CLOUD_PREMIUM and CLOUD_SSD.
* `snapshot_id` - (Optional, ForceNew) Data disk snapshot ID.

The `forward_balancer_ids` object supports the following:

* `listener_id` - (Required, ForceNew) Listener ID for application load balancers.
* `load_balancer_id` - (Required, ForceNew) ID of available load balancers.
* `target_attribute` - (Required, ForceNew) Attribute list of target rules.
* `rule_id` - (Optional, ForceNew) ID of forwarding rules.

The `target_attribute` object supports the following:

* `port` - (Required, ForceNew) Port number.
* `weight` - (Required, ForceNew) Weight.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



