---
subcategory: "Oceanus"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_oceanus_job_copy"
sidebar_current: "docs-tencentcloud-resource-oceanus_job_copy"
description: |-
  Provides a resource to create a oceanus job_copy
---

# tencentcloud_oceanus_job_copy

Provides a resource to create a oceanus job_copy

## Example Usage

```hcl
resource "tencentcloud_oceanus_job_copy" "example" {
  source_id         = "cql-0nob2hx8"
  target_cluster_id = "cluster-1kcd524h"
  source_name       = "keep_jar"
  target_name       = "tf_copy_example"
  target_folder_id  = "folder-7ctl246z"
  job_type          = 2
  work_space_id     = "space-2idq8wbr"
}
```

## Argument Reference

The following arguments are supported:

* `source_id` - (Required, String, ForceNew) The serial ID of the job to be copied.
* `target_cluster_id` - (Required, String, ForceNew) The cluster serial ID of the target cluster.
* `job_type` - (Optional, Int, ForceNew) The type of the source job.
* `source_name` - (Optional, String, ForceNew) The name of the job to be copied.
* `target_folder_id` - (Optional, String, ForceNew) The directory ID of the new job.
* `target_name` - (Optional, String, ForceNew) The name of the new job.
* `work_space_id` - (Optional, String, ForceNew) Workspace SerialId.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `job_id` - Copy Job ID.


