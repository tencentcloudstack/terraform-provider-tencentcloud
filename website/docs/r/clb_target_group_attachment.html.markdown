---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_target_group_attachment"
sidebar_current: "docs-tencentcloud-resource-clb_target_group_attachment"
description: |-
  Provides a resource to create a CLB target group attachment is bound to the load balancing listener or forwarding rule.
---

# tencentcloud_clb_target_group_attachment

Provides a resource to create a CLB target group attachment is bound to the load balancing listener or forwarding rule.

## Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-6"
}

// create vpc
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "vpc"
}

// create subnet
resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  name              = "subnet"
  cidr_block        = "10.0.1.0/24"
  is_multicast      = false
}

// create clb instance
resource "tencentcloud_clb_instance" "example" {
  clb_name     = "tf-example"
  network_type = "INTERNAL"
  vpc_id       = tencentcloud_vpc.vpc.id
  subnet_id    = tencentcloud_subnet.subnet.id
}

// create clb listener
resource "tencentcloud_clb_listener" "example" {
  clb_id        = tencentcloud_clb_instance.example.id
  listener_name = "tf-example"
  port          = 8080
  protocol      = "HTTP"
}

// create clb listener rule
resource "tencentcloud_clb_listener_rule" "example" {
  clb_id              = tencentcloud_clb_instance.example.id
  listener_id         = tencentcloud_clb_listener.example.listener_id
  domain              = "example.com"
  url                 = "/"
  session_expire_time = 60
  scheduler           = "WRR"
  target_type         = "TARGETGROUP"
}

// create clb target group
resource "tencentcloud_clb_target_group" "example" {
  target_group_name = "tf-example"
  vpc_id            = tencentcloud_vpc.vpc.id
}

// create clb target group attachment
resource "tencentcloud_clb_target_group_attachment" "example" {
  clb_id          = tencentcloud_clb_instance.example.id
  target_group_id = tencentcloud_clb_target_group.example.id
  listener_id     = tencentcloud_clb_listener.example.listener_id
  rule_id         = tencentcloud_clb_listener_rule.example.rule_id
}
```

## Argument Reference

The following arguments are supported:

* `clb_id` - (Required, String, ForceNew) ID of the CLB.
* `target_group_id` - (Required, String, ForceNew) ID of the CLB target group.
* `listener_id` - (Optional, String, ForceNew) ID of the CLB listener.
* `rule_id` - (Optional, String, ForceNew) ID of the CLB listener rule.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [AssociateTargetGroups](https://www.tencentcloud.com/document/product/214/33847) - This API is used to bind target groups to Cloud Load Balancer listeners (Layer-4 protocol) or forwarding rules (L7 protocol).
* Create: [DescribeLoadBalancers](https://www.tencentcloud.com/document/product/214/33830) - This API is used to query the list of CLB instances in a region.
* Create: [DescribeTargetGroups](https://www.tencentcloud.com/document/product/214/33826) - This API is used to query the target group information.
* Read: [DescribeTargetGroups](https://www.tencentcloud.com/document/product/214/33826) - This API is used to query the target group information.
* Delete: [DescribeTargetGroups](https://www.tencentcloud.com/document/product/214/33826) - This API is used to query the target group information.
* Delete: [DisassociateTargetGroups](https://www.tencentcloud.com/document/product/214/33822) - This API is used to disassociate a target group from a rule.

## Import

CLB target group attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_target_group_attachment.example lbtg-odareyb2#lbl-bicjmx3i#lb-cv0iz74c#loc-ac6uk7b6
```

