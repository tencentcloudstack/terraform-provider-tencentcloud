---
subcategory: "CLS"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_logset"
sidebar_current: "docs-tencentcloud-resource-cls_logset"
description: |-
  Provides a resource to create a cls logset
---

# tencentcloud_cls_logset

Provides a resource to create a cls logset

## Example Usage

```hcl
resource "tencentcloud_cls_logset" "logset" {
  logset_name = "logset"
  tags = {
    "test" = "test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `logset_name` - (Required, String) Logset name, which must be unique.
* `tags` - (Optional, Map) Tag description list. Up to 10 tag key-value pairs are supported and must be unique.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time.
* `role_name` - If AssumerUin is not empty, it indicates the service provider who creates the logset.
* `topic_count` - Number of log topics in logset.


## Import

cls logset can be imported using the id, e.g.

```
$ terraform import tencentcloud_cls_logset.logset 5cd3a17e-fb0b-418c-afd7-77b365397426
```

