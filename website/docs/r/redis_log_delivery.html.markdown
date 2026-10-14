---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_log_delivery"
sidebar_current: "docs-tencentcloud-resource-redis_log_delivery"
description: |-
  Provides a resource to create Redis instance log delivery land set its attributes.
---

# tencentcloud_redis_log_delivery

Provides a resource to create Redis instance log delivery land set its attributes.

~> **NOTE:** When you use an existing cls logset and topic to enable logging, there is no need to set parameters such
as `period`, `create_index`, `log_region`, etc.

## Example Usage

### Use cls logset and topic which existed

```hcl
resource "tencentcloud_redis_log_delivery" "example" {
  instance_id = "crs-dmjj8en7"
  logset_id   = "cc31d9d6-74c0-4888-8b2f-b8148c3bcc5c"
  topic_id    = "5c2333e9-0bab-41fd-9f75-c602b3f9545f"
}
```

### Use exist cls logset and create new topic

```hcl
resource "tencentcloud_redis_log_delivery" "example" {
  instance_id  = "crs-dmjj8en7"
  logset_id    = "cc31d9d6-74c0-4888-8b2f-b8148c3bcc5c"
  topic_name   = "tf-example"
  period       = 20
  create_index = true
}
```

### Create new cls logset and topic

```hcl
resource "tencentcloud_redis_log_delivery" "example" {
  instance_id      = "crs-dmjj8en7"
  log_region       = "ap-guangzhou"
  logset_name      = "tf-example132"
  topic_name       = "tf-example132"
  period           = 20
  create_index     = true
  is_delete_topic  = true
  is_delete_logset = true
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID.
* `create_index` - (Optional, Bool) Whether to create an index when creating a log topic.
* `is_delete_logset` - (Optional, Bool) Whether to delete the associated Logset when deleting the log delivery. Default is false.
* `is_delete_topic` - (Optional, Bool) Whether to delete the associated Topic when deleting the log delivery. Default is false.
* `log_region` - (Optional, String) The region where the log set is located; if not specified, the region where the instance is located will be used by default.
* `logset_id` - (Optional, String) The ID of the log set being delivered.
* `logset_name` - (Optional, String) Log set name. If LogsetId does not specify a specific log set ID, please configure this parameter to set the log set name, and the system will automatically create a new log set with the specified name.
* `period` - (Optional, Int) Log storage time, defaults to 30 days, with an optional range of 1-3600 days.
* `topic_id` - (Optional, String) The ID of the topic being delivered.
* `topic_name` - (Optional, String) Log topic name, required when TopicId is empty, a new log topic will be automatically created.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [ModifyInstanceLogDelivery](https://www.tencentcloud.com/document/product/239/66802) - This API is used to enable or disable the shipping of instance logs to CLS.
* Read: [DescribeInstanceLogDelivery](https://www.tencentcloud.com/document/product/239/66803) - This API is used to query the instance log shipping configuration.
* Read: [DescribeLogsets](https://www.tencentcloud.com/document/product/614/42778) - This API is used to get the list of logsets.
* Read: [DescribeTopics](https://www.tencentcloud.com/document/product/614/42783) - This API is used to obtain logs or metric topic lists and supports pagination.
* Update: [ModifyInstanceLogDelivery](https://www.tencentcloud.com/document/product/239/66802) - This API is used to enable or disable the shipping of instance logs to CLS.
* Delete: [DeleteLogset](https://www.tencentcloud.com/document/product/614/42779) - This API is used to delete a logset.
* Delete: [DeleteTopic](https://www.tencentcloud.com/document/product/614/42784) - This API is used to delete logs or metric topics.
* Delete: [ModifyInstanceLogDelivery](https://www.tencentcloud.com/document/product/239/66802) - This API is used to enable or disable the shipping of instance logs to CLS.

## Import

Redis log delivery can be imported using the instanceId, e.g.

```
terraform import tencentcloud_redis_log_delivery.example crs-dmjj8en7
```

