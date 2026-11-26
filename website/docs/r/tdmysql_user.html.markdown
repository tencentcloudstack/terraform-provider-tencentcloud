---
subcategory: "TDSQL Boundless(tdmysql)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmysql_user"
sidebar_current: "docs-tencentcloud-resource-tdmysql_user"
description: |-
  Provides a resource to manage a TDSQL-C for MySQL user.
---

# tencentcloud_tdmysql_user

Provides a resource to manage a TDSQL-C for MySQL user.

## Example Usage

```hcl
resource "tencentcloud_tdmysql_user" "example" {
  instance_id = "tdsqlc-c1q7xi8l"
  user_name   = "tf-example-user"
  host        = "10.0.0.1"
  password    = "TF@example-12345"
  description = "created by terraform"

  timeouts {
    create = "15m"
    delete = "15m"
  }
}
```

## Argument Reference

The following arguments are supported:

* `host` - (Required, String) Host. An IP range ends with %, which means all IPs in the range are allowed.
* `instance_id` - (Required, String) Instance ID.
* `user_name` - (Required, String) User name.
* `description` - (Optional, String) User description.
* `password` - (Optional, String) Unencrypted password.
* `timeouts` - (Optional, Object) The timeouts block allows you to specify the timeout for the create/update/delete operations.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - Creation time.
* `id` - ID of the resource.
* `update_time` - Update time.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateUsers](https://www.tencentcloud.com/document/api/1292/80901)
  - [DescribeFlow](https://www.tencentcloud.com/document/api/1292/78260)
  - [DescribeUsers](https://www.tencentcloud.com/document/api/1292/80119)
* Read: [DescribeUsers](https://www.tencentcloud.com/document/api/1292/80119)
* Update: [DescribeUsers](https://www.tencentcloud.com/document/api/1292/80119)
* Delete:
  - [DeleteUsers](https://www.tencentcloud.com/document/api/1292/80900)
  - [DescribeFlow](https://www.tencentcloud.com/document/api/1292/78260)

## Import

TDSQL-C for MySQL user can be imported using the composite id `instance_id#user_name#host`, e.g.

```
terraform import tencentcloud_tdmysql_user.example tdsqlc-c1q7xi8l#tf-example-user#10.0.0.1
```

