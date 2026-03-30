Use this data source to query detailed information of TEO realtime log delivery tasks.

Example Usage - Basic Query

```hcl
data "tencentcloud_teo_realtime_log_delivery_tasks" "tasks" {
}
```

Example Usage - Query with Offset and Limit

```hcl
data "tencentcloud_teo_realtime_log_delivery_tasks" "tasks" {
  offset = 0
  limit  = 20
}
```

Example Usage - Filter by Zone ID

```hcl
data "tencentcloud_teo_realtime_log_delivery_tasks" "tasks" {
  filters {
    name   = "zone-id"
    values = ["zone-abc123"]
  }
}
```

Example Usage - Filter by Task ID

```hcl
data "tencentcloud_teo_realtime_log_delivery_tasks" "tasks" {
  filters {
    name   = "task-id"
    values = ["task-123"]
  }
}
```

Example Usage - Filter by Task Type

```hcl
data "tencentcloud_teo_realtime_log_delivery_tasks" "tasks" {
  filters {
    name   = "task-type"
    values = ["cls"]
  }
}
```

Example Usage - Filter by Task Name with Fuzzy Matching

```hcl
data "tencentcloud_teo_realtime_log_delivery_tasks" "tasks" {
  filters {
    name   = "task-name"
    values = ["my-task"]
    fuzzy  = true
  }
}
```

Example Usage - Order Results

```hcl
data "tencentcloud_teo_realtime_log_delivery_tasks" "tasks" {
  order     = "task-id"
  direction = "asc"
}
```

Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter criteria. Supported filter fields: zone-id, task-id, task-name, task-type.
  * `name` - (Required, String) Field to be filtered.
  * `values` - (Required, Set) Value of the filtered field.
  * `fuzzy` - (Optional, Bool) Whether to enable fuzzy query.
* `offset` - (Optional, Int) Offset of the query result. Default is 0.
* `limit` - (Optional, Int) Limit on the number of query results. Default is 20.
* `order` - (Optional, String) Sort field, e.g., task-id, task-name.
* `direction` - (Optional, String) Sort direction: asc or desc.
* `result_output_file` - (Optional, String) Used to save results.

Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `realtime_log_delivery_tasks` - List of realtime log delivery tasks.
  * `task_id` - Task ID.
  * `zone_id` - ID of the site.
  * `task_name` - The name of the real-time log delivery task.
  * `task_type` - The real-time log delivery task type. Values: `cls`, `custom_endpoint`, `s3`.
  * `delivery_status` - The delivery status of the task.
  * `log_type` - Data delivery type.
  * `area` - Data delivery area. Values: `mainland`, `overseas`.
  * `entity_list` - List of entities corresponding to real-time log delivery tasks.
  * `fields` - A list of preset fields for delivery.
  * `custom_fields` - The list of custom fields delivered.
    * `name` - Extract data from the specified location. Values: `ReqHeader`, `RspHeader`, `Cookie`.
    * `value` - The name of the parameter whose value needs to be extracted.
    * `enabled` - Whether to deliver this field.
  * `sample` - Sampling ratio (1-100).
