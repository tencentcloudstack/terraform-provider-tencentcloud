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

* Create: [CreateCompareTask](https://cloud.tencent.com/document/api/571/82093) - 创建一致性校验任务
* Create: [DescribeCompareTasks](https://cloud.tencent.com/document/api/571/82088) - 查询一致性校验任务列表
* Create: [StartCompare](https://cloud.tencent.com/document/api/571/82074) - 启动一致性校验任务
* Read: [DescribeCompareTasks](https://cloud.tencent.com/document/api/571/82088) - 查询一致性校验任务列表
* Update: [ModifyCompareTask](https://cloud.tencent.com/document/api/571/82081) - 修改一致性校验任务
* Delete: [DeleteCompareTask](https://cloud.tencent.com/document/api/571/82090) - 删除一致性校验任务
* Delete: [DescribeCompareTasks](https://cloud.tencent.com/document/api/571/82088) - 查询一致性校验任务列表

