---
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

```hcl
resource "tencentcloud_tcaplus_cluster" "test" {
  idl_type                 = "PROTO"
  cluster_name             = "tf_tcaplus_cluster_test"
  vpc_id                   = "vpc-7k6gzox6"
  subnet_id                = "subnet-akwgvfa3"
  password                 = "1qaA2k1wgvfa3ZZZ"
  old_password_expire_last = 3600
}
```

## Argument Reference

The following arguments are supported:

* `cluster_name` - (Required) Name of the TcaplusDB cluster. Name length should be between 1 and 30.
* `idl_type` - (Required, ForceNew) IDL type of the TcaplusDB cluster. Valid values are PROTO and TDR.
* `password` - (Required) Password of the TcaplusDB cluster. Password length should be between 12 and 16. The password must be a *mix* of uppercase letters (A-Z), lowercase *letters* (a-z) and *numbers* (0-9).
* `subnet_id` - (Required, ForceNew) Subnet id of the TcaplusDB cluster.
* `vpc_id` - (Required, ForceNew) VPC id of the TcaplusDB cluster.
* `old_password_expire_last` - (Optional) Expiration time of old password after password update, unit: second.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `api_access_id` - Access id of the TcaplusDB cluster.For TcaplusDB SDK connect.
* `api_access_ip` - Access ip of the TcaplusDB cluster.For TcaplusDB SDK connect.
* `api_access_port` - Access port of the TcaplusDB cluster.For TcaplusDB SDK connect.
* `create_time` - Create time of the TcaplusDB cluster.
* `network_type` - Network type of the TcaplusDB cluster.
* `old_password_expire_time` - Expiration time of the old password. If `password_status` is `unmodifiable`, it means the old password has not yet expired.
* `password_status` - Password status of the TcaplusDB cluster. Valid values: `unmodifiable`, which means the password can not be changed in this moment; `modifiable`, which means the password can be changed in this moment.


## Import

tcaplus cluster can be imported using the id, e.g.

```
$ terraform import tencentcloud_tcaplus_cluster.test 26655801
```

