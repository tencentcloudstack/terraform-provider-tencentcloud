---
subcategory: "Oceanus"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_oceanus_cluster"
sidebar_current: "docs-tencentcloud-resource-oceanus_cluster"
description: |-
  Provides a resource to create an Oceanus (Stream Compute) cluster.
---

# tencentcloud_oceanus_cluster

Provides a resource to create an Oceanus (Stream Compute) cluster.

~> **NOTE:** `login_password` is write-only and is never stored in state.

## Example Usage

```hcl
resource "tencentcloud_oceanus_cluster" "example" {
  cluster_name                     = "tf-test4"
  cluster_type                     = ""
  cu                               = 12
  cu_memory                        = 4
  default_cos_bucket               = "keep-tf-test-1308726196"
  flink_web_ui_network_access_type = "NetworkAccess_EXTERNAL"
  instance_charge_type             = "PREPAID"
  login_password                   = "Password@123"
  period                           = 1
  region                           = "ap-guangzhou"
  zone                             = "ap-guangzhou-7"
  remark                           = "测试"
  renew_flag                       = "NOTIFY_AND_MANUAL_RENEW"
  vpc_descriptions = [{
    subnet_id = "subnet-4sp1mbni"
    vpc_id    = "vpc-2trcrynd"
  }]
}
```

## Argument Reference

The following arguments are supported:

* `cluster_name` - (Required, String) Cluster name. Supported characters: 1-50 English letters, Chinese characters, digits, hyphens (-) or underscores (_). Changing the name forces recreation.
* `default_cos_bucket` - (Required, String) Default COS bucket used to store checkpoints, jar packages, or delivered logs. Changing it forces recreation.
* `login_password` - (Required, String) Flink UI access password. The cluster user name is admin. This value is write-only and never persisted in state. Changing it forces recreation.
* `region` - (Required, String) Region name, e.g. `ap-guangzhou`. It is translated to the numeric region id expected by the cluster APIs via the region DescribeRegions API. Changing it forces recreation.
* `vpc_descriptions` - (Required, List) VPC used by the stream computing cluster to access other cloud resources in the same region via elastic network cards. Make sure the selected subnet has enough free IP addresses. Changing it forces recreation.
* `zone` - (Required, String) Availability zone name, e.g. `ap-guangzhou-7`. It is translated to the numeric zone id expected by the cluster APIs via the region DescribeZones API. Changing it forces recreation.
* `cluster_type` - (Optional, String) Cluster type. Valid values: `MULTI_AZ_CLUSTER` (multi availability zone cluster). Changing it forces recreation.
* `cu_memory` - (Optional, Int) CU to memory ratio. Valid values: 0, 2, 4, 8. Changing it forces recreation.
* `cu` - (Optional, Int) Number of cluster CUs. 12 CU is the minimal computing and billing unit; 1 CU contains 1 CPU and 4 GB memory. Management node fees are waived when the CU number is greater than or equal to 48. The cluster CU must satisfy 12 + 7*n (n >= 0). Resizing an existing cluster is applied in place via ScaleOceanusCluster.
* `flink_web_ui_network_access_type` - (Optional, String) Flink UI network access type. Valid values: `NetworkAccess_INTERNAL` (internal network), `NetworkAccess_EXTERNAL` (public network). Default value: `NetworkAccess_EXTERNAL`. Changing it forces recreation.
* `instance_charge_type` - (Optional, String) Cluster billing type. Valid values: `PREPAID` (monthly subscription), `POSTPAID_BY_SECOND` (pay as you go). Default value: `POSTPAID_BY_SECOND`. Changing it forces recreation.
* `period` - (Optional, Int) Purchase duration in months. Only effective for prepaid (PREPAID) clusters. Create accepts [1, 48]; renewal accepts [1, 36]. Changed in place via RenewOceanusCluster.
* `remark` - (Optional, String) Cluster description. Supported characters: 1-50 English letters, Chinese characters, digits, hyphens (-) or underscores (_). Changing it forces recreation.
* `renew_flag` - (Optional, String) Auto renewal flag, only valid for prepaid clusters. Valid values: `NOTIFY_AND_MANUAL_RENEW`, `NOTIFY_AND_AUTO_RENEW`, `DISABLE_NOTIFY_AND_MANUAL_RENEW`. Default value: `NOTIFY_AND_MANUAL_RENEW`. Changing it forces recreation.
* `slave_vpc_descriptions` - (Optional, List) VPC of the slave availability zone, only effective for multi availability zone clusters. Changing it forces recreation.
* `timeouts` - (Optional, Object) The `timeouts` block.

The `slave_vpc_descriptions` object supports the following:

* `subnet_id` - (Required, String) Subnet ID.
* `vpc_id` - (Required, String) VPC ID.

The `vpc_descriptions` object supports the following:

* `subnet_id` - (Required, String) Subnet ID.
* `vpc_id` - (Required, String) VPC ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the cluster, in the format of cluster-xxx.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateOceanusCluster](https://www.tencentcloud.com/search/?wd=CreateOceanusCluster&type=document&lang=en)
  - [DescribeClusters](https://www.tencentcloud.com/search/?wd=DescribeClusters&type=document&lang=en)
  - [DescribeRegions](https://www.tencentcloud.com/document/api/1271/71925)
  - [DescribeZones](https://www.tencentcloud.com/document/api/1271/71924)
* Read: [DescribeClusters](https://www.tencentcloud.com/search/?wd=DescribeClusters&type=document&lang=en)
* Update:
  - [DescribeClusters](https://www.tencentcloud.com/search/?wd=DescribeClusters&type=document&lang=en)
  - [RenewOceanusCluster](https://www.tencentcloud.com/search/?wd=RenewOceanusCluster&type=document&lang=en)
  - [ScaleOceanusCluster](https://www.tencentcloud.com/search/?wd=ScaleOceanusCluster&type=document&lang=en)
* Delete: [DeleteOceanusCluster](https://www.tencentcloud.com/search/?wd=DeleteOceanusCluster&type=document&lang=en)

## Import

Oceanus cluster can be imported using the cluster ID, e.g.

```
terraform import tencentcloud_oceanus_cluster.example cluster-tf123456
```

