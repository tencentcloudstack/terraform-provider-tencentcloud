---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_pod_instances"
sidebar_current: "docs-tencentcloud-datasource-tsf_pod_instances"
description: |-
  Use this data source to query detailed information of tsf pod_instances
---

# tencentcloud_tsf_pod_instances

Use this data source to query detailed information of tsf pod_instances

## Example Usage

```hcl
data "tencentcloud_tsf_pod_instances" "pod_instances" {
  group_id      = "group-ynd95rea"
  pod_name_list = ["keep-terraform-6f8f977688-zvphm"]
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, String) Instance&amp;#39;s group ID.
* `pod_name_list` - (Optional, Set: [`String`]) Filter, pod name list.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - pod instance list.
  * `content` - Content list.Note: This field may return null, which means no valid value was found.
    * `created_at` - Instance start time.Note: This field may return null, which means no valid value was found.
    * `instance_available_status` - Instance available status.Note: This field may return null, which means no valid value was found.
    * `instance_status` - Instance status.Note: This field may return null, which means no valid value was found.
    * `ip` - Instance ip.Note: This field may return null, which means no valid value was found.
    * `node_instance_id` - Instance node id.Note: This field may return null, which means no valid value was found.
    * `node_ip` - Instance node ip.Note: This field may return null, which means no valid value was found.
    * `pod_id` - Instance id (corresponding to the pod instance id in Kubernetes).Note: This field may return null, which means no valid value was found.
    * `pod_name` - Instance name (corresponding to the pod name in Kubernetes).Note: This field may return null, which means no valid value was found.
    * `ready_count` - Instance ready count.Note: This field may return null, which means no valid value was found.
    * `reason` - Instance reason for current status.Note: This field may return null, which means no valid value was found.
    * `restart_count` - Instance restart count.Note: This field may return null, which means no valid value was found.
    * `runtime` - Instance run time.Note: This field may return null, which means no valid value was found.
    * `service_instance_status` - Instance serve status.Note: This field may return null, which means no valid value was found.
    * `status` - Instance status. Please refer to the definition of instance and container status below. Starting (pod not ready): Starting; Running: Running; Abnormal: Abnormal; Stopped: Stopped;Note: This field may return null, which means no valid value was found.
  * `total_count` - Total number of records.Note: This field may return null, which means no valid value was found.


