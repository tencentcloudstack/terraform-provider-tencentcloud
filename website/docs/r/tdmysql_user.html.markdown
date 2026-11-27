---
subcategory: "TDSQL Boundless(tdmysql)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmysql_user"
sidebar_current: "docs-tencentcloud-resource-tdmysql_user"
description: |-
  Provides a resource to manage a TDSQL Boundless(tdmysql) user.
---

# tencentcloud_tdmysql_user

Provides a resource to manage a TDSQL Boundless(tdmysql) user.

## Example Usage

```hcl
resource "tencentcloud_tdmysql_user" "example" {
  instance_id = "tdsql3-ibpyogg7tg"
  user_name   = "tf-example"
  host        = "10.0.0.1"
  password    = "Password@2026"
  description = "created by terraform"

  timeouts {
    create = "3m"
    delete = "3m"
  }
}
```

## Argument Reference

The following arguments are supported:

* `host` - (Required, String) Host. An IP range ends with %, which means all IPs in the range are allowed.
* `instance_id` - (Required, String) Instance ID.
* `user_name` - (Required, String) User name.
* `description` - (Optional, String) User description.
* `password` - (Optional, String) Unencrypted password. It can be updated in-place; the change is applied through the ResetUsersPassword API.
* `timeouts` - (Optional, Object) The timeouts block allows you to specify the timeout for the create/update/delete operations.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateUsers](https://www.tencentcloud.com/document/api/1292/80901)
  - [DescribeFlow](https://www.tencentcloud.com/document/api/1292/78260)
  - [DescribeUsers](https://www.tencentcloud.com/document/api/1292/80119)
* Read: [DescribeUsers](https://www.tencentcloud.com/document/api/1292/80119)
* Update:
  - [DescribeFlow](https://www.tencentcloud.com/document/api/1292/78260)
  - [DescribeUsers](https://www.tencentcloud.com/document/api/1292/80119)
  - [ResetUsersPassword](https://www.tencentcloud.com/document/api/1292/81004)
* Delete:
  - [DeleteUsers](https://www.tencentcloud.com/document/api/1292/80900)
  - [DescribeFlow](https://www.tencentcloud.com/document/api/1292/78260)

## Import

TDSQL Boundless(tdmysql) user can be imported using the composite id `instance_id#user_name#host`, e.g.

```
terraform import tencentcloud_tdmysql_user.example tdsql3-ibpyogg7tg#tf-example#10.0.0.1
```

