---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_logset"
sidebar_current: "docs-tencentcloud-datasource-cls_logset"
description: |-
Use this data source to query detailed information of CLS logset

---

# tencentcloud_cls_logset

Use this data source to query detailed information of CLS logset

## Example Usage

```hcl
resource "tencentcloud_cls_logset" "logset_basic"{
     logset_name = "test"
}
data "tencentcloud_cls_logsets" "logsets" {
     filters {
                key = "logsetId"
                value = [tencentcloud_cls_logset.logset_basic.id]                                     
             }  
     limit = 2
}
```

## Argument Reference

The following arguments are supported:

* `offset` - (Optional) The offset of paging. The default value is 0.
* `limit` - (Optional) The limit number of single page paging. The default value is 20 and the maximum value is 100.
* `filters` - (Optional) filters of cls logsets,The upper limit of filters per request is 10.
  * `key` - (Required) Fields to be filtered, only supported logsetName,logsetId,tagKey,tag:tagKey
  * `value` - (Required) Values to be filtered,the upper limit of filter.values is 5.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `logsets` - Log set list. Each element contains the following attributes:
    * `logset_id` - Id of the logset.
    * `logset_name` - Name of the logset.
    * `create_time` - Creation time.
    * `topic_count` - Number of log topics under the logset.
    * `role_name` - If assumeruin is not empty, it indicates the server role that created the log set.
    * `tags`-Label of log set binding
      * `key`-Tag Key
      * `value`-Tag value
* `TotalCount` - Total number of pages



