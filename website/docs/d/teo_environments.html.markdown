---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_environments"
sidebar_current: "docs-tencentcloud-datasource-teo_environments"
description: |-
  Use this data source to query detailed information of teo environments
---

# tencentcloud_teo_environments

Use this data source to query detailed information of teo environments

## Example Usage

```hcl
data "tencentcloud_teo_environments" "teo_environments" {
  zone_id = "zone-2qtuhspy7cr6"
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String) Zone ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `env_infos` - Environment list.


