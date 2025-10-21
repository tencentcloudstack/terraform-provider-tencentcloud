---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_web_shell"
sidebar_current: "docs-tencentcloud-resource-waf_web_shell"
description: |-
  Provides a resource to create a waf web_shell
---

# tencentcloud_waf_web_shell

Provides a resource to create a waf web_shell

## Example Usage

```hcl
resource "tencentcloud_waf_web_shell" "example" {
  domain = "demo.waf.com"
  status = 0
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String) Domain.
* `status` - (Required, Int) Webshell status, 1: open; 0: closed; 2: log.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

waf web_shell can be imported using the id, e.g.

```
terraform import tencentcloud_waf_web_shell.example demo.waf.com
```

