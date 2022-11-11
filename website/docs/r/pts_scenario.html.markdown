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
  name       = "pts-js"
  project_id = "project-45vw7v82"
  type       = "pts-js"

  domain_name_config {
  }

  load {
    geo_regions_load_distribution {
      percentage = 100
      region     = "ap-guangzhou"
      region_id  = 1
    }

    load_spec {
      concurrency {
        graceful_stop_seconds   = 3
        iteration_count         = 0
        max_requests_per_second = 0

        stages {
          duration_seconds     = 120
          target_virtual_users = 2
        }
        stages {
          duration_seconds     = 120
          target_virtual_users = 4
        }
        stages {
          duration_seconds     = 120
          target_virtual_users = 5
        }
        stages {
          duration_seconds     = 240
          target_virtual_users = 5
        }
      }
    }
  }

  sla_policy {
  }

  test_scripts {
    encoded_content = <<-EOT
            // Send a http get request
            import http from 'pts/http';
            import { check, sleep } from 'pts';

            export default function () {
              // simple get request
              const resp1 = http.get('http://httpbin.org/get');
              console.log(resp1.body);
              // if resp1.body is a json string, resp1.json() transfer json format body to a json object
              console.log(resp1.json());
              check('status is 200', () => resp1.statusCode === 200);

              // sleep 1 second
              sleep(1);

              // get request with headers and parameters
              const resp2 = http.get('http://httpbin.org/get', {
                headers: {
                  Connection: 'keep-alive',
                  'User-Agent': 'pts-engine',
                },
                query: {
                  name1: 'value1',
                  name2: 'value2',
                },
              });

              console.log(resp2.json().args.name1); // 'value1'
              check('body.args.name1 equals value1', () => resp2.json().args.name1 === 'value1');
            }
        EOT
    load_weight     = 100
    name            = "script.js"
    size            = 838
    type            = "js"
    updated_at      = "2022-11-11T16:18:37+08:00"
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
* `status` - Scene statu Note: this field may return null, indicating that a valid value cannot be obtained.
* `sub_account_uin` - Sub-user ID Note: this field may return null, indicating that a valid value cannot be obtained.
* `uin` - User ID Note: this field may return null, indicating that a valid value cannot be obtained.
* `updated_at` - Scene modification time.


## Import

pts scenario can be imported using the id, e.g.
```
$ terraform import tencentcloud_pts_scenario.scenario scenario_id
```

