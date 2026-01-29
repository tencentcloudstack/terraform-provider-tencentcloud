---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_object"
sidebar_current: "docs-tencentcloud-resource-waf_object"
description: |-
  Provides a resource to create a Waf object
---

# tencentcloud_waf_object

Provides a resource to create a Waf object

~> **NOTE:** If you need to change field `instance_id`, you need to keep `status` at `0`; If you need to change field `proxy(ip_headers)`, you need to keep `status` at `1`.

## Example Usage

### Bind current account resources

```hcl
resource "tencentcloud_waf_object" "example" {
  object_id   = "lb-9h5x9lze"
  instance_id = "waf_2kxtlbky11b2v4fe"
  status      = 1
  proxy       = 3
  ip_headers = [
    "my-header1",
    "my-header2",
    "my-header3",
  ]
}
```

### Bind other member account resources

```hcl
resource "tencentcloud_waf_object" "example" {
  object_id     = "lb-0ljh3xew"
  instance_id   = "waf_2kxtlbky11b2v4fe"
  member_app_id = 1306832456
  member_uin    = "100987654164"
  status        = 1
  proxy         = 1
}
```

## Argument Reference

The following arguments are supported:

* `object_id` - (Required, String, ForceNew) Modifies the object identifier.
* `instance_id` - (Optional, String) New instance ID: considered a successful modification if identical to an already bound instance.
* `ip_headers` - (Optional, Set: [`String`]) This parameter indicates a custom header and is required when `proxy` is set to 3.
* `member_app_id` - (Optional, Int, ForceNew) The ID of the member to whom the listener belongs.
* `member_uin` - (Optional, String, ForceNew) Uin of the listener member.
* `proxy` - (Optional, Int) Whether to enable proxy. 0: do not enable; 1: use the first IP address in XFF as the client IP address; 2: use remote_addr as the client IP address; 3: obtain the client IP address from the specified header field that is given in `ip_headers`.
* `status` - (Optional, Int) New WAF switch status, considered successful if identical to existing status.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Waf object can be imported using the id, e.g.

```
terraform import tencentcloud_waf_object.example lb-9h5x9lze
```

