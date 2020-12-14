---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_target_groups"
sidebar_current: "docs-tencentcloud-datasource-clb_target_groups"
description: |-
  Use this data source to query target group information.
---

# tencentcloud_clb_target_groups

Use this data source to query target group information.

## Example Usage

```hcl
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-rule-basic"
}

resource "tencentcloud_clb_listener" "listener_basic" {
  clb_id        = tencentcloud_clb_instance.clb_basic.id
  port          = 1
  protocol      = "HTTP"
  listener_name = "listener_basic"
}

resource "tencentcloud_clb_listener_rule" "rule_basic" {
  clb_id              = tencentcloud_clb_instance.clb_basic.id
  listener_id         = tencentcloud_clb_listener.listener_basic.listener_id
  domain              = "abc.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
  target_type         = "TARGETGROUP"
}

resource "tencentcloud_clb_target_group" "test" {
  target_group_name = "test-target-keep-1"
}

resource "tencentcloud_clb_target_group_attachment" "group" {
  clb_id          = tencentcloud_clb_instance.clb_basic.id
  listener_id     = tencentcloud_clb_listener.listener_basic.listener_id
  rule_id         = tencentcloud_clb_listener_rule.rule_basic.rule_id
  targrt_group_id = tencentcloud_clb_target_group.test.id
}

data "tencentcloud_clb_target_groups" "target_group_info_id" {
  target_group_id = tencentcloud_clb_target_group.test.id
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional) Used to save results.
* `target_group_id` - (Optional) ID of Target group. Mutually exclusive with `vpc_id` and `target_group_name`. `target_group_id` is preferred.
* `target_group_name` - (Optional) Name of target group. Mutually exclusive with `target_group_id`. `target_group_id` is preferred.
* `vpc_id` - (Optional) Target group VPC ID. Mutually exclusive with `target_group_id`. `target_group_id` is preferred.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - Target group info list.
  * `associated_rule_list` - List of associated rules.
    * `domain` - Forwarding rule domain.
    * `listener_id` - Listener ID.
    * `listener_name` - Listener name.
    * `listener_port` - Listener port.
    * `load_balancer_id` - Load balance ID.
    * `load_balancer_name` - Load balance name.
    * `location_id` - Forwarding rule ID.
    * `protocol` - Listener protocol type.
    * `url` - Forwarding rule URL.
  * `create_time` - Creation time of the target group.
  * `port` - Port of target group.
  * `target_group_id` - ID of Target group.
  * `target_group_instance_list` - List of backend servers bound to the target group.
    * `eni_id` - ID of Elastic Network Interface.
    * `instance_id` - ID of backend service.
    * `instance_name` - The instance name of the backend service.
    * `private_ip_addresses` - Intranet IP list of back-end services.
    * `public_ip_addresses` - List of external network IP of back-end services.
    * `registered_time` - The time the backend service was bound.
    * `server_port` - Port of backend service.
    * `server_type` - Type of backend service.
    * `weight` - Forwarding weight of back-end services.
  * `target_group_name` - Target group VPC ID.
  * `update_time` - Modification time of the target group.
  * `vpc_id` - Name of target group.


