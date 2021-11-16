---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_instance"
sidebar_current: "docs-tencentcloud-resource-clb_instance"
description: |-
  Provides a resource to create a CLB instance.
---

# tencentcloud_clb_instance

Provides a resource to create a CLB instance.

## Example Usage

INTERNAL CLB

```hcl
resource "tencentcloud_clb_instance" "internal_clb" {
  network_type = "INTERNAL"
  clb_name     = "myclb"
  project_id   = 0
  vpc_id       = "vpc-7007ll7q"
  subnet_id    = "subnet-12rastkr"

  tags = {
    test = "tf"
  }
}
```

OPEN CLB

```hcl
resource "tencentcloud_clb_instance" "open_clb" {
  network_type              = "OPEN"
  clb_name                  = "myclb"
  project_id                = 0
  vpc_id                    = "vpc-da7ffa61"
  security_groups           = ["sg-o0ek7r93"]
  target_region_info_region = "ap-guangzhou"
  target_region_info_vpc_id = "vpc-da7ffa61"

  tags = {
    test = "tf"
  }
}
```

Default enable

```hcl
resource "tencentcloud_subnet" "subnet" {
  availability_zone = "ap-guangzhou-1"
  name              = "sdk-feature-test"
  vpc_id            = tencentcloud_vpc.foo.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

resource "tencentcloud_security_group" "sglab" {
  name        = "sg_o0ek7r93"
  description = "favourite sg"
  project_id  = 0
}

resource "tencentcloud_vpc" "foo" {
  name       = "for-my-open-clb"
  cidr_block = "10.0.0.0/16"

  tags = {
    "test" = "mytest"
  }
}

resource "tencentcloud_clb_instance" "open_clb" {
  network_type                 = "OPEN"
  clb_name                     = "my-open-clb"
  project_id                   = 0
  vpc_id                       = tencentcloud_vpc.foo.id
  load_balancer_pass_to_target = true

  security_groups           = [tencentcloud_security_group.sglab.id]
  target_region_info_region = "ap-guangzhou"
  target_region_info_vpc_id = tencentcloud_vpc.foo.id

  tags = {
    test = "open"
  }
}
```

CREATE multiple instance

```hcl
resource "tencentcloud_clb_instance" "open_clb1" {
  network_type   = "OPEN"
  clb_name       = "hello"
  master_zone_id = "ap-guangzhou-3"
}
~
```

CREATE instance with log

```hcl
resource "tencentcloud_vpc" "vpc_test" {
  name       = "clb-test"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_route_table" "rtb_test" {
  name   = "clb-test"
  vpc_id = "${tencentcloud_vpc.vpc_test.id}"
}

resource "tencentcloud_subnet" "subnet_test" {
  name              = "clb-test"
  cidr_block        = "10.0.1.0/24"
  availability_zone = "ap-guangzhou-3"
  vpc_id            = "${tencentcloud_vpc.vpc_test.id}"
  route_table_id    = "${tencentcloud_route_table.rtb_test.id}"
}

resource "tencentcloud_clb_log_set" "set" {
  period = 7
}

resource "tencentcloud_clb_log_topic" "topic" {
  log_set_id = "${tencentcloud_clb_log_set.set.id}"
  topic_name = "clb-topic"
}

resource "tencentcloud_clb_instance" "internal_clb" {
  network_type                 = "INTERNAL"
  clb_name                     = "myclb"
  project_id                   = 0
  vpc_id                       = "${tencentcloud_vpc.vpc_test.id}"
  subnet_id                    = "${tencentcloud_subnet.subnet_test.id}"
  load_balancer_pass_to_target = true
  log_set_id                   = "${tencentcloud_clb_log_set.set.id}"
  log_topic_id                 = "${tencentcloud_clb_log_topic.topic.id}"

  tags = {
    test = "tf"
  }
}
```

## Argument Reference

The following arguments are supported:

* `clb_name` - (Required) Name of the CLB. The name can only contain Chinese characters, English letters, numbers, underscore and hyphen '-'.
* `network_type` - (Required, ForceNew) Type of CLB instance. Valid values: `OPEN` and `INTERNAL`.
* `address_ip_version` - (Optional) IP version, only applicable to open CLB. Valid values are `ipv4`, `ipv6` and `IPv6FullChain`.
* `bandwidth_package_id` - (Optional) Bandwidth package id. If set, the `internet_charge_type` must be `BANDWIDTH_PACKAGE`.
* `internet_bandwidth_max_out` - (Optional) Max bandwidth out, only applicable to open CLB. Valid value ranges is [1, 2048]. Unit is MB.
* `internet_charge_type` - (Optional) Internet charge type, only applicable to open CLB. Valid values are `TRAFFIC_POSTPAID_BY_HOUR`, `BANDWIDTH_POSTPAID_BY_HOUR` and `BANDWIDTH_PACKAGE`.
* `load_balancer_pass_to_target` - (Optional) Whether the target allow flow come from clb. If value is true, only check security group of clb, or check both clb and backend instance security group.
* `log_set_id` - (Optional) The id of log set.
* `log_topic_id` - (Optional) the id of log topic.
* `master_zone_id` - (Optional) Setting master zone id of cross available zone disaster recovery, only applicable to open CLB.
* `project_id` - (Optional, ForceNew) ID of the project within the CLB instance, `0` - Default Project.
* `security_groups` - (Optional) Security groups of the CLB instance. Supports both `OPEN` and `INTERNAL` CLBs.
* `slave_zone_id` - (Optional) Setting slave zone id of cross available zone disaster recovery, only applicable to open CLB. this zone will undertake traffic when the master is down.
* `subnet_id` - (Optional, ForceNew) Subnet ID of the CLB. Effective only for CLB within the VPC. Only supports `INTERNAL` CLBs. Default is `ipv4`.
* `tags` - (Optional) The available tags within this CLB.
* `target_region_info_region` - (Optional) Region information of backend services are attached the CLB instance. Only supports `OPEN` CLBs.
* `target_region_info_vpc_id` - (Optional) Vpc information of backend services are attached the CLB instance. Only supports `OPEN` CLBs.
* `vpc_id` - (Optional, ForceNew) VPC ID of the CLB.
* `zone_id` - (Optional) Available zone id, only applicable to open CLB.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `clb_vips` - The virtual service address table of the CLB.
* `vip_isp` - Network operator, only applicable to open CLB. Valid values are `CMCC`(China Mobile), `CTCC`(Telecom), `CUCC`(China Unicom) and `BGP`. If this ISP is specified, network billing method can only use the bandwidth package billing (BANDWIDTH_PACKAGE).


## Import

CLB instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_instance.foo lb-7a0t6zqb
```

