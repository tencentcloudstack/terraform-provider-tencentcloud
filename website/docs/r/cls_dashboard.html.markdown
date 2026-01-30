---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_dashboard"
sidebar_current: "docs-tencentcloud-resource-cls_dashboard"
description: |-
  Provides a resource to create a CLS Dashboard.
---

# tencentcloud_cls_dashboard

Provides a resource to create a CLS Dashboard.

## Example Usage

```hcl
resource "tencentcloud_cls_dashboard" "dashboard" {
  dashboard_name = "my-dashboard"
}
```

### With configuration data

```hcl
resource "tencentcloud_cls_dashboard" "dashboard" {
  dashboard_name = "production-dashboard"
  data = jsonencode({
    timezone = "browser"
    subType  = "CLS_Host"
  })

  tags = {
    "team"        = "ops"
    "environment" = "production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `dashboard_name` - (Required, String) Dashboard name, which must be unique within the account.
* `data` - (Optional, String) Dashboard configuration data in JSON format. If not specified, an empty dashboard will be created.
* `tags` - (Optional, Map) Tag key-value pairs. Maximum of 10 tags.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time.
* `dashboard_id` - Dashboard ID (globally unique identifier).
* `update_time` - Last update time.


## Import

CLS dashboard can be imported using the id, e.g.

```
$ terraform import tencentcloud_cls_dashboard.dashboard dashboard-xxxx-xxxx-xxxx-xxxx
```

