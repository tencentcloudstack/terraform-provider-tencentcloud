---
subcategory: "TDMQ for MQTT(MQTT)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mqtt_instance"
sidebar_current: "docs-tencentcloud-resource-mqtt_instance"
description: |-
  Provides a resource to create a MQTT instance
---

# tencentcloud_mqtt_instance

Provides a resource to create a MQTT instance

## Example Usage

### Create a POSTPAID instance

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

// create mqtt instance
resource "tencentcloud_mqtt_instance" "example" {
  instance_type = "PRO"
  name          = "tf-example"
  sku_code      = "pro_6k_1"
  remark        = "remarks."
  vpc_list {
    vpc_id    = tencentcloud_vpc.vpc.id
    subnet_id = tencentcloud_subnet.subnet.id
  }
  pay_mode = 0
  tags = {
    createBy = "Terraform"
  }
}
```

### Create a PREPAID instance

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

// create mqtt instance
resource "tencentcloud_mqtt_instance" "example" {
  instance_type = "PRO"
  name          = "tf-example"
  sku_code      = "pro_10k_2"
  remark        = "remarks."
  vpc_list {
    vpc_id    = tencentcloud_vpc.vpc.id
    subnet_id = tencentcloud_subnet.subnet.id
  }
  pay_mode             = 1
  time_span            = 1
  renew_flag           = 1
  force_delete         = false
  automatic_activation = true
  authorization_policy = true
  tags = {
    createBy = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_type` - (Required, String) Instance type. PRO for Professional Edition; PLATINUM for Platinum Edition.
* `name` - (Required, String) Instance name.
* `sku_code` - (Required, String) Product SKU, available SKUs can be queried via the DescribeProductSKUList API.
* `authorization_policy` - (Optional, Bool) Authorization policy switch. Default is false.
* `automatic_activation` - (Optional, Bool) Is the automatic registration certificate automatically activated. Default is false.
* `force_delete` - (Optional, Bool) Indicate whether to force delete the instance. Default is `false`. If set true, the instance will be permanently deleted instead of being moved into the recycle bin. Note: only works for `PREPAID` instance.
* `pay_mode` - (Optional, Int) Payment mode (0: Postpaid; 1: Prepaid).
* `remark` - (Optional, String) Remarks.
* `renew_flag` - (Optional, Int) Whether to enable auto-renewal (0: Disabled; 1: Enabled).
* `tags` - (Optional, Map) Tags of the MQTT instance.
* `time_span` - (Optional, Int) Purchase duration (unit: months).
* `vpc_list` - (Optional, List) VPC information bound to the instance.

The `vpc_list` object supports the following:

* `subnet_id` - (Required, String) Subnet ID.
* `vpc_id` - (Required, String) VPC ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `device_certificate_provision_type` - Client certificate registration method: JITP: Automatic registration; API: Manually register through the API.


