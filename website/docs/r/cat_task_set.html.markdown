---
subcategory: "Cloud Automated Testing(CAT)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cat_task_set"
sidebar_current: "docs-tencentcloud-resource-cat_task_set"
description: |-
  Provides a resource to create a cat task_set
---

# tencentcloud_cat_task_set

Provides a resource to create a cat task_set

## Example Usage

```hcl
resource "tencentcloud_cat_task_set" "task_set" {
  batch_tasks {
    name           = "demo"
    target_address = "http://www.baidu.com"
  }
  task_type = 5
  nodes     = ["12136", "12137", "12138", "12141", "12144"]
  interval  = 6
  parameters = jsonencode(
    {
      "ipType"            = 0,
      "grabBag"           = 0,
      "filterIp"          = 0,
      "netIcmpOn"         = 1,
      "netIcmpActivex"    = 0,
      "netIcmpTimeout"    = 20,
      "netIcmpInterval"   = 0.5,
      "netIcmpNum"        = 20,
      "netIcmpSize"       = 32,
      "netIcmpDataCut"    = 1,
      "netDnsOn"          = 1,
      "netDnsTimeout"     = 5,
      "netDnsQuerymethod" = 1,
      "netDnsNs"          = "",
      "netDigOn"          = 1,
      "netDnsServer"      = 2,
      "netTracertOn"      = 1,
      "netTracertTimeout" = 60,
      "netTracertNum"     = 30,
      "whiteList"         = "",
      "blackList"         = "",
      "netIcmpActivexStr" = ""
    }
  )
  task_category = 1
  cron          = "* 0-6 * * *"
  tags = {
    "createdBy" = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `batch_tasks` - (Required, List) Batch task name address.
* `interval` - (Required, Int) Task interval minutes in (1,5,10,15,30,60,120,240).
* `nodes` - (Required, Set: [`String`]) Task Nodes.
* `parameters` - (Required, String) tasks parameters.
* `task_category` - (Required, Int) Task category,1:PC,2:Mobile.
* `task_type` - (Required, Int) Task Type 1:Page Performance, 2:File upload,3:File Download,4:Port performance 5:Audio and video.
* `cron` - (Optional, String) Timer task cron expression.
* `operate` - (Optional, String) The input is valid when the parameter is modified, `suspend`/`resume`, used to suspend/resume the dial test task.
* `tags` - (Optional, Map) Tag description list.

The `batch_tasks` object supports the following:

* `name` - (Required, String) Task name.
* `target_address` - (Required, String) Target address.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `status` - Task status 1:TaskPending, 2:TaskRunning,3:TaskRunException,4:TaskSuspending 5:TaskSuspendException,6:TaskSuspendException,7:TaskSuspended,9:TaskDeleted.
* `task_id` - Task Id.


## Import

cat task_set can be imported using the id, e.g.
```
$ terraform import tencentcloud_cat_task_set.task_set taskSet_id
```

