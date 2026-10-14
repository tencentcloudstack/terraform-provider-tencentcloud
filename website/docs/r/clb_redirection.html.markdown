---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_redirection"
sidebar_current: "docs-tencentcloud-resource-clb_redirection"
description: |-
  Provides a resource to create a CLB redirection.
---

# tencentcloud_clb_redirection

Provides a resource to create a CLB redirection.

## Example Usage

### Manual Rewrite

```hcl
resource "tencentcloud_clb_redirection" "example" {
  clb_id             = "lb-ab09jtd2"
  source_listener_id = "lbl-qgtfowas"
  target_listener_id = "lbl-lpwdkukk"
  source_rule_id     = "loc-liz99mtg"
  target_rule_id     = "loc-4f53xn52"
  rewrite_code       = 307
  take_url           = true
  source_domian      = "www.demo.com"
}
```

### Auto Rewrite

```hcl
resource "tencentcloud_clb_redirection" "example" {
  clb_id             = "lb-ab09jtd2"
  target_listener_id = "lbl-l7550kum"
  target_rule_id     = "loc-op7uz010"
  is_auto_rewrite    = true
}
```

## Argument Reference

The following arguments are supported:

* `clb_id` - (Required, String, ForceNew) ID of CLB instance.
* `target_listener_id` - (Required, String, ForceNew) ID of source listener.
* `target_rule_id` - (Required, String, ForceNew) Rule ID of target listener.
* `delete_all_auto_rewrite` - (Optional, Bool) Indicates whether delete all auto redirection. Default is `false`. It will take effect only when this redirection is auto-rewrite and this auto-rewrite auto redirected more than one rules. All the auto-rewrite relations will be deleted when this parameter set true.
* `is_auto_rewrite` - (Optional, Bool, ForceNew) Indicates whether automatic forwarding is enable, default is `false`. If enabled, the source listener and location should be empty, the target listener must be https protocol and port is 443.
* `rewrite_code` - (Optional, Int, ForceNew) Redirection status codes, with possible values of `301`, `302`, `307`.
* `source_domian` - (Optional, String, ForceNew) The domain name for source forwarding must be the domain name corresponding to `source_rule_id`, which is required when configuring `rewrite_code`. Only support `is_auto_rewrite` is `false`.
* `source_listener_id` - (Optional, String, ForceNew) ID of source listener.
* `source_rule_id` - (Optional, String, ForceNew) Rule ID of source listener.
* `take_url` - (Optional, Bool, ForceNew) Whether the redirect carries a matching URL is required when configuring `rewrite_code`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [AutoRewrite](https://www.tencentcloud.com/document/product/214/33846) - An HTTPS:443 listener needs to be created first, along with a forwarding rule. When this API is called, an HTTP:80 listener will be created automatically if it did not exist and a forwarding rule corresponding to `Domains` (specified in ...
* Create: [DescribeListeners](https://www.tencentcloud.com/document/product/214/33831) - This API is used to get the list of listeners by CLB ID, listener protocol, or listener port. If no filter is specified, all listeners for the CLB instance will be returned.
* Create: [ManualRewrite](https://www.tencentcloud.com/document/product/214/33821) - After the original access address and the address to be redirected are configured manually, the system will automatically redirect requests made to the original access address to the target address of the corresponding path.
* Read: [DescribeRewrite](https://www.tencentcloud.com/document/product/214/33829) - This API (DescribeRewrite) is used to query the redirection relationship between the forwarding rules of a CLB instance by instance ID.
* Delete: [DeleteRewrite](https://www.tencentcloud.com/document/product/214/33836) - This API (DeleteRewrite) is used to delete the redirection relationship between the specified forwarding rules.
* Delete: [DescribeRewrite](https://www.tencentcloud.com/document/product/214/33829) - This API (DescribeRewrite) is used to query the redirection relationship between the forwarding rules of a CLB instance by instance ID.

## Import

CLB redirection can be imported using the sourceLocId#targetLocId#sourceListenerId#targetListenerId#clbId, e.g.

```
terraform import tencentcloud_clb_redirection.example loc-ft8fmngv#loc-4xxr2cy7#lbl-jc1dx6ju#lbl-asj1hzuo#lb-p7olt9e5
```

