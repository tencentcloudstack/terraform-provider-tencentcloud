---
subcategory: "Intelligent Global Traffic Manager(IGTM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_igtm_instance"
sidebar_current: "docs-tencentcloud-resource-igtm_instance"
description: |-
  Provides a resource to create a IGTM instance
---

# tencentcloud_igtm_instance

Provides a resource to create a IGTM instance

~> **NOTE:** Currently, executing the `terraform destroy` command to delete this resource is not supported. If you need to destroy it, please contact Tencent Cloud IGTM through a ticket.

## Example Usage

```hcl
resource "tencentcloud_igtm_instance" "example" {
  domain            = "domain.com"
  access_type       = "CUSTOM"
  global_ttl        = 60
  package_type      = "STANDARD"
  instance_name     = "tf-example"
  access_domain     = "domain.com"
  access_sub_domain = "sub_domain.com"
  remark            = "remark."
  resource_id       = "ins-lnpnnwvwxgs"
}
```

## Argument Reference

The following arguments are supported:

* `access_domain` - (Required, String) Access main domain.
* `access_sub_domain` - (Required, String) Access subdomain.
* `access_type` - (Required, String) CUSTOM: Custom access domain
SYSTEM: System access domain.
* `domain` - (Required, String) Business domain.
* `global_ttl` - (Required, Int) Resolution effective time.
* `instance_name` - (Required, String) Instance name.
* `package_type` - (Required, String, ForceNew) Package type
FREE: Free version
STANDARD: Standard version
ULTIMATE: Ultimate version.
* `resource_id` - (Required, String, ForceNew) Package resource ID.
* `remark` - (Optional, String) Remark.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `instance_id` - Instance ID.


## Import

IGTM instance can be imported using the id, e.g.

```
terraform import tencentcloud_igtm_instance.example gtm-uukztqtoaru
```

