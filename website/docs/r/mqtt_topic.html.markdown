---
subcategory: "TDMQ for MQTT(MQTT)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mqtt_topic"
sidebar_current: "docs-tencentcloud-resource-mqtt_topic"
description: |-
  Provides a resource to create a MQTT topic
---

# tencentcloud_mqtt_topic

Provides a resource to create a MQTT topic

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

// create mqtt instance
resource "tencentcloud_mqtt_instance" "example" {
  instance_type = "BASIC"
  name          = "tf-example"
  sku_code      = "basic_2k"
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

// create topic
resource "tencentcloud_mqtt_topic" "example" {
  instance_id = tencentcloud_mqtt_instance.example.id
  topic       = "tf-example"
  remark      = "Remark."
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) MQTT instance ID.
* `topic` - (Required, String) Topic.
* `remark` - (Optional, String) Remarks.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

MQTT topic can be imported using the id, e.g.

```
terraform import tencentcloud_mqtt_topic.example mqtt-emb2v5wk#tf-example
```

