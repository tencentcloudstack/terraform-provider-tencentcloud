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

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateProCluster](https://cloud.tencent.com/document/api/1179/101433) - This api is used to create a professional cluster with prepayment via api calls.
* Create: [DescribePulsarProInstances](https://cloud.tencent.com/document/api/1179/90376) - This API is used to query the list of the purchased TDMQ for Pulsar pro instances.
* Create: [ModifyResourceTags](https://www.tencentcloud.com/document/product/651/33012) - This API is used to modify all tags associated with a resource.
* Read: [DescribePulsarProInstanceDetail](https://cloud.tencent.com/document/api/1179/90377) - This API is used to obtain the information of a TDMQ for Pulsar pro cluster instance.
* Read: [DescribePulsarProInstances](https://cloud.tencent.com/document/api/1179/90376) - This API is used to query the list of the purchased TDMQ for Pulsar pro instances.
* Read: [DescribeResourceTagsByResourceIds](https://www.tencentcloud.com/document/product/651/33022) - This API is used to query the tag key-value pairs associated with an existing resource.
* Update: [ModifyCluster](https://cloud.tencent.com/document/api/1179/52189) - This API is used to update a cluster.
* Update: [ModifyResourceTags](https://www.tencentcloud.com/document/product/651/33012) - This API is used to modify all tags associated with a resource.
* Delete: [DeleteProCluster](https://cloud.tencent.com/document/api/1179/101986) - This API is used to delete a professional cluster with prepayment via API calls.

## Import

TDMQ professional cluster can be imported using the id, e.g.

```
terraform import tencentcloud_tdmq_professional_cluster.example pulsar-x4r939zkwmm2
```

