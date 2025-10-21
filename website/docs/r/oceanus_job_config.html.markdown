---
subcategory: "Oceanus"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_oceanus_job_config"
sidebar_current: "docs-tencentcloud-resource-oceanus_job_config"
description: |-
  Provides a resource to create a oceanus job_config
---

# tencentcloud_oceanus_job_config

Provides a resource to create a oceanus job_config

## Example Usage

### is 2

```hcl
resource "tencentcloud_oceanus_job_config" "example" {
  job_id           = "cql-4xwincyn"
  entrypoint_class = "tf_example"
  program_args     = "--conf Key=Value"
  remark           = "remark."
  resource_refs {
    resource_id = "resource-q22ntswy"
    version     = 1
    type        = 1
  }
  default_parallelism = 1
  properties {
    key   = "pipeline.max-parallelism"
    value = "2048"
  }
  log_collect       = true
  job_manager_spec  = "1"
  task_manager_spec = "1"
  cls_logset_id     = "cd9adbb5-6b7d-48d2-9870-77658959c7a4"
  cls_topic_id      = "cec4c2f1-0bf3-470e-b1a5-b1c451e88838"
  log_collect_type  = 2
  work_space_id     = "space-2idq8wbr"
  log_level         = "INFO"
  auto_recover      = 1
  expert_mode_on    = false
}
```

### is 3

```hcl
resource "tencentcloud_oceanus_job_config" "example" {
  job_id           = "cql-4xwincyn"
  entrypoint_class = "tf_example"
  program_args     = "--conf Key=Value"
  remark           = "remark."
  resource_refs {
    resource_id = "resource-q22ntswy"
    version     = 1
    type        = 1
  }
  default_parallelism = 1
  properties {
    key   = "pipeline.max-parallelism"
    value = "2048"
  }
  log_collect       = true
  job_manager_spec  = "1"
  task_manager_spec = "1"
  cls_logset_id     = "cd9adbb5-6b7d-48d2-9870-77658959c7a4"
  cls_topic_id      = "cec4c2f1-0bf3-470e-b1a5-b1c451e88838"
  log_collect_type  = 3
  work_space_id     = "space-2idq8wbr"
  log_level         = "INFO"
  auto_recover      = 1
  expert_mode_on    = false
  cos_bucket        = "autotest-gz-bucket-1257058945"
}
```

## Argument Reference

The following arguments are supported:

* `job_id` - (Required, String) Job ID.
* `auto_recover` - (Optional, Int) Oceanus platform job recovery switch 1: on -1: off.
* `clazz_levels` - (Optional, List) Class log level.
* `cls_logset_id` - (Optional, String) CLS logset ID.
* `cls_topic_id` - (Optional, String) CLS log topic ID.
* `cos_bucket` - (Optional, String) COS storage bucket name used by the job.
* `default_parallelism` - (Optional, Int) Job default parallelism.
* `entrypoint_class` - (Optional, String) Main class.
* `expert_mode_configuration` - (Optional, List) Expert mode configuration.
* `expert_mode_on` - (Optional, Bool) Whether to enable expert mode.
* `job_manager_spec` - (Optional, Float64) JobManager specification.
* `log_collect_type` - (Optional, Int) Log collection type 2:CLS; 3:COS.
* `log_collect` - (Optional, Bool) Whether to collect job logs.
* `log_level` - (Optional, String) Log level.
* `program_args` - (Optional, String) Main class parameters.
* `properties` - (Optional, List) System parameters.
* `python_version` - (Optional, String) Python version used by the pyflink job at runtime.
* `remark` - (Optional, String) Remarks.
* `resource_refs` - (Optional, List) Resource reference array.
* `task_manager_spec` - (Optional, Float64) TaskManager specification.
* `work_space_id` - (Optional, String) Workspace SerialId.

The `clazz_levels` object supports the following:

