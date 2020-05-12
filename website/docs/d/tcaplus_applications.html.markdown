---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcaplus_applications"
sidebar_current: "docs-tencentcloud-datasource-tcaplus_applications"
description: |-
  Use this data source to query tcaplus applications
---

# tencentcloud_tcaplus_applications

Use this data source to query tcaplus applications

## Example Usage

```hcl
data "tencentcloud_tcaplus_applications" "name" {
  app_name = "app"
}
data "tencentcloud_tcaplus_applications" "id" {
  app_id = tencentcloud_tcaplus_application.test.id
}
data "tencentcloud_tcaplus_applications" "idname" {
  app_id   = tencentcloud_tcaplus_application.test.id
  app_name = "app"
}
```

## Argument Reference

The following arguments are supported:

* `app_id` - (Optional) Id of the tcapplus application to be query.
* `app_name` - (Optional) Name of the tcapplus application to be query.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of tcaplus application. Each element contains the following attributes.
  * `api_access_id` - Access id of the tcapplus application.For TcaplusDB SDK connect.
  * `api_access_ip` - Access ip of the tcapplus application.For TcaplusDB SDK connect.
  * `api_access_port` - Access port of the tcapplus application.For TcaplusDB SDK connect.
  * `app_id` - Id of the tcapplus application.
  * `app_name` - Name of the tcapplus application.
  * `create_time` - Create time of the tcapplus application.
  * `idl_type` - Idl type of the tcapplus application.
  * `network_type` - Network type of the tcapplus application.
  * `old_password_expire_time` - This field will display the old password expiration time,if password_status is `unmodifiable` means the old password has not yet expired, otherwise `-`.
  * `password_status` - Password status of the tcapplus application.`unmodifiable` means:can not change password now,`modifiable` means:can change password now.
  * `password` - Password of the tcapplus application.
  * `subnet_id` - Subnet id of the tcapplus application.
  * `vpc_id` - VPC id of the tcapplus application.


