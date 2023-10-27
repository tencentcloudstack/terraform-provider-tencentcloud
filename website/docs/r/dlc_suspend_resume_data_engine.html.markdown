---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_suspend_resume_data_engine"
sidebar_current: "docs-tencentcloud-resource-dlc_suspend_resume_data_engine"
description: |-
  Provides a resource to create a dlc suspend_resume_data_engine
---

# tencentcloud_dlc_suspend_resume_data_engine

Provides a resource to create a dlc suspend_resume_data_engine

## Example Usage

```hcl
resource "tencentcloud_dlc_suspend_resume_data_engine" "suspend_resume_data_engine" {
  data_engine_name = "example-iac"
  operate          = "suspend"
}
```

## Argument Reference

The following arguments are supported:

* `data_engine_name` - (Required, String, ForceNew) Engine name.
* `operate` - (Required, String, ForceNew) Engine operate tye: suspend/resume.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

dlc suspend_resume_data_engine can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_suspend_resume_data_engine.suspend_resume_data_engine suspend_resume_data_engine_id
```

