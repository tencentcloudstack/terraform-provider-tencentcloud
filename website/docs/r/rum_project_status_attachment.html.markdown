---
subcategory: "Real User Monitoring(RUM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_rum_project_status_attachment"
sidebar_current: "docs-tencentcloud-resource-rum_project_status_attachment"
description: |-
  Provides a resource to create a rum project_status_attachment
---

# tencentcloud_rum_project_status_attachment

Provides a resource to create a rum project_status_attachment

## Example Usage

```hcl
resource "tencentcloud_rum_project_status_attachment" "project_status_attachment" {
  project_id = 131407
  operate    = "stop"
}
```

## Argument Reference

The following arguments are supported:

* `operate` - (Required, String) `resume`, `stop`.
* `project_id` - (Required, Int, ForceNew) Project ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

rum project_status_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_rum_project_status_attachment.project_status_attachment project_status_attachment_id
```

