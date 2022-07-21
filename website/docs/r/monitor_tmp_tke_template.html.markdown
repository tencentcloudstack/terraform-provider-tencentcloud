---
subcategory: "Monitor"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_tmp_tke_template"
sidebar_current: "docs-tencentcloud-resource-monitor_tmp_tke_template"
description: |-
  Provides a resource to create a tmp tke template
---

# tencentcloud_monitor_tmp_tke_template

Provides a resource to create a tmp tke template

## Example Usage



## Argument Reference

The following arguments are supported:

* `template` - (Required, List) Template settings.

The `pod_monitors` object supports the following:

* `config` - (Required, String) Config.
* `name` - (Required, String) Name.
* `template_id` - (Optional, String) Used for the argument, if the configuration comes to the template, the template id.

The `raw_jobs` object supports the following:

* `config` - (Required, String) Config.
* `name` - (Required, String) Name.
* `template_id` - (Optional, String) Used for the argument, if the configuration comes to the template, the template id.

The `record_rules` object supports the following:

* `config` - (Required, String) Config.
* `name` - (Required, String) Name.
* `template_id` - (Optional, String) Used for the argument, if the configuration comes to the template, the template id.

The `service_monitors` object supports the following:

* `config` - (Required, String) Config.
* `name` - (Required, String) Name.
* `template_id` - (Optional, String) Used for the argument, if the configuration comes to the template, the template id.

The `template` object supports the following:

* `level` - (Required, String) Template dimensions, the following types are supported `instance` instance level, `cluster` cluster level.
* `name` - (Required, String) Template name.
* `describe` - (Optional, String) Template description.
* `is_default` - (Optional, Bool) Whether the system-supplied default template is used for outgoing references.
* `pod_monitors` - (Optional, List) Effective when Level is a cluster, A list of PodMonitors rules in the template.
* `raw_jobs` - (Optional, List) Effective when Level is a cluster, A list of RawJobs rules in the template.
* `record_rules` - (Optional, List) Effective when Level is instance, A list of aggregation rules in the template.
* `service_monitors` - (Optional, List) Effective when Level is a cluster, A list of ServiceMonitor rules in the template.
* `template_id` - (Optional, String) The ID of the template, which is used for the outgoing reference.
* `update_time` - (Optional, String) Last updated, for outgoing references.
* `version` - (Optional, String) Whether the system-supplied default template is used for outgoing references.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



