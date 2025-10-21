---
subcategory: "TDMQ for CMQ(tcmq)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcmq_subscribe"
sidebar_current: "docs-tencentcloud-datasource-tcmq_subscribe"
description: |-
  Use this data source to query detailed information of tcmq subscribe
---

# tencentcloud_tcmq_subscribe

Use this data source to query detailed information of tcmq subscribe

## Example Usage

```hcl
data "tencentcloud_tcmq_subscribe" "subscribe" {
  topic_name        = "topic_name"
  subscription_name = "subscription_name" ;
}
```

## Argument Reference

The following arguments are supported:

* `topic_name` - (Required, String) Topic name, which must be unique in the same topic under the same account in the same region. It can contain up to 64 letters, digits, and hyphens and must begin with a letter.
* `limit` - (Optional, Int) Number of topics to be returned per page in case of paginated return. If this parameter is not passed in, 20 will be used by default. Maximum value: 50.
* `offset` - (Optional, Int) Starting position of the list of topics to be returned on the current page in case of paginated return. If a value is entered, limit is required. If this parameter is left empty, 0 will be used by default.
* `result_output_file` - (Optional, String) Used to save results.
* `subscription_name` - (Optional, String) Fuzzy search by SubscriptionName.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `subscription_list` - Set of subscription attributes.
  * `binding_key` - Filtering policy for subscribing to and receiving messages.
  * `create_time` - Subscription creation time. A Unix timestamp accurate down to the millisecond will be returned.
  * `endpoint` - Endpoint that receives notifications, which varies by `protocol`: for HTTP, the endpoint must start with `http://`, and the `host` can be a domain or IP; for `queue`, `queueName` should be entered.
  * `filter_tags` - Filtering policy selected when a subscription is created:If `filterType` is 1, `filterTag` will be used for filtering. If `filterType` is 2, `bindingKey` will be used for filtering.
  * `last_modify_time` - Time when the subscription attribute is last modified. A Unix timestamp accurate down to the millisecond will be returned.
  * `msg_count` - Number of messages to be delivered in the subscription.
  * `notify_content_format` - Push content format. Valid values: 1. `JSON`; 2. `SIMPLIFIED`, i.e., the raw format. If `Protocol` is `queue`, this value must be `SIMPLIFIED`. If `Protocol` is `http`, both options are acceptable, and the default value is `JSON`.
  * `notify_strategy` - CMQ push server retry policy in case an error occurs while pushing a message to `Endpoint`. Valid values: 1. `BACKOFF_RETRY`: backoff retry, which is to retry at a fixed interval, discard the message after a certain number of retries, and continue to push the next message; 2. `EXPONENTIAL_DECAY_RETRY`: exponential decay retry, which is to retry at an exponentially increasing interval, such as 1s, 2s, 4s, 8s, and so on. As a message can be retained in a topic for one day, failed messages will be discarded at most after one day of retry. Default value: `EXPONENTIAL_DECAY_RETRY`.
  * `protocol` - Subscription protocol. Currently, two protocols are supported: HTTP and queue. To use the HTTP protocol, you need to build your own web server to receive messages. With the queue protocol, messages are automatically pushed to a CMQ queue and you can pull them concurrently.
  * `subscription_id` - Subscription ID, which will be used during monitoring data pull.
  * `subscription_name` - Subscription name, which must be unique in the same topic under the same account in the same region. It can contain up to 64 letters, digits, and hyphens and must begin with a letter.
  * `topic_owner` - Subscription owner APPID.


