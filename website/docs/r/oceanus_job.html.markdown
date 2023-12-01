---
subcategory: "Oceanus"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_oceanus_job"
sidebar_current: "docs-tencentcloud-resource-oceanus_job"
description: |-
  Provides a resource to create a oceanus job
---

# tencentcloud_oceanus_job

Provides a resource to create a oceanus job

## Example Usage

```hcl
resource "tencentcloud_oceanus_job" "example" {
  name          = "tf_example_job"
  job_type      = 1
  cluster_type  = 2
  cluster_id    = "cluster-1kcd524h"
  cu_mem        = 4
  remark        = "remark."
  folder_id     = "folder-7ctl246z"
  flink_version = "Flink-1.16"
  work_space_id = "space-2idq8wbr"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_type` - (Required, Int) The type of the cluster. 1 indicates shared cluster, and 2 indicates exclusive cluster.
* `job_type` - (Required, Int) The type of the job. 1 indicates SQL job, and 2 indicates JAR job.
* `name` - (Required, String) The name of the job. It can be composed of Chinese, English, numbers, hyphens (-), underscores (_), and periods (.), and the length cannot exceed 50 characters. Note that the job name cannot be the same as an existing job.
* `cluster_id` - (Optional, String) When ClusterType=2, it is required to specify the ID of the exclusive cluster to which the job is submitted.
* `cu_mem` - (Optional, Int) Set the memory specification of each CU, in GB. It supports 2, 4, 8, and 16 (which needs to apply for the whitelist before use). The default is 4, that is, 1 CU corresponds to 4 GB of running memory.
* `flink_version` - (Optional, String) The Flink version that the job runs.
* `folder_id` - (Optional, String) The folder ID to which the job name belongs. The root directory is root.
* `remark` - (Optional, String) The remark information of the job. It can be set arbitrarily.
* `work_space_id` - (Optional, String) The workspace SerialId.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



