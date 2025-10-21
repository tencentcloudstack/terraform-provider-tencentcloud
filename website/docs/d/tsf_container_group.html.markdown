---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_container_group"
sidebar_current: "docs-tencentcloud-datasource-tsf_container_group"
description: |-
  Use this data source to query detailed information of tsf container_group
---

# tencentcloud_tsf_container_group

Use this data source to query detailed information of tsf container_group

## Example Usage

```hcl
data "tencentcloud_tsf_container_group" "container_group" {
  application_id = "application-a24x29xv"
  search_word    = "keep"
  order_by       = "createTime"
  order_type     = 0
  cluster_id     = "cluster-vwgj5e6y"
  namespace_id   = "namespace-aemrg36v"
}
```

## Argument Reference

The following arguments are supported:

* `application_id` - (Optional, String) ApplicationId, required.
* `cluster_id` - (Optional, String) Cluster Id.
* `namespace_id` - (Optional, String) Namespace Id.
* `order_by` - (Optional, String) The sorting field. By default, it is the createTime field. Supports id, name, createTime.
* `order_type` - (Optional, Int) The sorting order. By default, it is 1, indicating descending order. 0 indicates ascending order, and 1 indicates descending order.
* `result_output_file` - (Optional, String) Used to save results.
* `search_word` - (Optional, String) search word, support group name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - result list.
  * `content` - List of deployment groups.Note: This field may return null, indicating that no valid value was found.
    * `alias` - The Group description.Note: This field may return null, indicating that no valid value was found.
    * `cluster_id` - Cluster Id.Note: This field may return null, indicating that no valid value was found.
    * `cluster_name` - Cluster name.Note: This field may return null, indicating that no valid value was found.
    * `cpu_limit` - The maximum amount of CPU, corresponding to K8S limit.Note: This field may return null, indicating that no valid value was found.
    * `cpu_request` - The initial amount of CPU, corresponding to K8S request.Note: This field may return null, indicating that no valid value was found.
    * `create_time` - Create time.Note: This field may return null, indicating that no valid value was found.
    * `group_id` - Group Id.Note: This field may return null, indicating that no valid value was found.
    * `group_name` - Group name.Note: This field may return null, indicating that no valid value was found.
    * `kube_inject_enable` - The value of KubeInjectEnable.Note: This field may return null, indicating that no valid value was found.
    * `mem_limit` - The maximum amount of memory allocated in MiB, corresponding to K8S limit.Note: This field may return null, indicating that no valid value was found.
    * `mem_request` - The initial amount of memory allocated in MiB, corresponding to K8S request.Note: This field may return null, indicating that no valid value was found.
    * `namespace_id` - Namespace Id.Note: This field may return null, indicating that no valid value was found.
    * `namespace_name` - Namespace name.Note: This field may return null, indicating that no valid value was found.
    * `repo_name` - Image name.Note: This field may return null, indicating that no valid value was found.
    * `server` - Image server.Note: This field may return null, indicating that no valid value was found.
    * `tag_name` - Image version Name.Note: This field may return null, indicating that no valid value was found.
    * `updated_time` - Update type.Note: This field may return null, indicating that no valid value was found.
  * `total_count` - Total count.


