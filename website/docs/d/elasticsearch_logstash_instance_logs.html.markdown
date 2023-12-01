---
subcategory: "Elasticsearch Service(ES)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_elasticsearch_logstash_instance_logs"
sidebar_current: "docs-tencentcloud-datasource-elasticsearch_logstash_instance_logs"
description: |-
  Use this data source to query detailed information of elasticsearch logstash_instance_logs
---

# tencentcloud_elasticsearch_logstash_instance_logs

Use this data source to query detailed information of elasticsearch logstash_instance_logs

## Example Usage

```hcl
data "tencentcloud_elasticsearch_logstash_instance_logs" "logstash_instance_logs" {
  instance_id = "ls-xxxxxx"
  log_type    = 1
  start_time  = "2023-10-31 10:30:00"
  end_time    = "2023-10-31 10:30:10"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance id.
* `log_type` - (Required, Int) Log type. Default 1, Valid values:
 - 1: Main Log
 - 2: Slow log
 - 3: GC Log.
* `end_time` - (Optional, String) Log end time, in YYYY-MM-DD HH:MM:SS format, such as 2019-01-22 20:15:53.
* `order_by_type` - (Optional, Int) Time sort method. Default is 0. 0: descending; 1: ascending order.
* `result_output_file` - (Optional, String) Used to save results.
* `search_key` - (Optional, String) Search terms, support LUCENE syntax, such as level:WARN, ip:1.1.1.1, message:test-index, etc.
* `start_time` - (Optional, String) Log start time, in YYYY-MM-DD HH:MM:SS format, such as 2019-01-22 20:15:53.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_log_list` - List of log details.
  * `ip` - Cluster node ip.
  * `level` - Log level.
  * `message` - Log content.
  * `node_id` - Cluster node id.
  * `time` - Log time.


