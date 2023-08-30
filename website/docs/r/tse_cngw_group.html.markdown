---
subcategory: "Tencent Cloud Service Engine(TSE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tse_cngw_group"
sidebar_current: "docs-tencentcloud-resource-tse_cngw_group"
description: |-
  Provides a resource to create a tse cngw_group
---

# tencentcloud_tse_cngw_group

Provides a resource to create a tse cngw_group

## Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-4"
}

resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_tse_vpc"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  name              = "tf_tse_subnet"
  cidr_block        = "10.0.1.0/24"
}

resource "tencentcloud_tse_cngw_gateway" "cngw_gateway" {
  description                = "terraform test1"
  enable_cls                 = true
  engine_region              = "ap-guangzhou"
  feature_version            = "STANDARD"
  gateway_version            = "2.5.1"
  ingress_class_name         = "tse-nginx-ingress"
  internet_max_bandwidth_out = 0
  name                       = "terraform-gateway1"
  trade_type                 = 0
  type                       = "kong"

  node_config {
    number        = 2
    specification = "1c2g"
  }

  vpc_config {
    subnet_id = tencentcloud_subnet.subnet.id
    vpc_id    = tencentcloud_vpc.vpc.id
  }

  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tse_cngw_group" "cngw_group" {
  description = "terraform desc"
  gateway_id  = tencentcloud_tse_cngw_gateway.cngw_gateway.id
  name        = "terraform-group"
  subnet_id   = tencentcloud_subnet.subnet.id

  node_config {
    number        = 2
    specification = "1c2g"
  }
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, String) gateway IDonly postpaid gateway supported.
* `name` - (Required, String) gateway group name.
* `node_config` - (Required, List) group node configration.
* `subnet_id` - (Required, String) subnet ID. Assign an IP address to the engine in the VPC subnet. Reference value:- subnet-ahde9me9.
* `description` - (Optional, String) description information of group.
* `internet_config` - (Optional, List) internet configration.
* `internet_max_bandwidth_out` - (Optional, Int) public network outbound traffic bandwidth,[1,2048]Mbps.

The `internet_config` object supports the following:

* `description` - (Optional, String) description of clb.
* `internet_address_version` - (Optional, String) internet type. Reference value:- IPV4 (default value)- IPV6.
* `internet_max_bandwidth_out` - (Optional, Int) public network bandwidth.
* `internet_pay_mode` - (Optional, String) trade type of internet. Reference value:- BANDWIDTH- TRAFFIC (default value).
* `master_zone_id` - (Optional, String) primary availability zone.
* `multi_zone_flag` - (Optional, Bool) Whether load balancing has multiple availability zones.
* `sla_type` - (Optional, String) specification type of clb. Default shared type when this parameter is empty. Reference value:- SLA LCU-supported.
* `slave_zone_id` - (Optional, String) alternate availability zone.

The `node_config` object supports the following:

* `number` - (Required, Int) group node number, 2-50.
* `specification` - (Required, String) group specification, 1c2g|2c4g|4c8g|8c16g.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



