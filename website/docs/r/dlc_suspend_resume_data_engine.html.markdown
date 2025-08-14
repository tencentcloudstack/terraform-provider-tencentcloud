---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_suspend_resume_data_engine"
sidebar_current: "docs-tencentcloud-resource-dlc_suspend_resume_data_engine"
description: |-
  Provides a resource to create a DLC suspend resume data engine
---

# tencentcloud_dlc_suspend_resume_data_engine

Provides a resource to create a DLC suspend resume data engine

## Example Usage

```hcl
resource "tencentcloud_dlc_suspend_resume_data_engine" "example" {
  data_engine_name = "tf-example"
  operate          = "suspend"
}
```

## Argument Reference

The following arguments are supported:

* `data_engine_name` - (Required, String, ForceNew) The name of a virtual cluster.
* `operate` - (Required, String, ForceNew) The operation type: `suspend` or `resume`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



