---
subcategory: "Serverless Cloud Function(SCF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_scf_terminate_async_event"
sidebar_current: "docs-tencentcloud-resource-scf_terminate_async_event"
description: |-
  Provides a resource to create a scf terminate_async_event
---

# tencentcloud_scf_terminate_async_event

Provides a resource to create a scf terminate_async_event

## Example Usage

```hcl
resource "tencentcloud_scf_terminate_async_event" "terminate_async_event" {
  function_name     = "keep-1676351130"
  invoke_request_id = "9de9405a-e33a-498d-bb59-e80b7bed1191"
  namespace         = "default"
  grace_shutdown    = true
}
```

## Argument Reference

The following arguments are supported:

* `function_name` - (Required, String, ForceNew) Function name.
* `invoke_request_id` - (Required, String, ForceNew) Terminated invocation request ID.
* `grace_shutdown` - (Optional, Bool, ForceNew) Whether to enable grace shutdown. If it's true, a SIGTERM signal is sent to the specified request. See [Sending termination signal](https://www.tencentcloud.com/document/product/583/63969?from_cn_redirect=1#.E5.8F.91.E9.80.81.E7.BB.88.E6.AD.A2.E4.BF.A1.E5.8F.B7]. It's set to false by default.
* `namespace` - (Optional, String, ForceNew) Namespace.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



