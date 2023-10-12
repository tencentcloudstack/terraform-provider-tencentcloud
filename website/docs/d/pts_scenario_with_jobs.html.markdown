---
subcategory: "Performance Testing Service(PTS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_pts_scenario_with_jobs"
sidebar_current: "docs-tencentcloud-datasource-pts_scenario_with_jobs"
description: |-
  Use this data source to query detailed information of pts scenario_with_jobs
---

# tencentcloud_pts_scenario_with_jobs

Use this data source to query detailed information of pts scenario_with_jobs

## Example Usage

```hcl
data "tencentcloud_pts_scenario_with_jobs" "scenario_with_jobs" {
  project_ids    = ["project-45vw7v82"]
  scenario_ids   = ["scenario-koakp3h6"]
  scenario_name  = "pts-jmeter"
  ascend         = true
  ignore_script  = true
  ignore_dataset = true
  scenario_type  = "pts-jmeter"
}
```

## Argument Reference

The following arguments are supported:

* `project_ids` - (Required, Set: [`String`]) Project ID list.
* `ascend` - (Optional, Bool) Whether to use ascending order.
* `ignore_dataset` - (Optional, Bool) Whether to ignore the dataset.
* `ignore_script` - (Optional, Bool) Whether to ignore the script content.
* `order_by` - (Optional, String) The field column used for ordering.
* `owner` - (Optional, String) The job owner.
* `result_output_file` - (Optional, String) Used to save results.
* `scenario_ids` - (Optional, Set: [`String`]) Scenario ID list.
* `scenario_name` - (Optional, String) Scenario name.
* `scenario_type` - (Optional, String) Scenario type, e.g.: pts-http, pts-js, pts-trpc, pts-jmeter.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `scenario_with_jobs_set` - The scenario configuration and its jobs.
  * `jobs` - Jobs related to the scenario.
    * `abort_reason` - The reason for aborting the job.
    * `configs` - Deprecated.
    * `created_at` - The job creation time.
    * `cron_id` - Cron job ID.
    * `datasets` - The test data sets.
      * `file_id` - File ID.
      * `head_lines` - The header lines of the file.
      * `header_columns` - The parameter name list.
      * `header_in_file` - Whether the first line contains the parameter names.
      * `line_count` - The line count of the file.
      * `name` - Test data set name.
      * `size` - File size.
      * `split` - Whether to split the test data.
      * `tail_lines` - The tail lines of the file.
      * `type` - File type.
      * `updated_at` - The time of the most recent update.
    * `debug` - Whether to run the job in the debug mode. The default value is false.
    * `domain_name_config` - The configuration for parsing domain names.
      * `dns_config` - The DNS configuration.
        * `nameservers` - DNS IP list.
      * `host_aliases` - The configuration for host aliases.
        * `host_names` - Host names.
        * `ip` - IP.
    * `duration` - Job running duration.
    * `end_time` - The job ending time.
    * `error_rate` - Error rate.
    * `extensions` - Deprecated.
    * `job_id` - Job ID.
    * `job_owner` - Job owner.
    * `load_source_infos` - The load source information.
      * `ip` - The IP of the load source.
      * `pod_name` - The pod name of the load source.
      * `region` - Region.
    * `load_sources` - Deprecated.
      * `ip` - The IP of the load source.
      * `pod_name` - The pod name of the load source.
      * `region` - Region.
    * `load` - The configuration of the load.
      * `geo_regions_load_distribution` - The geographical distribution of the load source.
        * `percentage` - Percentage.
        * `region_id` - Region ID.
        * `region` - Region.
      * `load_spec` - The specification of the load configuration.
        * `concurrency` - The configuration of the concurrency load test mode.
          * `graceful_stop_seconds` - The waiting period for a graceful shutdown.
          * `iteration_count` - The iteration count.
          * `max_requests_per_second` - The maximum RPS.
          * `resources` - The count of the load test resource.
          * `stages` - The multi-stage configuration.
            * `duration_seconds` - The execution time.
            * `target_virtual_users` - The target count of the virtual users.
        * `requests_per_second` - The configuration of the RPS mode load test.
          * `duration_seconds` - The execution time.
          * `graceful_stop_seconds` - The waiting period for a gracefulshutdown.
          * `max_requests_per_second` - The maximum RPS.
          * `resources` - The count of the load test resource.
          * `start_requests_per_second` - The starting minimum RPS.
          * `target_requests_per_second` - The target RPS.
          * `target_virtual_users` - Deprecated.
        * `script_origin` - The script origin.
          * `duration_seconds` - The execution time.
          * `machine_number` - Machine number.
          * `machine_specification` - Machine specification.
      * `vpc_load_distribution` - The distribution of the load source.
        * `region_id` - Region ID.
        * `region` - Region.
        * `subnet_ids` - The subnet ID list.
        * `vpc_id` - VPC ID.
    * `max_requests_per_second` - The maximum RPS.
    * `max_virtual_user_count` - The maximum VU of the job.
    * `message` - The message describing the job running status.
    * `network_receive_rate` - The rate of receiving bytes.
    * `network_send_rate` - The rate of sending bytes.
    * `note` - The note of the job.
    * `notification_hooks` - Notification hooks.
      * `events` - Notification hook.
      * `url` - The callback URL.
    * `plugins` - Plugins.
      * `file_id` - File ID.
      * `name` - File name.
      * `size` - File size.
      * `type` - File type.
      * `updated_at` - The time of the most recent update.
    * `project_id` - Project ID.
    * `project_name` - Project name.
    * `protocols` - The protocol file.
      * `file_id` - File ID.
      * `name` - File name.
      * `size` - File size.
      * `type` - File type.
      * `updated_at` - The time of the most recent update.
    * `request_files` - The files in the request.
      * `file_id` - File ID.
      * `name` - File name.
      * `size` - File size.
      * `type` - File type.
      * `updated_at` - The time of the most recent update.
    * `request_total` - The total reqeust count.
    * `requests_per_second` - RPS.
    * `response_time_average` - The average response time.
    * `response_time_max` - The maximum response time.
    * `response_time_min` - The minimum response time.
    * `response_time_p90` - The 90 percentile of the response time.
    * `response_time_p95` - The 95 percentile of the response time.
    * `response_time_p99` - The 99 percentile of the response time.
    * `scenario_id` - Scenario ID.
    * `scenario_name` - Scenario name.
    * `scripts` - Deprecated.
    * `start_time` - The job starting time.
    * `status` - Job running status. JobUnknown: 0,JobCreated:1,JobPending:2, JobPreparing:3,JobSelectClustering:4,JobCreateTasking:5,JobSyncTasking:6 JobRunning:11,JobFinished:12,JobPrepareException:13,JobFinishException:14,JobAborting:15,JobAborted:16,JobAbortException:17,JobDeleted:18, JobSelectClusterException:19,JobCreateTaskException:20,JobSyncTaskException:21.
    * `test_scripts` - Test scripts.
      * `encoded_content` - The base64 encoded content.
      * `encoded_http_archive` - The base64 encoded HAR.
      * `file_id` - File ID.
      * `load_weight` - The weight of the script, ranging from 1 to 100.
      * `name` - File name.
      * `size` - File size.
      * `type` - File type.
      * `updated_at` - The time of the most recent update.
    * `type` - Scenario type.
  * `scenario` - The returned scenario.
    * `app_id` - AppId.
    * `configs` - Deprecated.
    * `created_at` - The creation time of the scenario.
    * `cron_id` - The cron job ID.
    * `datasets` - The test data sets for the load test.
      * `file_id` - The file ID.
      * `head_lines` - The header lines of the file.
      * `header_columns` - The parameter name list.
      * `header_in_file` - Whether the first line contains the parameter names.
      * `line_count` - The line count of the file.
      * `name` - The file name of the test data sets.
      * `size` - The byte count of the file.
      * `split` - Whether to split the test data.
      * `tail_lines` - The tail lines of the file.
      * `type` - The file type.
      * `updated_at` - The update time of the file.
    * `description` - Scenario description.
    * `domain_name_config` - The configuration for parsing domain names.
      * `dns_config` - The DNS configuration.
        * `nameservers` - DNS IP list.
      * `host_aliases` - The configuration for host aliases.
        * `host_names` - Host names.
        * `ip` - IP.
    * `encoded_scripts` - Deprecated.
    * `extensions` - Deprecated.
    * `load` - Scenario is load test configuration.
      * `geo_regions_load_distribution` - The geographical distribution of the load source.
        * `percentage` - Percentage.
        * `region_id` - Region ID.
        * `region` - Region.
      * `load_spec` - Scenario is load specification.
        * `concurrency` - The configuration for the concurrency mode.
          * `graceful_stop_seconds` - The waiting period for a graceful shutdown.
          * `iteration_count` - The iteration count of the load test.
          * `max_requests_per_second` - The maximum RPS.
          * `resources` - The resource count of the load test.
          * `stages` - The configuration for the multi-stage load test.
            * `duration_seconds` - The execution time for the load test.
            * `target_virtual_users` - The number of the target virtual users.
        * `requests_per_second` - The configuration of the RPS mode load test.
          * `duration_seconds` - The execution time of the load test.
          * `graceful_stop_seconds` - The waiting period for a graceful shutdown.
          * `max_requests_per_second` - The maximum RPS.
          * `resources` - The recource count of the load test.
          * `start_requests_per_second` - The starting minimum RPS.
          * `target_requests_per_second` - The target RPS.
          * `target_virtual_users` - Deprecated.
        * `script_origin` - The script origin.
          * `duration_seconds` - The load test execution time.
          * `machine_number` - The load test machine number.
          * `machine_specification` - The load test machine specification.
      * `vpc_load_distribution` - The distribution of the load source.
        * `region_id` - Region ID.
        * `region` - Region.
        * `subnet_ids` - The subnet ID list.
        * `vpc_id` - The VPC ID.
    * `name` - Scenario name.
    * `notification_hooks` - The notification hooks.
      * `events` - The notification hook.
      * `url` - The callback URL.
    * `owner` - The owner.
    * `plugins` - Plugins.
      * `file_id` - File ID.
      * `name` - File name.
      * `size` - File size.
      * `type` - File type.
      * `updated_at` - The time of the most recent update.
    * `project_id` - Project ID.
    * `project_name` - Project name.
    * `protocols` - The protocol file.
      * `file_id` - File ID.
      * `name` - File name.
      * `size` - File size.
      * `type` - File type.
      * `updated_at` - The time of the most recent update.
    * `request_files` - The files in the request.
      * `file_id` - File ID.
      * `name` - File name.
      * `size` - File size.
      * `type` - File type.
      * `updated_at` - The time of the most recent update.
    * `scenario_id` - Scenario ID.
    * `sla_id` - The ID of the SLA policy.
    * `sla_policy` - The SLA policy.
      * `alert_channel` - The alert channel.
        * `amp_consumer_id` - AMP consumer ID.
        * `notice_id` - The notice ID bound with this alert channel.
      * `sla_rules` - The SLA rules.
        * `abort_flag` - Whether to abort the load test job.
        * `aggregation` - The aggregation method of the metrics.
        * `condition` - The operator for checking the condition.
        * `for` - The duration for checking the condition.
        * `label_filter` - The label filter.
          * `label_name` - Label name.
          * `label_value` - Label value.
        * `metric` - The load test metrics.
        * `value` - The threshold in the condition.
    * `status` - Scenario status.
    * `sub_account_uin` - SubAccountUin.
    * `test_scripts` - The script of the load test.
      * `encoded_content` - The base64 encoded content.
      * `encoded_http_archive` - The base64 encoded HAR.
      * `file_id` - File ID.
      * `load_weight` - The weight of the script, ranging from 1 to 100.
      * `name` - File name.
      * `size` - File size.
      * `type` - File type.
      * `updated_at` - The time of the most recent update.
    * `type` - Scenario type, e.g.: pts-http, pts-js, pts-trpc, pts-jmeter.
    * `uin` - Uin.
    * `updated_at` - The updating time of the scenario.


