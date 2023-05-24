---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_renew_db_instance_operation"
sidebar_current: "docs-tencentcloud-resource-mysql_renew_db_instance_operation"
description: |-
  Provides a resource to create a mysql renew_db_instance_operation
---

# tencentcloud_mysql_renew_db_instance_operation

Provides a resource to create a mysql renew_db_instance_operation

## Example Usage

```hcl
resource "tencentcloud_mysql_renew_db_instance_operation" "renew_db_instance_operation" {
  instance_id     = "cdb-c1nl9rpv"
  time_span       = 1
  modify_pay_type = "PREPAID"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) The instance ID to be renewed, the format is: cdb-c1nl9rpv, which is the same as the instance ID displayed on the cloud database console page, you can use [Query Instance List](https://cloud.tencent.com/document/api/236/ 15872).
* `time_span` - (Required, Int, ForceNew) Renewal duration, unit: month, optional values include [1,2,3,4,5,6,7,8,9,10,11,12,24,36].
* `modify_pay_type` - (Optional, String, ForceNew) If you need to renew the Pay-As-You-Go instance to a Subscription instance, the value of this input parameter needs to be specified as `PREPAID`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `deadline_time` - Instance expiration time.
* `deal_id` - Deal id.


