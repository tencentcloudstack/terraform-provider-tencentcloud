---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eni_sg_attachment"
sidebar_current: "docs-tencentcloud-resource-eni_sg_attachment"
description: |-
  Provides a resource to create a eni_sg_attachment
---

# tencentcloud_eni_sg_attachment

Provides a resource to create a eni_sg_attachment

## Example Usage

```hcl
resource "tencentcloud_eni_sg_attachment" "eni_sg_attachment" {
  network_interface_ids = ["eni-p0hkgx8p"]
  security_group_ids    = ["sg-902tl7t7", "sg-edmur627"]
}
```

## Argument Reference

The following arguments are supported:

* `network_interface_ids` - (Required, Set: [`String`], ForceNew) ENI instance ID. Such as:eni-pxir56ns. It Only support set one eni instance now.
* `security_group_ids` - (Required, Set: [`String`], ForceNew) Security group instance ID, for example:sg-33ocnj9n, can be obtained through DescribeSecurityGroups. There is a limit of 100 instances per request.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

vpc eni_sg_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_eni_sg_attachment.eni_sg_attachment eni_sg_attachment_id
```

