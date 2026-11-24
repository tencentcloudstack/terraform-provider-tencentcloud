---
subcategory: "Oceanus"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_oceanus_oceanus_cluster"
sidebar_current: "docs-tencentcloud-resource-oceanus_oceanus_cluster"
description: |-
  Provides a resource to create an Oceanus (Stream Compute) cluster.
---

# tencentcloud_oceanus_oceanus_cluster

Provides a resource to create an Oceanus (Stream Compute) cluster.

Cluster creation is asynchronous: the provider persists the returned ClusterId
immediately and then polls DescribeClusters until the cluster reports status 2
(running). Only two in-place updates are supported: CU resizing via
ScaleOceanusCluster and prepaid renewal via RenewOceanusCluster. Changing any
other attribute forces the cluster to be recreated.

Note that `region_id`, `zone_id`, `period`, `instance_charge_type`,
`cluster_type`, `renew_flag`, `flink_web_ui_network_access_type`, `cu_memory`
and the VPC lists have no reliable read-back path in DescribeClusters, so drift
on those attributes cannot be detected and their values are preserved from
state. `login_password` is write-only and is never stored in state.

## Example Usage

```hcl
resource "tencentcloud_oceanus_oceanus_cluster" "example" {
  cluster_name       = "tf-example-cluster"
  region_id          = 1
  zone_id            = 100001
  login_password     = "Password@123"
  default_cos_bucket = "tf-example-bucket-1250000000"
  cu                 = 12
  remark             = "managed by terraform"

  vpc_descriptions {
    vpc_id    = "vpc-tf1234"
    subnet_id = "subnet-tf1234"
  }

  slave_vpc_descriptions {
    vpc_id    = "vpc-tf1234"
    subnet_id = "subnet-tf1234"
  }
}
```

### ### Resize an existing cluster

```hcl
resource "tencentcloud_oceanus_oceanus_cluster" "example" {
  cluster_name       = "tf-example-cluster"
  region_id          = 1
  zone_id            = 100001
  login_password     = "Password@123"
  default_cos_bucket = "tf-example-bucket-1250000000"

  # Changing cu triggers ScaleOceanusCluster in place.
  cu = 19

  vpc_descriptions {
    vpc_id    = "vpc-tf1234"
    subnet_id = "subnet-tf1234"
  }
}
```

### ### Renew a prepaid cluster

```hcl
resource "tencentcloud_oceanus_oceanus_cluster" "example" {
  cluster_name         = "tf-example-cluster"
  region_id            = 1
  zone_id              = 100001
  login_password       = "Password@123"
  default_cos_bucket   = "tf-example-bucket-1250000000"
  cu                   = 12
  instance_charge_type = "PREPAID"
  renew_flag           = "NOTIFY_AND_MANUAL_RENEW"

  # Changing period triggers RenewOceanusCluster in place.
  period = 12

  vpc_descriptions {
    vpc_id    = "vpc-tf1234"
    subnet_id = "subnet-tf1234"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_name` - (Required, String) Cluster name. Supported characters: 1-50 English letters, Chinese characters, digits, hyphens (-) or underscores (_). Changing the name forces recreation.
* `default_cos_bucket` - (Required, String) Default COS bucket used to store checkpoints, jar packages, or delivered logs. Changing it forces recreation.
* `login_password` - (Required, String) Flink UI access password. The cluster user name is admin. This value is write-only and never persisted in state. Changing it forces recreation.
* `region_id` - (Required, Int) Region ID, which can be queried via the DescribeRegions API with the Product parameter set to oceanus. Changing it forces recreation.
* `vpc_descriptions` - (Required, List) VPC used by the stream computing cluster to access other cloud resources in the same region via elastic network cards. Make sure the selected subnet has enough free IP addresses. Changing it forces recreation.
* `zone_id` - (Required, Int) Availability zone ID, which can be queried via the DescribeZones API with the Product parameter set to oceanus. Changing it forces recreation.
* `cluster_type` - (Optional, String) Cluster type. Valid values: `MULTI_AZ_CLUSTER` (multi availability zone cluster). Changing it forces recreation.
* `cu_memory` - (Optional, Int) CU to memory ratio. Valid values: 0, 2, 4, 8. Changing it forces recreation.
* `cu` - (Optional, Int) Number of cluster CUs. 12 CU is the minimal computing and billing unit; 1 CU contains 1 CPU and 4 GB memory. Management node fees are waived when the CU number is greater than or equal to 48. The cluster CU must satisfy 12 + 7*n (n >= 0). Resizing an existing cluster is applied in place via ScaleOceanusCluster.
* `flink_web_ui_network_access_type` - (Optional, String) Flink UI network access type. Valid values: `NetworkAccess_INTERNAL` (internal network), `NetworkAccess_EXTERNAL` (public network). Default value: `NetworkAccess_EXTERNAL`. Changing it forces recreation.
* `instance_charge_type` - (Optional, String) Cluster billing type. Valid values: `PREPAID` (monthly subscription), `POSTPAID_BY_SECOND` (pay as you go). Default value: `POSTPAID_BY_SECOND`. Changing it forces recreation.
* `period` - (Optional, Int) Purchase duration in months. Only effective for prepaid (PREPAID) clusters. Create accepts [1, 48]; renewal accepts [1, 36]. Changed in place via RenewOceanusCluster.
* `remark` - (Optional, String) Cluster description. Supported characters: 1-50 English letters, Chinese characters, digits, hyphens (-) or underscores (_). Changing it forces recreation.
* `renew_flag` - (Optional, String) Auto renewal flag, only valid for prepaid clusters. Valid values: `NOTIFY_AND_MANUAL_RENEW`, `NOTIFY_AND_AUTO_RENEW`, `DISABLE_NOTIFY_AND_MANUAL_RENEW`. Default value: `NOTIFY_AND_MANUAL_RENEW`. Changing it forces recreation.
* `slave_vpc_descriptions` - (Optional, List) VPC of the slave availability zone, only effective for multi availability zone clusters. Changing it forces recreation.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the cluster, in the format of cluster-xxx.


## Import

Oceanus cluster can be imported using the cluster ID, e.g.

```
$ terraform import tencentcloud_oceanus_oceanus_cluster.example cluster-tf123456
```

