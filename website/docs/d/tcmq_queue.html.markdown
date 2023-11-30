---
subcategory: "TDMQ for CMQ(tcmq)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcmq_queue"
sidebar_current: "docs-tencentcloud-datasource-tcmq_queue"
description: |-
  Use this data source to query detailed information of tcmq queue
---

# tencentcloud_tcmq_queue

Use this data source to query detailed information of tcmq queue

## Example Usage

```hcl
data "tencentcloud_tcmq_queue" "queue" {
  queue_name = "queue_name"
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter. Currently, you can filter by tag. The tag name must be prefixed with `tag:`, such as `tag: owner`, `tag: environment`, or `tag: business`.
* `is_tag_filter` - (Optional, Bool) For filtering by tag, this parameter must be set to `true`.
* `limit` - (Optional, Int) The number of queues to be returned per page in case of paginated return. If this parameter is not passed in, 20 will be used by default. Maximum value: 50.
* `offset` - (Optional, Int) Starting position of a queue list to be returned on the current page in case of paginated return. If a value is entered, limit must be specified. If this parameter is left empty, 0 will be used by default.
* `queue_name_list` - (Optional, Set: [`String`]) Filter by CMQ queue name.
* `queue_name` - (Optional, String) Filter by QueueName.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Optional, String) Filter parameter name.
* `values` - (Optional, Set) Value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `queue_list` - Queue list.
  * `active_msg_num` - Total number of messages in `Active` status (i.e., unconsumed) in the queue, which is an approximate value.
  * `bps` - Bandwidth limit.
  * `create_time` - Queue creation time. A Unix timestamp accurate down to the millisecond will be returned.
  * `create_uin` - Creator `Uin`.
  * `dead_letter_policy` - Dead letter queue policy.
    * `dead_letter_queue` - Dead letter queue.
    * `max_receive_count` - Maximum number of receipts.
    * `max_time_to_live` - Maximum period in seconds before an unconsumed message expires, which is required if `Policy` is 1. Value range: 300-43200. This value should be smaller than `MsgRetentionSeconds` (maximum message retention period).
    * `policy` - Dead letter queue policy.
  * `dead_letter_source` - Dead letter queue.
    * `queue_id` - Message queue ID.
    * `queue_name` - Message queue name.
  * `delay_msg_num` - Number of delayed messages.
  * `inactive_msg_num` - Total number of messages in `Inactive` status (i.e., being consumed) in the queue, which is an approximate value.
  * `last_modify_time` - Time when the queue attribute is last modified. A Unix timestamp accurate down to the millisecond will be returned.
  * `max_delay_seconds` - Maximum retention period for inflight messages.
  * `max_msg_backlog_size` - Maximum size of heaped messages in bytes.
  * `max_msg_heap_num` - Maximum number of heaped messages. The value range is 1,000,000-10,000,000 during the beta test and can be 1,000,000-1,000,000,000 after the product is officially released. The default value is 10,000,000 during the beta test and will be 100,000,000 after the product is officially released.
  * `max_msg_size` - Maximum message length. Value range: 1024-65536 bytes (i.e., 1-64 KB). Default value: 65536.
  * `max_unacked_msg_num` - The maximum number of unacknowledged messages.
  * `min_msg_time` - Minimum unconsumed time of message in seconds.
  * `msg_retention_seconds` - The max period during which a message is retained before it is automatically acknowledged. Value range: 30-43,200 seconds (30 seconds to 12 hours). Default value: 3600 seconds (1 hour).
  * `namespace_name` - Namespace name.
  * `polling_wait_seconds` - Long polling wait time for message reception. Value range: 0-30 seconds. Default value: 0.
  * `qps` - Limit of the number of messages produced per second. The value for consumed messages is 1.1 times this value.
  * `queue_id` - Message queue ID.
  * `queue_name` - Message queue name.
  * `retention_size_in_mb` - Queue storage space configured for message rewind. Value range: 1,024-10,240 MB (if message rewind is enabled). The value `0` indicates that message rewind is not enabled.
  * `rewind_msg_num` - Number of retained messages which have been deleted by the `DelMsg` API but are still within their rewind time range.
  * `rewind_seconds` - Rewindable time of messages in the queue. Value range: 0-1,296,000s (if message rewind is enabled). The value `0` indicates that message rewind is not enabled.
  * `status` - Cluster status. `0`: creating; `1`: normal; `2`: terminating; `3`: deleted; `4`: isolated; `5`: creation failed; `6`: deletion failed.
  * `tags` - Associated tag.
    * `tag_key` - Value of the tag key.
    * `tag_value` - Value of the tag value.
  * `tenant_id` - Tenant ID.
  * `trace` - Message trace. true: enabled; false: not enabled.
  * `transaction_policy` - Transaction message policy.
    * `first_query_interval` - First lookback time.
    * `max_query_count` - Maximum number of queries.
  * `transaction` - 1: transaction queue; 0: general queue.
  * `visibility_timeout` - Message visibility timeout period. Value range: 1-43200 seconds (i.e., 12 hours). Default value: 30.


