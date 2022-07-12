---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_as_scaling_group"
sidebar_current: "docs-tencentcloud-resource-kubernetes_as_scaling_group"
description: |-
  Auto scaling group for kubernetes cluster (offlined).
---

# tencentcloud_kubernetes_as_scaling_group

Auto scaling group for kubernetes cluster (offlined).

~> **NOTE:**  This resource was offline no longer suppored.

## Example Usage



## Argument Reference

The following arguments are supported:

* `auto_scaling_config` - (Required, List, ForceNew) Auto scaling config parameters.
* `auto_scaling_group` - (Required, List) Auto scaling group parameters.
* `cluster_id` - (Required, String, ForceNew) ID of the cluster.
* `extra_args` - (Optional, List: [`String`], ForceNew) Custom parameter information related to the node.
* `labels` - (Optional, Map, ForceNew) Labels of kubernetes AS Group created nodes.
* `unschedulable` - (Optional, Int, ForceNew) Sets whether the joining node participates in the schedule. Default is '0'. Participate in scheduling.

The `auto_scaling_config` object supports the following:

* `configuration_name` - (Required, String, ForceNew) Name of a launch configuration.
* `instance_type` - (Required, String, ForceNew) Specified types of CVM instance.
* `data_disk` - (Optional, List, ForceNew) Configurations of data disk.
* `enhanced_monitor_service` - (Optional, Bool, ForceNew) To specify whether to enable cloud monitor service. Default is TRUE.
* `enhanced_security_service` - (Optional, Bool, ForceNew) To specify whether to enable cloud security service. Default is TRUE.
* `instance_tags` - (Optional, Map, ForceNew) A list of tags used to associate different resources.
* `internet_charge_type` - (Optional, String, ForceNew) Charge types for network traffic. Valid value: `BANDWIDTH_PREPAID`, `TRAFFIC_POSTPAID_BY_HOUR`, `TRAFFIC_POSTPAID_BY_HOUR` and `BANDWIDTH_PACKAGE`.
* `internet_max_bandwidth_out` - (Optional, Int) Max bandwidth of Internet access in Mbps. Default is `0`.
* `key_ids` - (Optional, List, ForceNew) ID list of keys.
* `password` - (Optional, String, ForceNew) Password to access.
* `project_id` - (Optional, Int, ForceNew) Specifys to which project the configuration belongs.
* `public_ip_assigned` - (Optional, Bool, ForceNew) Specify whether to assign an Internet IP address.
* `security_group_ids` - (Optional, List, ForceNew) Security groups to which a CVM instance belongs.
* `system_disk_size` - (Optional, Int, ForceNew) Volume of system disk in GB. Default is `50`.
* `system_disk_type` - (Optional, String, ForceNew) Type of a CVM disk. Valid value: `CLOUD_PREMIUM` and `CLOUD_SSD`. Default is `CLOUD_PREMIUM`.

The `auto_scaling_group` object supports the following:

* `max_size` - (Required, Int) Maximum number of CVM instances (0~2000).
* `min_size` - (Required, Int) Minimum number of CVM instances (0~2000).
* `scaling_group_name` - (Required, String, ForceNew) Name of a scaling group.
* `vpc_id` - (Required, String, ForceNew) ID of VPC network.
* `default_cooldown` - (Optional, Int, ForceNew) Default cooldown time in second, and default value is 300.
* `desired_capacity` - (Optional, Int, ForceNew) Desired volume of CVM instances, which is between max_size and min_size.
* `forward_balancer_ids` - (Optional, List, ForceNew) List of application load balancers, which can't be specified with load_balancer_ids together.
* `load_balancer_ids` - (Optional, List, ForceNew) ID list of traditional load balancers.
* `project_id` - (Optional, Int, ForceNew) Specifys to which project the scaling group belongs.
* `retry_policy` - (Optional, String, ForceNew) Available values for retry policies include `IMMEDIATE_RETRY` and `INCREMENTAL_INTERVALS`.
* `subnet_ids` - (Optional, List, ForceNew) ID list of subnet, and for VPC it is required.
* `tags` - (Optional, Map, ForceNew) Tags of a scaling group.
* `termination_policies` - (Optional, List, ForceNew) Available values for termination policies include `OLDEST_INSTANCE` and `NEWEST_INSTANCE`.
* `zones` - (Optional, List, ForceNew) List of available zones, for Basic network it is required.

The `data_disk` object supports the following:

* `disk_size` - (Optional, Int, ForceNew) Volume of disk in GB. Default is `0`.
* `disk_type` - (Optional, String, ForceNew) Types of disk. Valid value: `CLOUD_PREMIUM` and `CLOUD_SSD`.
* `snapshot_id` - (Optional, String, ForceNew) Data disk snapshot ID.

The `forward_balancer_ids` object supports the following:

* `listener_id` - (Required, String, ForceNew) Listener ID for application load balancers.
* `load_balancer_id` - (Required, String, ForceNew) ID of available load balancers.
* `target_attribute` - (Required, List, ForceNew) Attribute list of target rules.
* `rule_id` - (Optional, String, ForceNew) ID of forwarding rules.

The `target_attribute` object supports the following:

* `port` - (Required, Int, ForceNew) Port number.
* `weight` - (Required, Int, ForceNew) Weight.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



