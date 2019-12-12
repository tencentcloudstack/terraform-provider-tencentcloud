---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcaplus_application"
sidebar_current: "docs-tencentcloud-resource-tcaplus_application"
description: |-
  Use this resource to create tcaplus application
---

# tencentcloud_tcaplus_application

Use this resource to create tcaplus application

## Example Usage



## Argument Reference

The following arguments are supported:

* `app_name` - (Required) Name of the tcapplus application. length should between 1 and 30.
* `idl_type` - (Required, ForceNew) ID of the tcapplus application.Valid values are PROTO,TDR,MIX.
* `password` - (Required) Password of the tcapplus application. length should between 12 and 16,a-z and 0-9 and A-Z must contain.
* `subnet_id` - (Required, ForceNew) Subnet id of the tcapplus application.
* `vpc_id` - (Required, ForceNew) VPC id of the tcapplus application.
* `old_password_expire_last` - (Optional) Old password expected expiration seconds after change password,must >= 300.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `api_access_id` - access id of the tcapplus application.For TcaplusDB SDK connect.
* `api_access_ip` - access ip of the tcapplus application.For TcaplusDB SDK connect.
* `api_access_port` - access port of the tcapplus application.For TcaplusDB SDK connect.
* `create_time` - create time of the tcapplus application.
* `network_type` - network type of the tcapplus application.
* `old_password_expire_time` - This field will display the old password expiration time,if password_status is `unmodifiable` means the old password has not yet expired, otherwise `-`.
* `password_status` - password status of the tcapplus application.`unmodifiable` means:can not change password now,`modifiable` means:can change password now.


## Import

tcaplus application can be imported using the id, e.g.

```
$ terraform import tencentcloud_tcaplus_application.test 26655801
```

```

