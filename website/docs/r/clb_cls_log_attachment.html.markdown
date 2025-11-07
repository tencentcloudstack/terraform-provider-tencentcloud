---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_cls_log_attachment"
sidebar_current: "docs-tencentcloud-resource-clb_cls_log_attachment"
description: |-
  Provides a resource to create a CLB cls log attachment
---

# tencentcloud_clb_cls_log_attachment

Provides a resource to create a CLB cls log attachment

## Example Usage

```hcl
resource "tencentcloud_clb_log_topic" "example" {
  log_set_id = "2ed70190-bf06-4777-980d-2d8a327a2554"
  topic_name = "tf-example"
  status     = true
}

resource "tencentcloud_clb_cls_log_attachment" "example" {
  load_balancer_id = "lb-n26tx0bm"
  log_set_id       = "2ed70190-bf06-4777-980d-2d8a327a2554"
  log_topic_id     = tencentcloud_clb_log_topic.example.id
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required, String, ForceNew) CLB instance ID.
* `log_set_id` - (Required, String, ForceNew) Logset ID of the Cloud Log Service (CLS).<li>When adding or updating a log topic, call the [DescribeLogsets](https://intl.cloud.tencent.com/document/product/614/58624?from_cn_redirect=1) API to obtain the logset ID.</li><li>When deleting a log topic, set this parameter to null.</li>.
* `log_topic_id` - (Required, String, ForceNew) Log topic ID of the CLS.<li>When adding or updating a log topic, call the [DescribeTopics](https://intl.cloud.tencent.com/document/product/614/56454?from_cn_redirect=1) API to obtain the log topic ID.</li><li>When deleting a log topic, set this parameter to null.</li>.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

CLB cls log attachment can be imported using the loadBalancerId#logSetId#logTopicId, e.g.

```
terraform import tencentcloud_clb_cls_log_attachment.example lb-n26tx0bm#2ed70190-bf06-4777-980d-2d8a327a2554#ac2fda28-3e79-4b51-b193-bfcf1aeece24
```

