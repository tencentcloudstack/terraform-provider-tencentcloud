---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_security_group_attachment"
sidebar_current: "docs-tencentcloud-resource-clb_security_group_attachment"
description: |-
  Provides a resource to create a clb security_group_attachment
---

# tencentcloud_clb_security_group_attachment

Provides a resource to create a clb security_group_attachment

## Example Usage

```hcl
# create security group
resource "tencentcloud_security_group" "example" {
  name        = "tf-example"
  description = "sg desc."
  project_id  = 0

  tags = {
    "example" = "test"
  }
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = "ap-guangzhou-6"
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create clb
resource "tencentcloud_clb_instance" "example" {
  network_type = "INTERNAL"
  clb_name     = "tf-example"
  project_id   = 0
  vpc_id       = tencentcloud_vpc.vpc.id
  subnet_id    = tencentcloud_subnet.subnet.id

  tags = {
    "example" = "test"
  }
}

# attachment
resource "tencentcloud_clb_security_group_attachment" "example" {
  security_group    = tencentcloud_security_group.example.id
  load_balancer_ids = [tencentcloud_clb_instance.example.id]
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_ids` - (Required, Set: [`String`], ForceNew) Array of CLB instance IDs. Only support set one security group now.
* `security_group` - (Required, String, ForceNew) Security group ID, such as esg-12345678.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

clb security_group_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_clb_security_group_attachment.example sg-5275dorp#lb-5dnrkgry
```

