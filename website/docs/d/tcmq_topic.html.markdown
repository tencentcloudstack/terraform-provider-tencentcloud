---
subcategory: "TDMQ for CMQ(tcmq)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcmq_topic"
sidebar_current: "docs-tencentcloud-datasource-tcmq_topic"
description: |-
  Use this data source to query detailed information of tcmq topic
---

# tencentcloud_tcmq_topic

Use this data source to query detailed information of tcmq topic

## Example Usage

```hcl
data "tencentcloud_tcmq_topic" "topic" {
  topic_name = "topic_name"
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter. Currently, you can filter by tag. The tag name must be prefixed with `tag:`, such as `tag: owner`, `tag: environment`, or `tag: business`.
* `is_tag_filter` - (Optional, Bool) For filtering by tag, this parameter must be set to `true`.
* `limit` - (Optional, Int) Number of topics to be returned per page in case of paginated return. If this parameter is not passed in, 20 will be used by default. Maximum value: 50.
* `offset` - (Optional, Int) Starting position of the list of topics to be returned on the current page in case of paginated return. If a value is entered, limit is required. If this parameter is left empty, 0 will be used by default.
* `result_output_file` - (Optional, String) Used to save results.
* `topic_name_list` - (Optional, Set: [`String`]) Filter by CMQ topic name.
* `topic_name` - (Optional, String) Fuzzy search by TopicName.

The `filters` object supports the following:

* `name` - (Optional, String) Filter parameter name.
* `values` - (Optional, Set) Value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `topic_list` - Topic list.
  * `broker_type` - Valid values: `0` (Pulsar), `1` (RocketMQ).
  * `create_time` - Topic creation time. A Unix timestamp accurate down to the millisecond will be returned.
  * `create_uin` - Creator `Uin`. The `resource` field for CAM authentication is composed of this field.
  * `filter_type` - Filtering policy selected when a subscription is created: If `filterType` is 1, `FilterTag` will be used for filtering. If `filterType` is 2, `BindingKey` will be used for filtering.
  * `last_modify_time` - Time when the topic attribute is last modified. A Unix timestamp accurate down to the millisecond will be returned.
  * `max_msg_size` - Maximum message size, which ranges from 1,024 to 1,048,576 bytes (i.e., 1-1,024 KB). The default value is 65,536.
  * `msg_count` - Number of current messages in the topic (number of retained messages).
  * `msg_retention_seconds` - Maximum lifecycle of message in topic. After the period specified by this parameter has elapsed since a message is sent to the topic, the message will be deleted no matter whether it has been successfully pushed to the user. This parameter is measured in seconds and defaulted to one day (86,400 seconds), which cannot be modified.
  * `namespace_name` - Namespace name.
  * `qps` - Number of messages published per second.
  * `status` - Cluster status. `0`: creating; `1`: normal; `2`: terminating; `3`: deleted; `4`: isolated; `5`: creation failed; `6`: deletion failed.
  * `tags` - Associated tag.
    * `tag_key` - Value of the tag key.
    * `tag_value` - Value of the tag value.
  * `tenant_id` - Tenant ID.
  * `topic_id` - Topic ID.
  * `topic_name` - Topic name.
  * `trace` - Message trace. true: enabled; false: not enabled.


