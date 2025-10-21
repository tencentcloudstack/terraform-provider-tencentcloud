---
subcategory: "Serverless Cloud Function(SCF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_scf_account_info"
sidebar_current: "docs-tencentcloud-datasource-scf_account_info"
description: |-
  Use this data source to query detailed information of scf account_info
---

# tencentcloud_scf_account_info

Use this data source to query detailed information of scf account_info

## Example Usage

```hcl
data "tencentcloud_scf_account_info" "account_info" {}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `account_limit` - Namespace limit information.
  * `namespace` - Namespace limit information.
    * `concurrent_executions` - Concurrency.
    * `functions_count` - Total number of functions.
    * `init_timeout_limit` - Initialization timeout limit.
    * `max_msg_ttl` - Upper limit of message retention time for async retry.
    * `min_msg_ttl` - Lower limit of message retention time for async retry.
    * `namespace` - Namespace name.
    * `retry_num_limit` - Limit of async retry attempt quantity.
    * `test_model_limit` - Test event limit Note: this field may return null, indicating that no valid values can be obtained.
    * `timeout_limit` - Timeout limit.
    * `trigger` - Trigger information.
      * `apigw` - Number of API Gateway triggers.
      * `ckafka` - Number of CKafka triggers.
      * `clb` - Number of CLB triggers.
      * `cls` - Number of CLS triggers.
      * `cm` - Number of CM triggers.
      * `cmq` - Number of CMQ triggers.
      * `cos` - Number of COS triggers.
      * `eb` - Number of EventBridge triggers Note: This field may return null, indicating that no valid values can be obtained.
      * `mps` - Number of MPS triggers.
      * `timer` - Number of timer triggers.
      * `total` - Total number of triggers.
      * `vod` - Number of VOD triggers.
  * `namespaces_count` - Limit of namespace quantity.
* `account_usage` - Namespace usage information.
  * `namespace` - Namespace details.
    * `functions_count` - Number of functions in namespace.
    * `functions` - Function array.
    * `namespace` - Namespace name.
    * `total_allocated_concurrency_mem` - Concurrency usage of the namespace Note: This field may return null, indicating that no valid value can be obtained.
    * `total_allocated_provisioned_mem` - Provisioned concurrency usage of the namespace Note: This field may return null, indicating that no valid value can be obtained.
    * `total_concurrency_mem` - Total memory quota of the namespace Note: This field may return null, indicating that no valid values can be obtained.
  * `namespaces_count` - Number of namespaces.
  * `total_allocated_concurrency_mem` - Quota of configured user concurrency memory in the current region.
  * `total_concurrency_mem` - Upper limit of user concurrency memory in the current region.
  * `user_concurrency_mem_limit` - Quota of account concurrency actually configured by user.


