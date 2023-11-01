---
subcategory: "Elasticsearch Service(ES)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_elasticsearch_instance_logs"
sidebar_current: "docs-tencentcloud-datasource-elasticsearch_instance_logs"
description: |-
  Use this data source to query detailed information of es elasticsearch_instance_logs
---

# tencentcloud_elasticsearch_instance_logs

Use this data source to query detailed information of es elasticsearch_instance_logs

## Example Usage

```hcl
data "tencentcloud_elasticsearch_instance_logs" "elasticsearch_instance_logs" {
  instance_id = "es-xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance id.
* `end_time` - (Optional, String) End time. The format is YYYY-MM-DD HH:MM:SS, such as 2019-01-22 20:15:53.
* `log_type` - (Optional, Int) Log type. Log type, default is 1, Valid values:
- 1: master log
- 2: Search slow log
- 3: Index slow log
- 4: GC log.
* `order_by_type` - (Optional, Int) Order type. Time sort method. Default is 0, valid values:
- 0: descending;
- 1: ascending order.
* `result_output_file` - (Optional, String) Used to save results.
* `search_key` - (Optional, String) Search key. Support LUCENE syntax, such as level:WARN, ip:1.1.1.1, message:test-index, etc.
* `start_time` - (Optional, String) Start time. The format is YYYY-MM-DD HH:MM:SS, such as 2019-01-22 20:15:53.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_log_list` - List of log details.
  * `ip` - Cluster node ip.
  * `level` - Log level.
  * `message` - Log message.
  * `node_id` - Cluster node id.
  * `time` - Log time.


