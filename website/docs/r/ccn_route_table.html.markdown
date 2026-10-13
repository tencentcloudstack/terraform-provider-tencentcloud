---
subcategory: "Cloud Connect Network(CCN)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ccn_route_table"
sidebar_current: "docs-tencentcloud-resource-ccn_route_table"
description: |-
  Provides a resource to create a CCN Route table.
---

# tencentcloud_ccn_route_table

Provides a resource to create a CCN Route table.

## Example Usage

```hcl
variable "region" {
  default = "ap-guangzhou"
}

variable "availability_zone" {
  default = "ap-guangzhou-4"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "172.16.0.0/16"
}

# create subnet
resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create ccn
resource "tencentcloud_ccn" "example" {
  name                 = "tf-example"
  description          = "desc."
  qos                  = "AG"
  charge_type          = "PREPAID"
  bandwidth_limit_type = "INTER_REGION_LIMIT"
  tags = {
    createBy = "terraform"
  }
}

# create ccn route table
resource "tencentcloud_ccn_route_table" "example" {
  ccn_id      = tencentcloud_ccn.example.id
  name        = "tf-example"
  description = "desc."
}
```

## Argument Reference

The following arguments are supported:

* `ccn_id` - (Required, String, ForceNew) CCN Instance ID.
* `description` - (Required, String) Description of CCN Route table.
* `name` - (Required, String) CCN Route table name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - create time.
* `is_default_table` - True: default routing table False: non default routing table.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateCcnRouteTables](https://cloud.tencent.com/document/api/215/108749) - 新建云联网路由表
* Read: [DescribeCcnRouteTables](https://cloud.tencent.com/document/api/215/108745) - 查询云联网路由表信息
* Update: [ModifyCcnRouteTables](https://cloud.tencent.com/document/api/215/108742) - 修改云联网路由表
* Delete: [DeleteCcnRouteTables](https://cloud.tencent.com/document/api/215/108748) - 删除云联网路由表

## Import

Ccn instance can be imported, e.g.

```
$ terraform import tencentcloud_ccn_route_table.example ccnrtb-r5hrr417
```

