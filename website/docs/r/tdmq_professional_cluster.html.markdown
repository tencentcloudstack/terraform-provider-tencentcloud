---
subcategory: "TDMQ for Pulsar(tpulsar)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_professional_cluster"
sidebar_current: "docs-tencentcloud-resource-tdmq_professional_cluster"
description: |-
  Provides a resource to create a tdmq professional_cluster
---

# tencentcloud_tdmq_professional_cluster

Provides a resource to create a tdmq professional_cluster

## Example Usage

### single-zone Professional Cluster

```hcl
resource "tencentcloud_tdmq_professional_cluster" "professional_cluster" {
  auto_renew_flag = 1
  cluster_name    = "single_zone_cluster"
  product_name    = "PULSAR.P1.MINI2"
  storage_size    = 600
  tags = {
    "createby" = "terrafrom"
  }
  zone_ids = [
    100004,
  ]

  vpc {
    subnet_id = "subnet-xxxx"
    vpc_id    = "vpc-xxxx"
  }
}
```

### Multi-zone Professional Cluster

```hcl
resource "tencentcloud_tdmq_professional_cluster" "professional_cluster" {
  auto_renew_flag = 1
  cluster_name    = "multi_zone_cluster"
  product_name    = "PULSAR.P1.MINI2"
  storage_size    = 200
  tags = {
    "key"  = "value1"
    "key2" = "value2"
  }
  zone_ids = [
    330001,
    330002,
    330003,
  ]

  vpc {
    subnet_id = "subnet-xxxx"
    vpc_id    = "vpc-xxxx"
  }
}
```

## Argument Reference

The following arguments are supported:

* `auto_renew_flag` - (Required, Int) Whether to turn on automatic monthly renewal. `1`: turn on, `0`: turn off.
* `cluster_name` - (Required, String) Name of cluster. It does not support Chinese characters and special characters except dashes and underscores and cannot exceed 64 characters.
* `product_name` - (Required, String) Cluster specification code. Reference[Professional Cluster Specifications](https://cloud.tencent.com/document/product/1179/83705).
* `storage_size` - (Required, Int) Storage specifications. Reference[Professional Cluster Specifications](https://cloud.tencent.com/document/product/1179/83705).
* `zone_ids` - (Required, Set: [`Int`]) Multi-AZ deployment select three Availability Zones, like: [200002,200003,200004]. Single availability zone deployment selects an availability zone, like [200002].
* `auto_voucher` - (Optional, Int, ForceNew) Whether to automatically select vouchers. `1`: Yes, `0`: No. Default is `0`.
* `tags` - (Optional, Map) Tag description list.
* `time_span` - (Optional, Int, ForceNew) Purchase duration, value range: 1~50. Default: 1.
* `vpc` - (Optional, List) Label of VPC network.

The `vpc` object supports the following:

* `subnet_id` - (Required, String) Id of Subnet.
* `vpc_id` - (Required, String) Id of VPC.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tdmq professional_cluster can be imported using the id, e.g.

```
terraform import tencentcloud_tdmq_professional_cluster.professional_cluster professional_cluster_id
```

