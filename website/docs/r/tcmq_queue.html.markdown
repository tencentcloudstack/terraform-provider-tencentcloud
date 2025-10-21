---
subcategory: "TDMQ for CMQ(tcmq)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcmq_queue"
sidebar_current: "docs-tencentcloud-resource-tcmq_queue"
description: |-
  Provides a resource to create a tcmq queue
---

# tencentcloud_tcmq_queue

Provides a resource to create a tcmq queue

## Example Usage

```hcl
resource "tencentcloud_tcmq_queue" "queue" {
  queue_name = "queue_name"
}
```

## Argument Reference

The following arguments are supported:

* `queue_name` - (Required, String) Queue name, which must be unique under the same account in the same region. It can contain up to 64 letters, digits, and hyphens and must begin with a letter.
* `dead_letter_queue_name` - (Optional, String) Dead letter queue name.
* `first_query_interval` - (Optional, Int) First lookback interval.
* `max_msg_heap_num` - (Optional, Int) Maximum number of heaped messages. The value range is 1,000,000-10,000,000 during the beta test and can be 1,000,000-1,000,000,000 after the product is officially released. The default value is 10,000,000 during the beta test and will be 100,000,000 after the product is officially released.
* `max_msg_size` - (Optional, Int) Maximum message length. Value range: 1024-65536 bytes (i.e., 1-64 KB). Default value: 65536.
* `max_query_count` - (Optional, Int) Maximum number of lookbacks.
* `max_receive_count` - (Optional, Int) Maximum receipt times. Value range: 1-1000.
* `max_time_to_live` - (Optional, Int) Maximum period in seconds before an unconsumed message expires, which is required if `policy` is 1. Value range: 300-43200. This value should be smaller than `msgRetentionSeconds` (maximum message retention period).
* `msg_retention_seconds` - (Optional, Int) The max period during which a message is retained before it is automatically acknowledged. Value range: 30-43,200 seconds (30 seconds to 12 hours). Default value: 3600 seconds (1 hour).
* `policy` - (Optional, Int) Dead letter policy. 0: message has been consumed multiple times but not deleted; 1: `Time-To-Live` has elapsed.
* `polling_wait_seconds` - (Optional, Int) Long polling wait time for message reception. Value range: 0-30 seconds. Default value: 0.
* `retention_size_in_mb` - (Optional, Int) Queue storage space configured for message rewind. Value range: 10,240-512,000 MB (if message rewind is enabled). The value `0` indicates that message rewind is not enabled.
* `rewind_seconds` - (Optional, Int) Rewindable time of messages in the queue. Value range: 0-1,296,000s (if message rewind is enabled). The value `0` indicates that message rewind is not enabled.
* `trace` - (Optional, Bool) Whether to enable message trace. true: yes; false: no. If this field is not configured, the feature will not be enabled.
* `transaction` - (Optional, Int) 1: transaction queue; 0: general queue.
* `visibility_timeout` - (Optional, Int) Message visibility timeout period. Value range: 1-43200 seconds (i.e., 12 hours). Default value: 30.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tcmq queue can be imported using the id, e.g.

```
terraform import tencentcloud_tcmq_queue.queue queue_id
```

