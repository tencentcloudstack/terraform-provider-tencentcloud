---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_cls_log_attachment"
sidebar_current: "docs-tencentcloud-resource-clb_cls_log_attachment"
description: |-
  Provides a resource to create a clb clb_cls_log_attachment
---

# tencentcloud_clb_cls_log_attachment

Provides a resource to create a clb clb_cls_log_attachment

## Example Usage

```hcl
resource "tencentcloud_clb_cls_log_attachment" "clb_cls_log_attachment" {
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required, String, ForceNew) CLB instance ID.
* `log_set_id` - (Required, String, ForceNew) Logset ID of the Cloud Log Service (CLS).<li>When adding or updating a log topic, call the [DescribeLogsets](https://intl.cloud.tencent.com/document/product/614/58624?from_cn_redirect=1) API to obtain the logset ID.</li><li>When deleting a log topic, set this parameter to null.</li>.
* `log_topic_id` - (Required, String, ForceNew) Log topic ID of the CLS.<li>When adding or updating a log topic, call the [DescribeTopics](https://intl.cloud.tencent.com/document/product/614/56454?from_cn_redirect=1) API to obtain the log topic ID.</li><li>When deleting a log topic, set this parameter to null.</li>.
* `log_type` - (Optional, String, ForceNew) Log type:
<li>`ACCESS`: access logs</li>
<li>`HEALTH`: health check logs</li>
Default: `ACCESS`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

clb clb_cls_log_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_clb_cls_log_attachment.clb_cls_log_attachment clb_cls_log_attachment_id
```

