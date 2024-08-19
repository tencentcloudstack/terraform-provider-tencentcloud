---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_logsets"
sidebar_current: "docs-tencentcloud-datasource-cls_logsets"
description: |-
  Use this data source to query detailed information of cls logsets
---

# tencentcloud_cls_logsets

Use this data source to query detailed information of cls logsets

## Example Usage

### Query all cls logsets

```hcl
data "tencentcloud_cls_logsets" "logsets" {}
```

### Query by filters

```hcl
# Query by `logsetName`
data "tencentcloud_cls_logsets" "logsets" {
  filters {
    key    = "logsetName"
    values = ["cls_service_logging"]
  }
}

# Query by `logsetId`
data "tencentcloud_cls_logsets" "logsets" {
  filters {
    key    = "logsetId"
    values = ["50d499a8-c4c0-4442-aa04-e8aa8a02437d"]
  }
}

# Query by `tagKey`
data "tencentcloud_cls_logsets" "logsets" {
  filters {
    key    = "tagKey"
    values = ["createdBy"]
  }
}

# Query by `tag:tagKey`
data "tencentcloud_cls_logsets" "logsets" {
  filters {
    key    = "tag:createdBy"
    values = ["terraform"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Query by filter.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `key` - (Required, String) Fields that need to be filtered. Support: `logsetName`, `logsetId`, `tagKey`, `tag:tagKey`.
* `values` - (Required, Set) The values that need to be filtered.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `logsets` - logset lists.
  * `assumer_name` - Cloud product identification, when the log set is created by another cloud product, this field will display the cloud product name, such as CDN, TKE.
  * `create_time` - Create time.
  * `logset_id` - Logset Id.
  * `logset_name` - Logset name.
  * `role_name` - If `assumer_name` is not empty, it indicates the service role that created the log set.
  * `tags` - Tags.
    * `key` - Tag key.
    * `value` - Tag value.
  * `topic_count` - Topic count.