* `clazz` - (Required, String) Java class full pathNote: This field may return null, indicating that no valid value can be obtained.
* `level` - (Required, String) Log level TRACE, DEBUG, INFO, WARN, ERRORNote: This field may return null, indicating that no valid value can be obtained.

The `configuration` object of `node_config` supports the following:

* `key` - (Required, String) System configuration key.
* `value` - (Required, String) System configuration value.

The `edges` object of `job_graph` supports the following:

* `source` - (Required, Int) Starting node ID of the edgeNote: This field may return null, indicating that no valid value can be obtained.
* `target` - (Required, Int) Target node ID of the edgeNote: This field may return null, indicating that no valid value can be obtained.

The `expert_mode_configuration` object supports the following:

* `job_graph` - (Optional, List) Job graphNote: This field may return null, indicating that no valid value can be obtained.
* `node_config` - (Optional, List) Node configurationNote: This field may return null, indicating that no valid value can be obtained.
* `slot_sharing_groups` - (Optional, List) Slot sharing groupsNote: This field may return null, indicating that no valid value can be obtained.

The `job_graph` object of `expert_mode_configuration` supports the following:

* `edges` - (Optional, List) Edge set of the running graphNote: This field may return null, indicating that no valid value can be obtained.
* `nodes` - (Optional, List) Point set of the running graphNote: This field may return null, indicating that no valid value can be obtained.

The `node_config` object of `expert_mode_configuration` supports the following:

* `id` - (Required, Int) Node IDNote: This field may return null, indicating that no valid value can be obtained.
* `configuration` - (Optional, List) Configuration propertiesNote: This field may return null, indicating that no valid value can be obtained.
* `parallelism` - (Optional, Int) Node parallelismNote: This field may return null, indicating that no valid value can be obtained.
* `slot_sharing_group` - (Optional, String) Slot sharing groupNote: This field may return null, indicating that no valid value can be obtained.
* `state_ttl` - (Optional, String) State TTL configuration of the node, separated by semicolonsNote: This field may return null, indicating that no valid value can be obtained.

The `nodes` object of `job_graph` supports the following:

* `description` - (Required, String) Node descriptionNote: This field may return null, indicating that no valid value can be obtained.
* `id` - (Required, Int) Node IDNote: This field may return null, indicating that no valid value can be obtained.
* `name` - (Required, String) Node nameNote: This field may return null, indicating that no valid value can be obtained.
* `parallelism` - (Required, Int) Node parallelismNote: This field may return null, indicating that no valid value can be obtained.

The `properties` object supports the following:

* `key` - (Required, String) System configuration key.
* `value` - (Required, String) System configuration value.

The `resource_refs` object supports the following:

* `resource_id` - (Required, String) Resource ID.
* `type` - (Required, Int) Reference resource type, for example, setting the main resource to 1 represents the jar package where the main class is located.
* `version` - (Required, Int) Resource version ID, -1 indicates the latest version.

The `slot_sharing_groups` object of `expert_mode_configuration` supports the following:

* `name` - (Required, String) Name of the SlotSharingGroupNote: This field may return null, indicating that no valid value can be obtained.
* `spec` - (Required, List) Specification of the SlotSharingGroupNote: This field may return null, indicating that no valid value can be obtained.
* `description` - (Optional, String) Description of the SlotSharingGroupNote: This field may return null, indicating that no valid value can be obtained.

The `spec` object of `slot_sharing_groups` supports the following:

* `cpu` - (Required, Float64) Applicable CPUNote: This field may return null, indicating that no valid value can be obtained.
* `heap_memory` - (Required, String) Default is b, supporting units are b, kb, mb, gbNote: This field may return null, indicating that no valid value can be obtained.
* `managed_memory` - (Optional, String) Default is b, supporting units are b, kb, mb, gbNote: This field may return null, indicating that no valid value can be obtained.
* `off_heap_memory` - (Optional, String) Default is b, supporting units are b, kb, mb, gbNote: This field may return null, indicating that no valid value can be obtained.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



