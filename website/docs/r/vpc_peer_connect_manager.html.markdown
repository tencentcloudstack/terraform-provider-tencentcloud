---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_peer_connect_manager"
sidebar_current: "docs-tencentcloud-resource-vpc_peer_connect_manager"
description: |-
  Provides a resource to create a vpc peer_connect_manager
---

# tencentcloud_vpc_peer_connect_manager

Provides a resource to create a vpc peer_connect_manager

## Example Usage

```hcl
data "tencentcloud_user_info" "info" {}

locals {
  owner_uin = data.tencentcloud_user_info.info.owner_uin
}

resource "tencentcloud_vpc" "vpc" {
  name       = "tf-example-pcx"
  cidr_block = "10.0.0.0/16"
}
resource "tencentcloud_vpc" "des_vpc" {
  name       = "tf-example-pcx-des"
  cidr_block = "172.16.0.0/16"
}
resource "tencentcloud_vpc_peer_connect_manager" "peer_connect_manager" {
  source_vpc_id           = tencentcloud_vpc.vpc.id
  peering_connection_name = "example-iac"
  destination_vpc_id      = tencentcloud_vpc.des_vpc.id
  destination_uin         = local.owner_uin
  destination_region      = "ap-guangzhou"
}
```

## Argument Reference

The following arguments are supported:

* `destination_region` - (Required, String) Peer region.
* `destination_uin` - (Required, String) Peer user UIN.
* `destination_vpc_id` - (Required, String) The unique ID of the peer VPC.
* `peering_connection_name` - (Required, String) Peer connection name.
* `source_vpc_id` - (Required, String) The unique ID of the local VPC.
* `bandwidth` - (Optional, Int) Bandwidth upper limit, unit Mbps.
* `charge_type` - (Optional, String) Billing mode, daily peak value POSTPAID_BY_DAY_MAX, monthly value 95 POSTPAID_BY_MONTH_95.
* `qos_level` - (Optional, String) Service classification PT, AU, AG.
* `type` - (Optional, String) Interworking type, VPC_PEER interworking between VPCs; VPC_BM_PEER interworking between VPC and BM Network.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

vpc peer_connect_manager can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_peer_connect_manager.peer_connect_manager peer_connect_manager_id
```

