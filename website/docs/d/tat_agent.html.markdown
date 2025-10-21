---
subcategory: "TencentCloud Automation Tools(TAT)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tat_agent"
sidebar_current: "docs-tencentcloud-datasource-tat_agent"
description: |-
  Use this data source to query detailed information of tat agent
---

# tencentcloud_tat_agent

Use this data source to query detailed information of tat agent

## Example Usage

```hcl
data "tencentcloud_tat_agent" "agent" {
  # instance_ids = ["ins-f9jr4bd2"]
  filters {
    name   = "environment"
    values = ["Linux"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter conditions. agent-status - String - Required: No - (Filter condition) Filter by agent status. Valid values: Online, Offline. environment - String - Required: No - (Filter condition) Filter by the agent environment. Valid value: Linux. instance-id - String - Required: No - (Filter condition) Filter by the instance ID. Up to 10 Filters allowed in one request. For each filter, five Filter.Values can be specified. InstanceIds and Filters cannot be specified at the same time.
* `instance_ids` - (Optional, Set: [`String`]) List of instance IDs for the query.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Field to be filtered.
* `values` - (Required, Set) Filter values of the field.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `automation_agent_set` - List of agent message.
  * `agent_status` - Agent status.Ranges:&lt;li&gt; Online:Online&lt;li&gt; Offline:Offline.
  * `environment` - Environment for Agent.Ranges:&lt;li&gt; Linux:Linux instance&lt;li&gt; Windows:Windows instance.
  * `instance_id` - InstanceId.
  * `last_heartbeat_time` - Time of last heartbeat.
  * `support_features` - List of feature Agent support.
  * `version` - Agent version.


