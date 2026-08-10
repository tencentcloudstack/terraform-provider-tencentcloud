---
subcategory: "TcaplusDB"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcaplus_cluster"
sidebar_current: "docs-tencentcloud-resource-tcaplus_cluster"
description: |-
  Use this resource to create TcaplusDB cluster.
---

# tencentcloud_tcaplus_cluster

Use this resource to create TcaplusDB cluster.

~> **NOTE:** TcaplusDB now only supports the following regions: `ap-shanghai,ap-hongkong,na-siliconvalley,ap-singapore,ap-seoul,ap-tokyo,eu-frankfurt, and na-ashburn`.

## Example Usage

### Create a new tcaplus cluster instance

```hcl
locals {
  vpc_id    = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  subnet_id = data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id
}

variable "availability_zone" {
  default = "ap-guangzhou-3"
}

data "tencentcloud_vpc_subnets" "vpc" {
  is_default        = true
  availability_zone = var.availability_zone
}

resource "tencentcloud_tcaplus_cluster" "example" {
  idl_type                 = "PROTO"
  cluster_name             = "tf_example_tcaplus_cluster"
  vpc_id                   = local.vpc_id
  subnet_id                = local.subnet_id
  password                 = "your_pw_123111"
  old_password_expire_last = 3600
}
```

### Create a dedicated tcaplus cluster instance with resource tags, server list and proxy list

```hcl
locals {
  vpc_id    = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  subnet_id = data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id
}

variable "availability_zone" {
  default = "ap-guangzhou-3"
}

data "tencentcloud_vpc_subnets" "vpc" {
  is_default        = true
  availability_zone = var.availability_zone
}

resource "tencentcloud_tcaplus_cluster" "dedicated_example" {
  idl_type     = "PROTO"
  cluster_name = "tf_example_dedicated_cluster"
  vpc_id       = local.vpc_id
  subnet_id    = local.subnet_id
  password     = "your_pw_123111"
  cluster_type = 2

  resource_tags {
    tag_key   = "env"
    tag_value = "prod"
  }

  resource_tags {
    tag_key   = "owner"
    tag_value = "terraform"
  }

  server_list {
    machine_type = "S5.LARGE8"
    machine_num  = 2
  }

  proxy_list {
    machine_type = "S5.LARGE8"
    machine_num  = 1
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_name` - (Required, String) Name of the TcaplusDB cluster. Name length should be between 1 and 30.
* `idl_type` - (Required, String, ForceNew) IDL type of the TcaplusDB cluster. Valid values: `PROTO` and `TDR`.
* `password` - (Required, String) Password of the TcaplusDB cluster. Password length should be between 12 and 16. The password must be a *mix* of uppercase letters (A-Z), lowercase *letters* (a-z) and *numbers* (0-9).
* `subnet_id` - (Required, String, ForceNew) Subnet id of the TcaplusDB cluster.
* `vpc_id` - (Required, String, ForceNew) VPC id of the TcaplusDB cluster.
* `cluster_type` - (Optional, Int) Cluster type of the TcaplusDB cluster. `1`: shared cluster, `2`: dedicated cluster. This parameter is only valid for CreateCluster API and cannot be modified once set.
* `old_password_expire_last` - (Optional, Int) Expiration time of old password after password update, unit: second.
* `proxy_list` - (Optional, List) Dedicated proxy machine list of the TcaplusDB cluster. Only valid when `cluster_type` is `2` (dedicated cluster). For creation, each element exposes `machine_type` and `machine_num`. This parameter is only valid for CreateCluster API and cannot be modified once set.
* `resource_tags` - (Optional, List) Resource tags of the TcaplusDB cluster. This parameter is only valid for CreateCluster API and cannot be modified once set. Note: This field is write-only and will not be refreshed on Read because the DescribeClusters API does not return cluster-level tags.
* `server_list` - (Optional, List) Dedicated server machine list of the TcaplusDB cluster. Only valid when `cluster_type` is `2` (dedicated cluster). For creation, each element exposes `machine_type` and `machine_num`. This parameter is only valid for CreateCluster API and cannot be modified once set.

The `proxy_list` object supports the following:

* `machine_num` - (Optional, Int) 
* `machine_type` - (Optional, String) 

The `resource_tags` object supports the following:

* `tag_key` - (Required, String) 
* `tag_value` - (Optional, String) 

The `server_list` object supports the following:

* `machine_num` - (Optional, Int) 
* `machine_type` - (Optional, String) 

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `api_access_id` - Access ID of the TcaplusDB cluster.For TcaplusDB SDK connect.
* `api_access_ip` - Access IP of the TcaplusDB cluster.For TcaplusDB SDK connect.
* `api_access_port` - Access port of the TcaplusDB cluster.For TcaplusDB SDK connect.
* `create_time` - Create time of the TcaplusDB cluster.
* `network_type` - Network type of the TcaplusDB cluster.
* `old_password_expire_time` - Expiration time of the old password. If `password_status` is `unmodifiable`, it means the old password has not yet expired.
* `password_status` - Password status of the TcaplusDB cluster. Valid values: `unmodifiable`, `modifiable`. `unmodifiable`. which means the password can not be changed in this moment; `modifiable`, which means the password can be changed in this moment.


## Import

tcaplus cluster can be imported using the id, e.g.

```
$ terraform import tencentcloud_tcaplus_cluster.example cluster_id
```

