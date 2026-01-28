Provides a list of ClickHouse (TCHouse-C) instances.

Example Usage

Query all instances

```hcl
data "tencentcloud_clickhouse_instances" "all" {
}
```

Query by instance ID

```hcl
data "tencentcloud_clickhouse_instances" "by_id" {
  instance_id = "cdwch-xxxxxx"
}
```

Query by instance name

```hcl
data "tencentcloud_clickhouse_instances" "by_name" {
  instance_name = "my-clickhouse-cluster"
}
```

Query by tags

```hcl
data "tencentcloud_clickhouse_instances" "by_tags" {
  tags = {
    env = "production"
    app = "analytics"
  }
}
```

Query with multiple filters

```hcl
data "tencentcloud_clickhouse_instances" "filtered" {
  instance_name = "test"
  tags = {
    env = "test"
  }
  is_simple = true
  result_output_file = "clickhouse_instances.json"
}

output "instance_count" {
  value = length(data.tencentcloud_clickhouse_instances.filtered.instance_list)
}

output "first_instance" {
  value = length(data.tencentcloud_clickhouse_instances.filtered.instance_list) > 0 ? data.tencentcloud_clickhouse_instances.filtered.instance_list[0] : null
}
```