---
subcategory: "Data Transmission Service(DTS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dts_compare_task"
sidebar_current: "docs-tencentcloud-resource-dts_compare_task"
description: |-
  Provides a resource to create a dts compare_task
---

# tencentcloud_dts_compare_task

Provides a resource to create a dts compare_task

## Example Usage

```hcl
resource "tencentcloud_dts_compare_task" "compare_task" {
  job_id      = ""
  task_name   = ""
  object_mode = ""
  objects {
    object_mode = ""
    object_items {
      db_name     = ""
      db_mode     = ""
      schema_name = ""
      table_mode  = ""
      tables {
        table_name = ""
      }
      view_mode = ""
      views {
        view_name = ""
      }
    }

  }
}
```

## Argument Reference

The following arguments are supported:

* `job_id` - (Required, String) job id.
* `object_mode` - (Optional, String) object mode.
* `objects` - (Optional, List) objects.
* `task_name` - (Optional, String) task name.

The `object_items` object of `objects` supports the following:

* `db_mode` - (Optional, String) database mode.
* `db_name` - (Optional, String) database name.
* `schema_name` - (Optional, String) schema name.
* `table_mode` - (Optional, String) table mode.
* `tables` - (Optional, List) table list.
* `view_mode` - (Optional, String) view mode.
* `views` - (Optional, List) view list.

The `objects` object supports the following:

* `object_mode` - (Required, String) object mode.
* `object_items` - (Optional, List) object items.

The `tables` object of `object_items` supports the following:

* `table_name` - (Optional, String) table name.

The `views` object of `object_items` supports the following:

* `view_name` - (Optional, String) view name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `compare_task_id` - compare task id.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateCompareTask](https://cloud.tencent.com/document/api/571/82093) - This API is used to create a data consistency check task. After the task is successfully created, its ID will be returned in the format of `dts-8yv4w2i1-cmp-37skmii9`, and you can call `StartCompare` to start it.
* Create: [DescribeCompareTasks](https://cloud.tencent.com/document/api/571/82088) - This API is used to query the list of data consistency check tasks under the current task.
* Create: [StartCompare](https://cloud.tencent.com/document/api/571/82074) - This API is used to start a data consistency check task after creating it by calling the `CreateCompareTask` API. After calling this API, you can call the `DescribeCompareTasks` API to query the latest task status.
* Read: [DescribeCompareTasks](https://cloud.tencent.com/document/api/571/82088) - This API is used to query the list of data consistency check tasks under the current task.
* Update: [ModifyCompareTask](https://cloud.tencent.com/document/api/571/82081) - This API is used to modify the parameters of a data consistency check task after it is created and before it starts.
* Delete: [DeleteCompareTask](https://cloud.tencent.com/document/api/571/82090) - This API is used to delete a data consistency check task, which can be called when the task status is `success`, `failed`, or `canceled`.
* Delete: [DescribeCompareTasks](https://cloud.tencent.com/document/api/571/82088) - This API is used to query the list of data consistency check tasks under the current task.

