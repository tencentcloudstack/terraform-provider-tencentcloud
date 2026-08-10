---
subcategory: "TcaplusDB"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcaplus_cluster"
sidebar_current: "docs-tencentcloud-resource-tcaplus_cluster"
description: |-
  Provides a resource to create a TcaplusDB cluster.
---

# tencentcloud_tcaplus_cluster

Provides a resource to create a TcaplusDB cluster.

~> **NOTE:** TcaplusDB now only supports the following regions: `ap-shanghai`, `ap-hongkong`, `na-siliconvalley`, `ap-singapore`, `ap-seoul`, `ap-tokyo`, `eu-frankfurt`, `and na-ashburn`.

## Example Usage

### Create a tcaplus cluster instance with cluster_type is 1

```hcl
resource "tencentcloud_tcaplus_cluster" "example" {
  idl_type                 = "MIX"
  cluster_name             = "tf_example"
  vpc_id                   = "vpc-jll1dzwr"
  subnet_id                = "subnet-ef14ogeu"
  password                 = "Password@2026"
  old_password_expire_last = 3600
  cluster_type             = 1
  resource_tags {
    tag_key   = "createBy"
    tag_value = "Terraform"
  }
}
```

### Create a tcaplus cluster instance with cluster_type is 2

```hcl
resource "tencentcloud_tcaplus_cluster" "example" {
  idl_type                 = "MIX"
  cluster_name             = "tf_example"
  vpc_id                   = "vpc-qtzga3pm"
  subnet_id                = "subnet-c063n9el"
  password                 = "Password@2026"
  old_password_expire_last = 3600
  cluster_type             = 2
  server_list {
    machine_type = "T1"
    machine_num  = 4
  }

  proxy_list {
    machine_type = "T1"
    machine_num  = 2
  }

  resource_tags {
    tag_key   = "createBy"
    tag_value = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_name` - (Required, String) Cluster name, Chinese or English characters can be used, maximum length is 32 characters.
* `idl_type` - (Required, String, ForceNew) Cluster data description language type, uniformly filled with `MIX`, enumeration value: `MIX`: supports both `PROTO` and `TDR` tables.
* `password` - (Required, String) Cluster access password, must be `a-zA-Z0-9` characters, and must contain numbers, uppercase and lowercase letters.
* `subnet_id` - (Required, String, ForceNew) The subnet instance ID bound to the cluster, such as: `subnet-pxir56ns`.
* `vpc_id` - (Required, String, ForceNew) The private network instance ID bound to the cluster, such as: `vpc-f49l6u0z`.
* `cluster_type` - (Optional, Int) Cluster type: `1` shared, `2` dedicated.
* `old_password_expire_last` - (Optional, Int) Expiration time of old password after password update, unit: second.
* `proxy_list` - (Optional, Set) Dedicated cluster occupied proxy machines. Only valid when `cluster_type` is `2` (dedicated cluster). For creation, each element exposes `machine_type` and `machine_num`.
* `resource_tags` - (Optional, Set) Cluster tag set. Note: this field cannot be modified after cluster creation via CreateCluster, but can be modified via ModifyClusterTags. Tags will be refreshed on Read via DescribeClusterTags.
* `server_list` - (Optional, Set) Dedicated cluster occupied svr machines. Only valid when `cluster_type` is `2` (dedicated cluster). For creation, each element exposes `machine_type` and `machine_num`.

The `proxy_list` object supports the following:

* `machine_num` - (Optional, Int) Machine quantity.
* `machine_type` - (Optional, String) Machine type.

The `resource_tags` object supports the following:

* `tag_key` - (Required, String) Tag key.
* `tag_value` - (Optional, String) Tag value.

The `server_list` object supports the following:

* `machine_num` - (Optional, Int) Machine quantity.
* `machine_type` - (Optional, String) Machine type.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `api_access_id` - Access ID of the TcaplusDB cluster. For TcaplusDB SDK connect.
* `api_access_ip` - Access IP of the TcaplusDB cluster. For TcaplusDB SDK connect.
* `api_access_port` - Access port of the TcaplusDB cluster. For TcaplusDB SDK connect.
* `cluster_id` - Cluster ID.
* `create_time` - Create time of the TcaplusDB cluster.
* `network_type` - Network type of the TcaplusDB cluster.
* `old_password_expire_time` - Expiration time of the old password. If `password_status` is `unmodifiable`, it means the old password has not yet expired.
* `password_status` - Password status of the TcaplusDB cluster. Valid values: `unmodifiable`, `modifiable`. `unmodifiable`. which means the password can not be changed in this moment; `modifiable`, which means the password can be changed in this moment.


## Import

TcaplusDB cluster can be imported using the id, e.g.

```
terraform import tencentcloud_tcaplus_cluster.example 35402666774
```

