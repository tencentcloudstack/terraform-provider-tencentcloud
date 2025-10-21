---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_log_post_cls_flow"
sidebar_current: "docs-tencentcloud-resource-waf_log_post_cls_flow"
description: |-
  Provides a resource to create a WAF log post cls flow
---

# tencentcloud_waf_log_post_cls_flow

Provides a resource to create a WAF log post cls flow

## Example Usage

### If log_type is 1

```hcl
resource "tencentcloud_waf_log_post_cls_flow" "example" {
  cls_region     = "ap-guangzhou"
  logset_name    = "waf_post_logset"
  log_type       = 1
  log_topic_name = "waf_post_logtopic"
}
```

### If log_type is 2

```hcl
resource "tencentcloud_waf_log_post_cls_flow" "example" {
  cls_region     = "ap-guangzhou"
  logset_name    = "waf_post_logset"
  log_type       = 2
  log_topic_name = "waf_post_logtopic"
}
```

## Argument Reference

The following arguments are supported:

* `cls_region` - (Optional, String) The region where the CLS is delivered. The default value is ap-shanghai.
* `log_topic_name` - (Optional, String) The name of the log subject where the submitted CLS is located. The default value is waf_post_logtopic.
* `log_type` - (Optional, Int) 1- Access log, 2- Attack log, the default is access log.
* `logset_name` - (Optional, String) The name of the log set where the delivered CLS is located. The default value is waf_post_logset.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `flow_id` - Unique ID for post cls flow.
* `log_topic_id` - CLS log topic ID.
* `logset_id` - CLS logset ID.
* `status` - Status 0- Off 1- On.


## Import

WAF log post cls flow can be imported using the id, e.g.

```
# If log_type is 1
terraform import tencentcloud_waf_log_post_cls_flow.example 111462#1

# If log_type is 2
terraform import tencentcloud_waf_log_post_cls_flow.example 111467#2
```

