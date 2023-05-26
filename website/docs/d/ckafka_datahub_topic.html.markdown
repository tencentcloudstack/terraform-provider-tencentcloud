---
subcategory: "Cloud Kafka(ckafka)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ckafka_datahub_topic"
sidebar_current: "docs-tencentcloud-datasource-ckafka_datahub_topic"
description: |-
  Use this data source to query detailed information of ckafka datahub_topic
---

# tencentcloud_ckafka_datahub_topic

Use this data source to query detailed information of ckafka datahub_topic

## Example Usage

```hcl
data "tencentcloud_ckafka_datahub_topic" "datahub_topic" {
}
```

## Argument Reference

The following arguments are supported:

* `limit` - (Optional, Int) The maximum number of results returned this time, the default is 50, and the maximum value is 50.
* `offset` - (Optional, Int) The offset position of this query, the default is 0.
* `result_output_file` - (Optional, String) Used to save results.
* `search_word` - (Optional, String) query key word.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `topic_list` - Topic list.
  * `name` - name.
  * `note` - Remark.
  * `partition_num` - number of partitions.
  * `retention_ms` - Expiration.
  * `status` - Status, 1 in use, 2 in deletion.
  * `topic_id` - Topic Id.
  * `topic_name` - Topic name.


