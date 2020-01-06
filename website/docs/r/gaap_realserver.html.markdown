---
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

* `name` - (Required) Name of the GAAP realserver, the maximum length is 30.
* `domain` - (Optional, ForceNew) Domain of the GAAP realserver, conflict with `ip`.
* `ip` - (Optional, ForceNew) IP of the GAAP realserver, conflict with `domain`.
* `project_id` - (Optional, ForceNew) ID of the project within the GAAP realserver, '0' means is default project.
* `tags` - (Optional) Tags of the GAAP realserver.


## Import

GAAP realserver can be imported using the id, e.g.

```
  $ terraform import tencentcloud_gaap_realserver.foo rs-4ftghy6
```

