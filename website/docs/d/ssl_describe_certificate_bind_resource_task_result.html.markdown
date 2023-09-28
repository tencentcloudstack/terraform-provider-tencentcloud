---
subcategory: "SSL Certificates"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssl_describe_certificate_bind_resource_task_result"
sidebar_current: "docs-tencentcloud-datasource-ssl_describe_certificate_bind_resource_task_result"
description: |-
  Use this data source to query detailed information of ssl describe_certificate_bind_resource_task_result
---

# tencentcloud_ssl_describe_certificate_bind_resource_task_result

Use this data source to query detailed information of ssl describe_certificate_bind_resource_task_result

## Example Usage

```hcl
data "tencentcloud_ssl_describe_certificate_bind_resource_task_result" "describe_certificate_bind_resource_task_result" {
  task_ids =
}
```

## Argument Reference

The following arguments are supported:

* `task_ids` - (Required, Set: [`String`]) Task ID, query the results of binding cloud resources according to the task ID, support the maximum support of 100.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `sync_task_bind_resource_result` - List of asynchronous tasks binding affiliated cloud resources resultsNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `bind_resource_result` - Related Cloud Resources ResultNote: This field may return NULL, indicating that the valid value cannot be obtained.
    * `bind_resource_region_result` - Binding resource area results.
      * `region` - areaNote: This field may return NULL, indicating that the valid value cannot be obtained.
      * `total_count` - Total number of related resources.
    * `resource_type` - Resource types: CLB, CDN, Waf, LIVE, VOD, DDOS, TKE, Apigateway, TCB, Teo (Edgeone).
  * `cache_time` - Current result cache time.
  * `error` - Associated Cloud Resource Error InformationNote: This field may return NULL, indicating that the valid value cannot be obtained.
    * `code` - Unusual error codeNote: This field may return NULL, indicating that the valid value cannot be obtained.
    * `message` - Unusual error messageNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `status` - Related Cloud Resources Inquiry results: 0 indicates that in the query, 1 means the query is successful.2 means the query is abnormal; if the status is 1, check the results of bindResourceResult; if the state is 2, check the reason for ERROR.
  * `task_id` - Task ID.


