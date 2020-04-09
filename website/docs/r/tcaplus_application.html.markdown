---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcaplus_application"
sidebar_current: "docs-tencentcloud-resource-tcaplus_application"
description: |-
  Use this resource to create tcaplus application
---

# tencentcloud_tcaplus_application

Use this resource to create tcaplus application

~> **NOTE:** tcaplus now only supports the following regions:ap-shanghai,ap-hongkong,na-siliconvalley,ap-singapore,ap-seoul,ap-tokyo,eu-frankfurt

## Example Usage



## Argument Reference

The following arguments are supported:

* `app_name` - (Required) Name of the tcapplus application. length should between 1 and 30.
* `idl_type` - (Required, ForceNew) Idl type of the tcapplus application.Valid values are PROTO,TDR,MIX.
* `password` - (Required) Password of the tcapplus application. length should between 12 and 16,a-z and 0-9 and A-Z must contain.
* `subnet_id` - (Required, ForceNew) Subnet id of the tcapplus application.
* `vpc_id` - (Required, ForceNew) VPC id of the tcapplus application.
* `old_password_expire_last` - (Optional) Old password expected expiration seconds after change password,must >= 300.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `api_access_id` - Access id of the tcapplus application.For TcaplusDB SDK connect.
* `api_access_ip` - Access ip of the tcapplus application.For TcaplusDB SDK connect.
* `api_access_port` - Access port of the tcapplus application.For TcaplusDB SDK connect.
* `create_time` - Create time of the tcapplus application.
* `network_type` - Network type of the tcapplus application.
* `old_password_expire_time` - This field will display the old password expiration time,if password_status is `unmodifiable` means the old password has not yet expired, otherwise `-`.
* `password_status` - Password status of the tcapplus application.`unmodifiable` means:can not change password now,`modifiable` means:can change password now.


## Import

tcaplus application can be imported using the id, e.g.

```
$ terraform import tencentcloud_tcaplus_application.test 26655801
```

```

