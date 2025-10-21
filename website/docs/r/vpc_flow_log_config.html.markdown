---
subcategory: "Flow Logs(FL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_flow_log_config"
sidebar_current: "docs-tencentcloud-resource-vpc_flow_log_config"
description: |-
  Provides a resource to create a vpc flow_log_config
---

# tencentcloud_vpc_flow_log_config

Provides a resource to create a vpc flow_log_config

## Example Usage

### If disable FlowLogs

```hcl
data "tencentcloud_availability_zones" "zones" {}

data "tencentcloud_images" "image" {
  image_type       = ["PUBLIC_IMAGE"]
  image_name_regex = "Final"
}

data "tencentcloud_instance_types" "instance_types" {
  filter {
    name   = "zone"
    values = [data.tencentcloud_availability_zones.zones.zones.0.name]
  }

  filter {
    name   = "instance-family"
    values = ["S5"]
  }

  cpu_core_count   = 2
  exclude_sold_out = true
}

resource "tencentcloud_cls_logset" "logset" {
  logset_name = "delogsetmo"
  tags = {
    "test" = "test"
  }
}

resource "tencentcloud_cls_topic" "topic" {
  topic_name           = "topic"
  logset_id            = tencentcloud_cls_logset.logset.id
  auto_split           = false
  max_split_partitions = 20
  partition_count      = 1
  period               = 10
  storage_type         = "hot"
  tags = {
    "test" = "test",
  }
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-flow-log-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones.zones.zones.0.name
  name              = "vpc-flow-log-subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_eni" "example" {
  name        = "vpc-flow-log-eni"
  vpc_id      = tencentcloud_vpc.vpc.id
  subnet_id   = tencentcloud_subnet.subnet.id
  description = "eni desc"
  ipv4_count  = 1
}

resource "tencentcloud_instance" "example" {
  instance_name            = "ci-test-eni-attach"
  availability_zone        = data.tencentcloud_availability_zones.zones.zones.0.name
  image_id                 = data.tencentcloud_images.image.images.0.image_id
  instance_type            = data.tencentcloud_instance_types.instance_types.instance_types.0.instance_type
  system_disk_type         = "CLOUD_PREMIUM"
  disable_security_service = true
  disable_monitor_service  = true
  vpc_id                   = tencentcloud_vpc.vpc.id
  subnet_id                = tencentcloud_subnet.subnet.id
}

resource "tencentcloud_eni_attachment" "example" {
  eni_id      = tencentcloud_eni.example.id
  instance_id = tencentcloud_instance.example.id
}

resource "tencentcloud_vpc_flow_log" "example" {
  flow_log_name        = "tf-example-vpc-flow-log"
  resource_type        = "NETWORKINTERFACE"
  resource_id          = tencentcloud_eni_attachment.example.eni_id
  traffic_type         = "ACCEPT"
  vpc_id               = tencentcloud_vpc.vpc.id
  flow_log_description = "this is a testing flow log"
  cloud_log_id         = tencentcloud_cls_topic.topic.id
  storage_type         = "cls"
  tags = {
    "testKey" = "testValue"
  }
}

resource "tencentcloud_vpc_flow_log_config" "config" {
  flow_log_id = tencentcloud_vpc_flow_log.example.id
  enable      = false
}
```

### If enable FlowLogs

```hcl
resource "tencentcloud_vpc_flow_log_config" "config" {
  flow_log_id = tencentcloud_vpc_flow_log.example.id
  enable      = true
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Required, Bool) If enable snapshot policy.
* `flow_log_id` - (Required, String) Flow log ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

vpc flow_log_config can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_flow_log_config.flow_log_config flow_log_id
```

