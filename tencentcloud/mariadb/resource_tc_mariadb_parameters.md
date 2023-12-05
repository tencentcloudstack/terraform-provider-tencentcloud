Provides a resource to create a mariadb parameters

Example Usage

```hcl
resource "tencentcloud_mariadb_parameters" "parameters" {
    instance_id = "tdsql-4pzs5b67"

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
        value = "utf8mb4_general_ci"
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
        param = "group_concat_max_len"
        value = "1024"
    }
    params {
        param = "innodb_concurrency_tickets"
        value = "5000"
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
        value = "OFF"
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
        param = "query_alloc_block_size"
        value = "16384"
    }
    params {
        param = "query_prealloc_size"
        value = "24576"
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
Import

mariadb parameters can be imported using the id, e.g.
```
$ terraform import tencentcloud_mariadb_parameters.parameters tdsql-4pzs5b67
```