---
subcategory: "Web & Agents Development Platform(Makers)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_makers_deployment_log"
sidebar_current: "docs-tencentcloud-data_source-makers_deployment_log"
description: |-
  Use this data source to query the build log of a Makers deployment.
---

# tencentcloud_makers_deployment_log

Use this data source to query the build log of a Makers deployment.

## Example Usage

```hcl
data "tencentcloud_makers_deployment_log" "example" {
  zone_id       = "zone-123456"
  project_id    = "project-123456"
  deployment_id = "deployment-123456"
}

output "log_url" {
  value = data.tencentcloud_makers_deployment_log.example.log_url
}
```

## Argument Reference

The following arguments are supported:

* `deployment_id` - (Required, String) Deployment ID.
* `project_id` - (Required, String) Project ID.
* `zone_id` - (Required, String) Makers site ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `log_url` - Build log URL. This field may return null, which means no valid value is available.


