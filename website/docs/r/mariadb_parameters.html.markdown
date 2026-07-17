---
subcategory: "TencentDB for MariaDB(MariaDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mariadb_parameters"
sidebar_current: "docs-tencentcloud-resource-mariadb_parameters"
description: |-
  Provides a resource to create a mariadb parameters
---

# tencentcloud_mariadb_parameters

Provides a resource to create a mariadb parameters

## Example Usage

```hcl
resource "tencentcloud_mariadb_parameters" "example" {
  instance_id = "tdsql-5n00ev40zl"
  params {
    param = "auto_increment_increment"
    value = "1"
  }
  params {
    param = "auto_increment_offset"
    value = "1"
  }
  params {
    param = "autocommit"
    value = "ON"
  }
  params {
    param = "binlog_transaction_dependency_history_size"
    value = "25000"
  }
  params {
    param = "binlog_write_threshold"
    value = "1610612736"
  }
  params {
    param = "character_set_server"
    value = "utf8mb4"
  }
  params {
    param = "collation_connection"
    value = "utf8mb4_general_ci"
  }
  params {
    param = "collation_database"
    value = "utf8mb4_general_ci"
  }
  params {
    param = "collation_server"
    value = "utf8mb4_general_ci"
  }
  params {
    param = "connect_timeout"
    value = "10"
  }
  params {
    param = "default_collation_for_utf8mb4"
    value = "utf8mb4_0900_ai_ci"
  }
  params {
    param = "default_week_format"
    value = "0"
  }
  params {
    param = "delay_key_write"
    value = "ON"
  }
  params {
    param = "delayed_insert_limit"
    value = "100"
  }
  params {
    param = "delayed_insert_timeout"
    value = "300"
  }
  params {
    param = "delayed_queue_size"
    value = "1000"
  }
  params {
    param = "div_precision_increment"
    value = "4"
  }
  params {
    param = "event_scheduler"
    value = "ON"
  }
  params {
    param = "explicit_defaults_for_timestamp"
    value = "ON"
  }
  params {
    param = "group_concat_max_len"
    value = "1024"
  }
  params {
    param = "innodb_autoinc_lock_mode"
    value = "2"
  }
  params {
    param = "innodb_backquery_enable"
    value = "OFF"
  }
  params {
    param = "innodb_backquery_window"
    value = "86400"
  }
  params {
    param = "innodb_concurrency_tickets"
    value = "5000"
  }
  params {
    param = "innodb_encryption_algorithm"
    value = "AES"
  }
  params {
    param = "innodb_flush_log_at_trx_commit"
    value = "1"
  }
  params {
    param = "innodb_lock_wait_timeout"
    value = "20"
  }
  params {
    param = "innodb_max_dirty_pages_pct"
    value = "70.000000"
  }
  params {
    param = "innodb_max_undo_log_size"
    value = "1073741824"
  }
  params {
    param = "innodb_old_blocks_pct"
    value = "37"
  }
  params {
    param = "innodb_old_blocks_time"
    value = "1000"
  }
  params {
    param = "innodb_purge_batch_size"
    value = "1000"
  }
  params {
    param = "innodb_read_ahead_threshold"
    value = "56"
  }
  params {
    param = "innodb_stats_method"
    value = "nulls_equal"
  }
  params {
    param = "innodb_stats_on_metadata"
    value = "OFF"
  }
  params {
    param = "innodb_strict_mode"
    value = "ON"
  }
  params {
    param = "innodb_table_locks"
    value = "ON"
  }
  params {
    param = "innodb_thread_concurrency"
    value = "0"
  }
  params {
    param = "interactive_timeout"
    value = "28800"
  }
  params {
    param = "join_buffer_size"
    value = "2097152"
  }
  params {
    param = "key_cache_age_threshold"
    value = "300"
  }
  params {
    param = "key_cache_block_size"
    value = "1024"
  }
  params {
    param = "key_cache_division_limit"
    value = "100"
  }
  params {
    param = "local_infile"
    value = "OFF"
  }
  params {
    param = "lock_wait_timeout"
    value = "5"
  }
  params {
    param = "log_queries_not_using_indexes"
    value = "OFF"
  }
  params {
    param = "long_query_time"
    value = "1.000000"
  }
  params {
    param = "low_priority_updates"
    value = "OFF"
  }
  params {
    param = "lower_case_table_names"
    value = "1"
  }
  params {
    param = "max_allowed_packet"
    value = "1073741824"
  }
  params {
    param = "max_binlog_size"
    value = "536870912"
  }
  params {
    param = "max_connect_errors"
    value = "2000"
  }
  params {
    param = "max_connections"
    value = "10000"
  }
  params {
    param = "max_execution_time"
    value = "0"
  }
  params {
    param = "max_prepared_stmt_count"
    value = "200000"
  }
  params {
    param = "myisam_sort_buffer_size"
    value = "4194304"
  }
  params {
    param = "net_buffer_length"
    value = "16384"
  }
  params {
    param = "net_read_timeout"
    value = "150"
  }
  params {
    param = "net_retry_count"
    value = "10"
  }
  params {
    param = "net_write_timeout"
    value = "300"
  }
  params {
    param = "optimizer_switch"
    value = "batched_key_access=off,block_nested_loop=on,condition_fanout_filter=on,csi_prefer_first_match_semi_join=on,csi_prefer_hash_group_by=on,csi_prefer_no_ref_access=on,csi_route_prefer=on,derived_condition_pushdown=on,derived_merge=on,duplicateweedout=on,engine_condition_pushdown=on,firstmatch=on,group_by_no_tmptable_for_csi=on,hash_join=on,hypergraph_optimizer=off,index_condition_pushdown=on,index_merge=on,index_merge_intersection=on,index_merge_sort_union=on,index_merge_union=on,loosescan=on,materialization=on,mrr=on,mrr_cost_based=on,prefer_ordering_index=on,semijoin=on,skip_scan=on,sort_merge_join=off,subquery_materialization_cost_based=on,subquery_to_derived=off,use_index_extensions=on,use_invisible_indexes=off,winmagic=off"
  }
  params {
    param = "performance_schema"
    value = "ON"
  }
  params {
    param = "query_alloc_block_size"
    value = "16384"
  }
  params {
    param = "query_prealloc_size"
    value = "24576"
  }
  params {
    param = "reject_table_no_pk"
    value = "1"
  }
  params {
    param = "slow_launch_time"
    value = "2"
  }
  params {
    param = "sort_buffer_size"
    value = "2097152"
  }
  params {
    param = "sql_mode"
    value = "NO_ENGINE_SUBSTITUTION,STRICT_TRANS_TABLES"
  }
  params {
    param = "sql_require_primary_key"
    value = "ON"
  }
  params {
    param = "sql_safe_updates"
    value = "OFF"
  }
  params {
    param = "sqlasyntimeout"
    value = "30"
  }
  params {
    param = "sync_binlog"
    value = "1"
  }
  params {
    param = "table_definition_cache"
    value = "10240"
  }
  params {
    param = "table_open_cache"
    value = "20480"
  }
  params {
    param = "thread_pool_oversubscribe"
    value = "30"
  }
  params {
    param = "thread_pool_size"
    value = "24"
  }
  params {
    param = "time_zone"
    value = "+08:00"
  }
  params {
    param = "tmp_table_size"
    value = "33554432"
  }
  params {
    param = "tx_isolation"
    value = "READ-COMMITTED"
  }
  params {
    param = "wait_timeout"
    value = "28800"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) instance id.
* `params` - (Required, Set) Number of days to keep, no more than 30.

The `params` object supports the following:

* `param` - (Required, String) parameter name.
* `value` - (Required, String) parameter value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mariadb parameters can be imported using the id, e.g.
```
terraform import tencentcloud_mariadb_parameters.example tdsql-4pzs5b67
```

