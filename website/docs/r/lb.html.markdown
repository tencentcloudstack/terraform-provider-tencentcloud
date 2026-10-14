---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lb"
sidebar_current: "docs-tencentcloud-resource-lb"
description: |-
  Provides a Load Balancer resource.
---

# tencentcloud_lb

Provides a Load Balancer resource.

~> **NOTE:** It has been deprecated and replaced by `tencentcloud_clb_instance`.

## Example Usage

```hcl
resource "tencentcloud_lb" "classic" {
  type       = "OPEN"
  forward    = "APPLICATION"
  name       = "tf-test-classic"
  project_id = 0
}
```

## Argument Reference

The following arguments are supported:

* `type` - (Required, String, ForceNew) The network type of the LB. Valid value: 'OPEN', 'INTERNAL'.
* `forward` - (Optional, String, ForceNew) The type of the LB. Valid value: 'CLASSIC', 'APPLICATION'.
* `name` - (Optional, String) The name of the LB.
* `project_id` - (Optional, Int, ForceNew) The project id of the LB, unspecified or 0 stands for default project.
* `vpc_id` - (Optional, String, ForceNew) The VPC ID of the LB, unspecified or 0 stands for CVM basic network.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `status` - The status of the LB.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateLoadBalancer](https://www.tencentcloud.com/document/product/214/33841) - This API is used to create a Cloud Load Balancer instance (this interface supports only pay-as-you-go CLB instances. For annual/monthly subscription, proceed to purchase through the console).
* Read: [DescribeLoadBalancers](https://www.tencentcloud.com/document/product/214/33830) - This API is used to query the list of CLB instances in a region.
* Update: [ModifyLoadBalancerAttributes](https://www.tencentcloud.com/document/product/214/33817) - This API is used to modify the attributes of a CLB instance, such as name and cross-region attributes.
* Delete: [DeleteLoadBalancer](https://www.tencentcloud.com/document/product/214/33837) - This API is used to delete one or more specified CLB instances. After successful deletion, the listeners and forwarding rules under the CLB instance will be deleted together, and the backend service will be unbound.

