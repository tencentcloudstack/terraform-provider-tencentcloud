---
subcategory: "Cloud Kafka(ckafka)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ckafka_group_info"
sidebar_current: "docs-tencentcloud-datasource-ckafka_group_info"
description: |-
  Use this data source to query detailed information of ckafka group_info
---

# tencentcloud_ckafka_group_info

Use this data source to query detailed information of ckafka group_info

## Example Usage

```hcl
data "tencentcloud_ckafka_group_info" "group_info" {
  instance_id = "ckafka-xxxxxx"
  group_list  = ["xxxxxx"]
}
```

## Argument Reference

The following arguments are supported:

* `group_list` - (Required, Set: [`String`]) Kafka consumption group, Consumer-group, here is an array format, format GroupList.0=xxx&amp;amp;GroupList.1=yyy.
* `instance_id` - (Required, String) InstanceId.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - result.
  * `error_code` - Error code, normally 0.
  * `group` - Kafka consumer group.
  * `members` - This array contains information only if state is Stable and protocol_type is consumer.
    * `assignment` - Stores the partition information assigned to the consumer.
      * `topics` - topic list.
        * `partitions` - Allocated partition information.
        * `topic` - Assigned topic name.
      * `version` - assignment version information.
    * `client_host` - Generally store the customer&#39;s IP address.
    * `client_id` - The client.id information set by the client consumer SDK itself.
    * `member_id` - ID that the coordinator generated for consumer.
  * `protocol_type` - The protocol type selected by the consumption group is normally the consumer, but some systems use their own protocol, such as kafka-connect, which uses connect. Only the standard consumer protocol, this interface knows the format of the specific allocation method, and can analyze the specific partition allocation.
  * `protocol` - Common consumer partition allocation algorithms are as follows (the default option for Kafka consumer SDK is range)  range|roundrobin|sticky.
  * `state` - Group state description (commonly Empty, Stable, and Dead states): Dead: The consumption group does not exist Empty: The consumption group does not currently have any consumer subscriptions PreparingRebalance: The consumption group is in the rebalance state CompletingRebalance: The consumption group is in the rebalance state Stable: Each consumer in the consumption group has joined and is in a stable state.


