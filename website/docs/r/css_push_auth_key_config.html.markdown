---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_push_auth_key_config"
sidebar_current: "docs-tencentcloud-resource-css_push_auth_key_config"
description: |-
  Provides a resource to create a css push_auth_key_config
---

# tencentcloud_css_push_auth_key_config

Provides a resource to create a css push_auth_key_config

## Example Usage

```hcl
resource "tencentcloud_css_push_auth_key_config" "push_auth_key_config" {
  domain_name     = "your_push_domain_name"
  enable          = 1
  master_auth_key = "testmasterkey"
  backup_auth_key = "testbackkey"
  auth_delta      = 1800
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, String) Domain Name.
* `auth_delta` - (Optional, Int) Valid time, unit: second.
* `backup_auth_key` - (Optional, String) Standby authentication key. No transfer means that the current value is not modified.
* `enable` - (Optional, Int) Enable or not, 0: Close, 1: Enable. No transfer means that the current value is not modified.
* `master_auth_key` - (Optional, String) Primary authentication key. No transfer means that the current value is not modified.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

css push_auth_key_config can be imported using the id, e.g.

```
terraform import tencentcloud_css_push_auth_key_config.push_auth_key_config push_auth_key_config_id
```

