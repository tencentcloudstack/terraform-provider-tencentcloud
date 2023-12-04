---
subcategory: "Elasticsearch Service(ES)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_elasticsearch_describe_index_list"
sidebar_current: "docs-tencentcloud-datasource-elasticsearch_describe_index_list"
description: |-
  Use this data source to query detailed information of elasticsearch index list
---

# tencentcloud_elasticsearch_describe_index_list

Use this data source to query detailed information of elasticsearch index list

## Example Usage

```hcl
data "tencentcloud_elasticsearch_describe_index_list" "describe_index_list" {
  index_type  = "normal"
  instance_id = "es-nni6pm4s"
}
```

## Argument Reference

The following arguments are supported:

* `index_type` - (Required, String) Index type. `auto`: Autonomous index; `normal`: General index.
* `index_name` - (Optional, String) Index name. If you fill in the blanks, get all indexes.
* `index_status_list` - (Optional, Set: [`String`]) Index status list.
* `instance_id` - (Optional, String) ES cluster id.
* `order_by` - (Optional, String) Sort field. Support index name: IndexName, index storage: IndexStorage, index creation time: IndexCreateTime.
* `order` - (Optional, String) Sort order, which supports asc and desc. The default is desc data format asc,desc.
* `password` - (Optional, String) Cluster access password.
* `result_output_file` - (Optional, String) Used to save results.
* `username` - (Optional, String) Cluster access user name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `index_meta_fields` - Index metadata field.
  * `app_id` - App id.
  * `backing_indices` - Backing indices.
    * `index_create_time` - Index create time.
    * `index_name` - Index name.
    * `index_phrase` - Index phrase.
    * `index_status` - Index status.
    * `index_storage` - Index storage.
  * `cluster_id` - Cluster id.
  * `cluster_name` - Cluster name.
  * `cluster_version` - Cluster version.
  * `index_create_time` - Index create time.
  * `index_docs` - Number of indexed documents.
  * `index_meta_json` - Index meta json.
  * `index_name` - Index name.
  * `index_options_field` - Index options field.
    * `expire_max_age` - Expire max age.
    * `expire_max_size` - Expire max size.
    * `rollover_dynamic` - Whether to turn on dynamic scrolling.
    * `rollover_max_age` - Rollover max age.
    * `shard_num_dynamic` - Whether to enable dynamic slicing.
    * `timestamp_field` - Time partition field.
    * `write_mode` - Write mode.
  * `index_policy_field` - Index lifecycle field.
    * `cold_action` - Cold action.
    * `cold_enable` - Whether to enable the cold phase.
    * `cold_min_age` - Cold phase transition time.
    * `frozen_enable` - Start frozen phase.
    * `frozen_min_age` - Frozen phase transition time.
    * `warm_enable` - Whether to enable warm.
    * `warm_min_age` - Warm phase transition time.
  * `index_settings_field` - Index settings field.
    * `number_of_replicas` - Number of index copy fragments.
    * `number_of_shards` - Number of index main fragments.
    * `refresh_interval` - Index refresh frequency.
  * `index_status` - Index status.
  * `index_storage` - Index storage.
  * `index_type` - Index type.


