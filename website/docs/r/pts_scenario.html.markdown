---
subcategory: "Performance Testing Service(PTS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_pts_scenario"
sidebar_current: "docs-tencentcloud-resource-pts_scenario"
description: |-
  Provides a resource to create a pts scenario
---

# tencentcloud_pts_scenario

Provides a resource to create a pts scenario

## Example Usage

```hcl
resource "tencentcloud_pts_scenario" "scenario" {
  name        = "pts"
  type        = ""
  project_id  = ""
  description = ""
  load {
    load_spec {
      concurrency {
        stages {
          duration_seconds     = ""
          target_virtual_users = ""
        }
        iteration_count         = ""
        max_requests_per_second = ""
        graceful_stop_seconds   = ""
      }
      requests_per_second {
        max_requests_per_second    = ""
        duration_seconds           = ""
        resources                  = ""
        start_requests_per_second  = ""
        target_requests_per_second = ""
        graceful_stop_seconds      = ""
      }
      script_origin {
        machine_number        = ""
        machine_specification = ""
        duration_seconds      = ""
      }
    }
    vpc_load_distribution {
      region_id  = ""
      region     = ""
      vpc_id     = ""
      subnet_ids = ""
    }
    geo_regions_load_distribution {
      region_id  = ""
      region     = ""
      percentage = ""
    }

  }
  datasets {
    name           = ""
    split          = ""
    header_in_file = ""
    header_columns = ""
    line_count     = ""
    updated_at     = ""
    size           = ""
    head_lines     = ""
    tail_lines     = ""
    type           = ""
    file_id        = ""

  }
  extensions = ""
  cron_id    = ""
  test_scripts {
    name                 = ""
    size                 = ""
    type                 = ""
    updated_at           = ""
    encoded_content      = ""
    encoded_http_archive = ""
    load_weight          = ""

  }
  protocols {
    name       = ""
    size       = ""
    type       = ""
    updated_at = ""
    file_id    = ""

  }
  request_files {
    name       = ""
    size       = ""
    type       = ""
    updated_at = ""
    file_id    = ""

  }
  sla_policy {
    sla_rules {
      metric      = ""
      aggregation = ""
      condition   = ""
      value       = ""
      label_filter {
        label_name  = ""
        label_value = ""
      }
      abort_flag = ""
      for        = ""
    }
    alert_channel {
      notice_id       = ""
      amp_consumer_id = ""
    }

  }
  plugins {
    name       = ""
    size       = ""
    type       = ""
    updated_at = ""
    file_id    = ""

  }
  domain_name_config {
    host_aliases {
      host_names = ""
      ip         = ""
    }
    dns_config {
      nameservers = ""
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Pts Scenario name.
* `project_id` - (Required, String) Project id.
* `type` - (Required, String) Pressure test engine type.
* `cron_id` - (Optional, String) cron job ID.
* `datasets` - (Optional, List) Test data set.
* `description` - (Optional, String) Pts Scenario Description.
* `domain_name_config` - (Optional, List) Domain name resolution configuration.
* `extensions` - (Optional, Set: [`String`]) deprecated.
* `load` - (Optional, List) Pressure allocation.
* `plugins` - (Optional, List) SLA strategy.
* `protocols` - (Optional, List) Protocol file path.
* `request_files` - (Optional, List) Request file path.
* `sla_policy` - (Optional, List) SLA strategy.
* `test_scripts` - (Optional, List) Test script file information.

The `alert_channel` object supports the following:

* `amp_consumer_id` - (Optional, String) AMP consumer ID.
* `notice_id` - (Optional, String) Notification template ID.

The `concurrency` object supports the following:

* `graceful_stop_seconds` - (Optional, Int) Wait time for graceful termination of the task.
* `iteration_count` - (Optional, Int) Number of runs.
* `max_requests_per_second` - (Optional, Int) Maximum RPS.
* `stages` - (Optional, List) Multi-phase configuration array.

The `datasets` object supports the following:

* `header_in_file` - (Required, Bool) Whether the first line is the parameter name.
* `name` - (Required, String) The file name where the test dataset is located.
* `split` - (Required, Bool) Test whether the dataset is fragmented.
* `file_id` - (Optional, String) File ID.
* `head_lines` - (Optional, Set) Header data row.
* `header_columns` - (Optional, Set) Parameter name array.
* `line_count` - (Optional, Int) Number of file lines.
* `size` - (Optional, Int) Number of file bytes.
* `tail_lines` - (Optional, Set) Trailing data row.
* `type` - (Optional, String) File type.
* `updated_at` - (Optional, String) Update time.

The `dns_config` object supports the following:

* `nameservers` - (Optional, Set) DNS IP List.

The `domain_name_config` object supports the following:

* `dns_config` - (Optional, List) DNS configuration.
* `host_aliases` - (Optional, List) Domain name binding configuration.

The `geo_regions_load_distribution` object supports the following:

* `region_id` - (Required, Int) Regional ID.
* `percentage` - (Optional, Int) Percentage.
* `region` - (Optional, String) Region.

The `host_aliases` object supports the following:

* `host_names` - (Optional, Set) List of domain names to be bound.
* `ip` - (Optional, String) The IP address to be bound.

The `label_filter` object supports the following:

* `label_name` - (Optional, String) Label name.
* `label_value` - (Optional, String) Label value.

The `load_spec` object supports the following:

* `concurrency` - (Optional, List) Configuration of concurrent pressure mode.
* `requests_per_second` - (Optional, List) Configuration of RPS pressure mode.
* `script_origin` - (Optional, List) Built-in stress mode in script.

The `load` object supports the following:

* `geo_regions_load_distribution` - (Optional, List) Pressure distribution.
* `load_spec` - (Optional, List) Pressure allocation.
* `vpc_load_distribution` - (Optional, List) Source of stress.

The `plugins` object supports the following:

* `file_id` - (Optional, String) File id.
* `name` - (Optional, String) File name.
* `size` - (Optional, Int) File size.
* `type` - (Optional, String) File type.
* `updated_at` - (Optional, String) Update time.

The `protocols` object supports the following:

* `file_id` - (Optional, String) File ID.
* `name` - (Optional, String) Protocol name.
* `size` - (Optional, Int) File name.
* `type` - (Optional, String) File type.
* `updated_at` - (Optional, String) Update time.

The `request_files` object supports the following:

* `file_id` - (Optional, String) File id.
* `name` - (Optional, String) File name.
* `size` - (Optional, Int) File size.
* `type` - (Optional, String) File type.
* `updated_at` - (Optional, String) Update time.

The `requests_per_second` object supports the following:

* `duration_seconds` - (Optional, Int) Pressure time.
* `graceful_stop_seconds` - (Optional, Int) Elegant shutdown waiting time.
* `max_requests_per_second` - (Optional, Int) Maximum RPS.
* `resources` - (Optional, Int) Number of resources.
* `start_requests_per_second` - (Optional, Int) Initial RPS.
* `target_requests_per_second` - (Optional, Int) Target RPS, invalid input parameter.

The `script_origin` object supports the following:

* `duration_seconds` - (Required, Int) Pressure testing time.
* `machine_number` - (Required, Int) Number of machines.
* `machine_specification` - (Required, String) Machine specification.

The `sla_policy` object supports the following:

* `alert_channel` - (Optional, List) Alarm notification channel.
* `sla_rules` - (Optional, List) SLA rules.

The `sla_rules` object supports the following:

* `abort_flag` - (Optional, Bool) Whether to stop the stress test task.
* `aggregation` - (Optional, String) Aggregation method of pressure test index.
* `condition` - (Optional, String) Pressure test index condition judgment symbol.
* `for` - (Optional, String) duraion.
* `label_filter` - (Optional, List) tag.
* `metric` - (Optional, String) Pressure test index.
* `value` - (Optional, Float64) Threshold value.

The `stages` object supports the following:

* `duration_seconds` - (Optional, Int) Pressure time.
* `target_virtual_users` - (Optional, Int) Number of virtual users.

The `test_scripts` object supports the following:

* `encoded_content` - (Optional, String) Base64 encoded file content.
* `encoded_http_archive` - (Optional, String) Base64 encoded har structure.
* `load_weight` - (Optional, Int) Script weight, range 1-100.
* `name` - (Optional, String) File name.
* `size` - (Optional, Int) File size.
* `type` - (Optional, String) File type.
* `updated_at` - (Optional, String) Update time.

The `vpc_load_distribution` object supports the following:

* `region_id` - (Required, Int) Regional ID.
* `region` - (Optional, String) Region.
* `subnet_ids` - (Optional, Set) Subnet ID list.
* `vpc_id` - (Optional, String) VPC ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `app_id` - App ID Note: this field may return null, indicating that a valid value cannot be obtained.
* `created_at` - Scene creation time.
* `notification_hooks` - Notification event callback Note: this field may return null, indicating that a valid value cannot be obtained.
* `status` - Scene statu Note: this field may return null, indicating that a valid value cannot be obtained.
* `sub_account_uin` - Sub-user ID Note: this field may return null, indicating that a valid value cannot be obtained.
* `uin` - User ID Note: this field may return null, indicating that a valid value cannot be obtained.
* `updated_at` - Scene modification time.


## Import

pts scenario can be imported using the id, e.g.
```
$ terraform import tencentcloud_pts_scenario.scenario scenario_id
```

