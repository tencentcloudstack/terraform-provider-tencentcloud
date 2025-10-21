---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eip_address_transform"
sidebar_current: "docs-tencentcloud-resource-eip_address_transform"
description: |-
  Provides a resource to create a eip address_transform
---

# tencentcloud_eip_address_transform

Provides a resource to create a eip address_transform

## Example Usage

```hcl
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

# create cvm
resource "tencentcloud_instance" "example" {
  instance_name              = "tf_example"
  availability_zone          = "ap-guangzhou-6"
  image_id                   = "img-9qrfy1xt"
  instance_type              = "SA3.MEDIUM4"
  system_disk_type           = "CLOUD_HSSD"
  system_disk_size           = 100
  hostname                   = "example"
  project_id                 = 0
  vpc_id                     = tencentcloud_vpc.vpc.id
  subnet_id                  = tencentcloud_subnet.subnet.id
  allocate_public_ip         = true
  internet_max_bandwidth_out = 10

  data_disks {
    data_disk_type = "CLOUD_HSSD"
    data_disk_size = 50
    encrypt        = false
  }

  tags = {
    tagKey = "tagValue"
  }
}

resource "tencentcloud_eip_address_transform" "example" {
  instance_id = tencentcloud_instance.example.id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) the instance ID of a normal public network IP to be operated. eg:ins-23mk45jn.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



