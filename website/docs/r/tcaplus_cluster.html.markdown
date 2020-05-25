---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcaplus_cluster"
sidebar_current: "docs-tencentcloud-resource-tcaplus_cluster"
description: |-
  Use this resource to create tcaplus cluster
---

# tencentcloud_tcaplus_cluster

Use this resource to create tcaplus cluster

~> **NOTE:** tcaplus now only supports the following regions:ap-shanghai,ap-hongkong,na-siliconvalley,ap-singapore,ap-seoul,ap-tokyo,eu-frankfurt

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

* `cluster_name` - (Required) Name of the tcaplus cluster. length should between 1 and 30.
* `idl_type` - (Required, ForceNew) Idl type of the tcaplus cluster.Valid values are PROTO,TDR,MIX.
* `password` - (Required) Password of the tcaplus cluster. length should between 12 and 16,a-z and 0-9 and A-Z must contain.
* `subnet_id` - (Required, ForceNew) Subnet id of the tcaplus cluster.
* `vpc_id` - (Required, ForceNew) VPC id of the tcaplus cluster.
* `old_password_expire_last` - (Optional) Old password expected expiration seconds after change password,must >= 300.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `api_access_id` - Access id of the tcaplus cluster.For TcaplusDB SDK connect.
* `api_access_ip` - Access ip of the tcaplus cluster.For TcaplusDB SDK connect.
* `api_access_port` - Access port of the tcaplus cluster.For TcaplusDB SDK connect.
* `create_time` - Create time of the tcaplus cluster.
* `network_type` - Network type of the tcaplus cluster.
* `old_password_expire_time` - This field will display the old password expiration time,if password_status is `unmodifiable` means the old password has not yet expired, otherwise `-`.
* `password_status` - Password status of the tcaplus cluster.`unmodifiable` means:can not change password now,`modifiable` means:can change password now.


## Import

tcaplus cluster can be imported using the id, e.g.

```
$ terraform import tencentcloud_tcaplus_cluster.test 26655801
```

```

