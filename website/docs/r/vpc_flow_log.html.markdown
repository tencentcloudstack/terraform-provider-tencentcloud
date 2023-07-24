---
subcategory: "Flow Logs(FL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_flow_log"
sidebar_current: "docs-tencentcloud-resource-vpc_flow_log"
description: |-
  Provides a resource to create a vpc flow_log
---

# tencentcloud_vpc_flow_log

Provides a resource to create a vpc flow_log

~> **NOTE:** The cloud server instance specifications that support stream log collection include: M6ce, M6p, SA3se, S4m, DA3, ITA3, I6t, I6, S5se, SA2, SK1, S4, S5, SN3ne, S3ne, S2ne, SA2a, S3ne, SW3a, SW3b, SW3ne, ITA3, IT5c, IT5, IT5c, IT3, I3, D3, DA2, D2, M6, MA2, M4, C6, IT3a, IT3b, IT3c, C4ne, CN3ne, C3ne, GI1, PNV4, GNV4v, GNV4, GT4, GI3X, GN7, GN7vw.

~> **NOTE:** The following models no longer support the collection of new stream logs, and the stock stream logs will no longer be reported for data from July 25, 2022: Standard models: S3, SA1, S2, S1;Memory type: M3, M2, M1;Calculation type: C4, CN3, C3, C2;Batch calculation type: BC1, BS1;HPCC: HCCIC5, HCCG5v.

## Example Usage

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
```

## Argument Reference

The following arguments are supported:

* `flow_log_name` - (Required, String) Specify flow log name.
* `resource_id` - (Required, String) Specify resource unique Id of `resource_type` configured.
* `resource_type` - (Required, String) Specify resource type. NOTE: Only support `NETWORKINTERFACE` for now. Values: `VPC`, `SUBNET`, `NETWORKINTERFACE`, `CCN`, `NAT`, `DCG`.
* `traffic_type` - (Required, String) Specify log traffic type, values: `ACCEPT`, `REJECT`, `ALL`.
* `cloud_log_id` - (Optional, String) Specify flow log storage id.
* `cloud_log_region` - (Optional, String) Specify flow log storage region, default using current.
* `flow_log_description` - (Optional, String) Specify flow Log description.
* `flow_log_storage` - (Optional, List) Specify consumer detail, required while `storage_type` is `ckafka`.
* `storage_type` - (Optional, String) Specify consumer type, values: `cls`, `ckafka`.
* `tags` - (Optional, Map) Tag description list.
* `vpc_id` - (Optional, String) Specify vpc Id, ignore while `resource_type` is `CCN` (unsupported) but required while other types.

The `flow_log_storage` object supports the following:

* `storage_id` - (Optional, String) Specify storage instance id, required while `storage_type` is `ckafka`.
* `storage_topic` - (Optional, String) Specify storage topic id, required while `storage_type` is `ckafka`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

vpc flow_log can be imported using the flow log Id combine vpc Id, e.g.

```
$ terraform import tencentcloud_vpc_flow_log.flow_log flow_log_id fl-xxxx1234#vpc-yyyy5678
```

