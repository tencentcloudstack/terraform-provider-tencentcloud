---
subcategory: "TDMQ for Pulsar(tpulsar)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_professional_cluster"
sidebar_current: "docs-tencentcloud-resource-tdmq_professional_cluster"
description: |-
  Provides a resource to create a TDMQ professional cluster
---

# tencentcloud_tdmq_professional_cluster

Provides a resource to create a TDMQ professional cluster

## Example Usage

### Single-zone Professional Cluster

```hcl
resource "tencentcloud_tdmq_professional_cluster" "example" {
  auto_renew_flag  = 1
  cluster_name     = "tf_example"
  product_name     = "PULSAR.P2.SMALL4"
  storage_size     = 600
  instance_version = "2.9.2"
  zone_ids = [
    100006,
  ]

  vpc {
    vpc_id    = "vpc-i5yyodl9"
    subnet_id = "subnet-hhi88a58"
  }

  tags = {
    createby = "Terrafrom"
  }
}
```

### Multi-zone Professional Cluster

```hcl
resource "tencentcloud_tdmq_professional_cluster" "example" {
  auto_renew_flag  = 1
  cluster_name     = "tf_example"
  product_name     = "PULSAR.P2.SMALL4"
  storage_size     = 600
  instance_version = "3.0.0"
  zone_ids = [
    100006,
    100007
  ]

  vpc {
    vpc_id    = "vpc-i5yyodl9"
    subnet_id = "subnet-hhi88a58"
  }

  tags = {
    createby = "Terrafrom"
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
* `instance_version` - (Optional, String, ForceNew) Cluster version information. User can specify a version when creating the cluster.
* `tags` - (Optional, Map) Tag description list.
* `time_span` - (Optional, Int, ForceNew) Purchase duration, value range: 1~50. Default: 1.
* `vpc` - (Optional, List) Label of VPC network.

The `vpc` object supports the following:

* `subnet_id` - (Required, String) Id of Subnet.
* `vpc_id` - (Required, String) Id of VPC.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cluster_id` - Id of cluster.


## Import

TDMQ professional cluster can be imported using the id, e.g.

```
terraform import tencentcloud_tdmq_professional_cluster.example pulsar-x4r939zkwmm2
```

