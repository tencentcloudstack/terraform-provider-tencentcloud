---
subcategory: "TDMQ for CMQ"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcmq_subscribe"
sidebar_current: "docs-tencentcloud-resource-tcmq_subscribe"
description: |-
  Provides a resource to create a tcmq subscribe
---

# tencentcloud_tcmq_subscribe

Provides a resource to create a tcmq subscribe

## Example Usage

```hcl
resource "tencentcloud_tcmq_subscribe" "subscribe" {
  topic_name        = "topic_name"
  subscription_name = "subscription_name"
  protocol          = "http"
  endpoint          = "http://xxxxxx" ;
}
```

## Argument Reference

The following arguments are supported:

* `endpoint` - (Required, String) `Endpoint` for notification receipt, which is distinguished by `Protocol`. For `http`, `Endpoint` must begin with `http://` and `host` can be a domain name or IP. For `Queue`, enter `QueueName`. Note that currently the push service cannot push messages to a VPC; therefore, if a VPC domain name or address is entered for `Endpoint`, pushed messages will not be received. Currently, messages can be pushed only to the public network and classic network.
* `protocol` - (Required, String) ubscription protocol. Currently, two protocols are supported: `http` and `queue`. To use the `http` protocol, you need to build your own web server to receive messages. With the `queue` protocol, messages are automatically pushed to a CMQ queue and you can pull them concurrently.
* `subscription_name` - (Required, String) Subscription name, which must be unique in the same topic under the same account in the same region. It can contain up to 64 letters, digits, and hyphens and must begin with a letter.
* `topic_name` - (Required, String) Topic name, which must be unique in the same topic under the same account in the same region. It can contain up to 64 letters, digits, and hyphens and must begin with a letter.
* `binding_key` - (Optional, Set: [`String`]) The number of `BindingKey` cannot exceed 5, and the length of each `BindingKey` cannot exceed 64 bytes. This field indicates the filtering policy for subscribing to and receiving messages. Each `BindingKey` includes up to 15 dots (namely up to 16 segments).
* `filter_tags` - (Optional, Set: [`String`]) Message body tag (used for message filtering). The number of tags cannot exceed 5, and each tag can contain up to 16 characters. It is used in conjunction with the `MsgTag` parameter of `(Batch)PublishMessage`. Rules: 1. If `FilterTag` is not configured, no matter whether `MsgTag` is configured, the subscription will receive all messages published to the topic; 2. If the array of `FilterTag` values has a value, only when at least one of the values in the array also exists in the array of `MsgTag` values (i.e., `FilterTag` and `MsgTag` have an intersection) can the subscription receive messages published to the topic; 3. If the array of `FilterTag` values has a value, but `MsgTag` is not configured, then no message published to the topic will be received, which can be considered as a special case of rule 2 as `FilterTag` and `MsgTag` do not intersect in this case. The overall design idea of rules is based on the intention of the subscriber.
* `notify_content_format` - (Optional, String) Push content format. Valid values: 1. JSON; 2. SIMPLIFIED, i.e., the raw format. If `Protocol` is `queue`, this value must be `SIMPLIFIED`. If `Protocol` is `http`, both options are acceptable, and the default value is `JSON`.
* `notify_strategy` - (Optional, String) CMQ push server retry policy in case an error occurs while pushing a message to `Endpoint`. Valid values: 1. `BACKOFF_RETRY`: backoff retry, which is to retry at a fixed interval, discard the message after a certain number of retries, and continue to push the next message; 2. `EXPONENTIAL_DECAY_RETRY`: exponential decay retry, which is to retry at an exponentially increasing interval, such as 1s, 2s, 4s, 8s, and so on. As a message can be retained in a topic for one day, failed messages will be discarded at most after one day of retry. Default value: `EXPONENTIAL_DECAY_RETRY`.
* `tags` - (Optional, Map) Tag description list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tcmq subscribe can be imported using the id, e.g.

```
terraform import tencentcloud_tcmq_subscribe.subscribe subscribe_id
```

