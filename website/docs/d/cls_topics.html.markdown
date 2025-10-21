---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_topics"
sidebar_current: "docs-tencentcloud-datasource-cls_topics"
description: |-
  Use this data source to query detailed information of CLS topics
---

# tencentcloud_cls_topics

Use this data source to query detailed information of CLS topics

## Example Usage

### Query all topics

```hcl
data "tencentcloud_cls_topics" "example" {}
```

### Query topics by filters

```hcl
data "tencentcloud_cls_topics" "example" {
  filters {
    key    = "topicId"
    values = ["88babc9b-ab8f-4dd1-9b04-3e2925cf9c4c"]
  }

  filters {
    key    = "topicName"
    values = ["tf-example"]
  }

  filters {
    key    = "logsetId"
    values = ["3e8e0521-32db-4532-beeb-9beefa56d3ea"]
  }

  filters {
    key    = "storageType"
    values = ["hot"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `biz_type` - (Optional, Int) Topic type
- 0 (default): Log topic.
- 1: Metric topic.
* `filters` - (Optional, List) <li>topicName: Filter by **log topic name**. Fuzzy match is implemented by default. You can use the `PreciseSearch` parameter to set exact match. Type: String. Required. No. <br><li>logsetName: Filter by **logset name**. Fuzzy match is implemented by default. You can use the `PreciseSearch` parameter to set exact match. Type: String. Required: No. <br><li>topicId: Filter by **log topic ID**. Type: String. Required: No. <br><li>logsetId: Filter by **logset ID**. You can call `DescribeLogsets` to query the list of created logsets or log in to the console to view them. You can also call `CreateLogset` to create a logset. Type: String. Required: No. <br><li>tagKey: Filter by **tag key**. Type: String. Required: No. <br><li>tag:tagKey: Filter by **tag key-value pair**. The `tagKey` should be replaced with a specified tag key, such as `tag:exampleKey`. Type: String. Required: No. <br><li>storageType: Filter by **log topic storage type**. Valid values: `hot` (standard storage) and `cold` (IA storage). Type: String. Required: No. Each request can have up to 10 `Filters` and 100 `Filter.Values`.
* `precise_search` - (Optional, Int) Match mode for `Filters` fields.
- 0: Fuzzy match for `topicName` and `logsetName`. This is the default value.
- 1: Exact match for `topicName`.
- 2: Exact match for `logsetName`.
- 3: Exact match for `topicName` and `logsetName`.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `key` - (Required, String) Field to be filtered.
* `values` - (Required, Set) Value to be filtered.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `topics` - Log topic list.


