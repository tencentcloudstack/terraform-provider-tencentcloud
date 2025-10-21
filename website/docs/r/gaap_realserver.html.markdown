---
subcategory: "Global Application Acceleration(GAAP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_realserver"
sidebar_current: "docs-tencentcloud-resource-gaap_realserver"
description: |-
  Provides a resource to create a GAAP realserver.
---

# tencentcloud_gaap_realserver

Provides a resource to create a GAAP realserver.

## Example Usage

```hcl
resource "tencentcloud_gaap_realserver" "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"

  tags = {
    test = "test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Name of the GAAP realserver, the maximum length is 30.
* `domain` - (Optional, String, ForceNew) Domain of the GAAP realserver, conflict with `ip`.
* `ip` - (Optional, String, ForceNew) IP of the GAAP realserver, conflict with `domain`.
* `project_id` - (Optional, Int, ForceNew) ID of the project within the GAAP realserver, '0' means is default project.
* `tags` - (Optional, Map) Tags of the GAAP realserver.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

GAAP realserver can be imported using the id, e.g.

```
  $ terraform import tencentcloud_gaap_realserver.foo rs-4ftghy6
```

