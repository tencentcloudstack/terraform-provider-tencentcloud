---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_realservers"
sidebar_current: "docs-tencentcloud-datasource-gaap_realservers"
description: |-
  Use this data source to query gaap realservers.
---

# tencentcloud_gaap_realservers

Use this data source to query gaap realservers.

## Example Usage

```hcl
resource "tencentcloud_gaap_realserver" "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"
}

data "tencentcloud_gaap_realservers" "foo" {
  ip = "${tencentcloud_gaap_realserver.foo.ip}"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Optional) Domain of the GAAP realserver to be queried, conflict with `ip`.
* `ip` - (Optional) IP of the GAAP realserver to be queried, conflict with `domain`.
* `name` - (Optional) Name of the GAAP realserver to be queried, the maximum length is 30.
* `project_id` - (Optional) ID of the project within the GAAP realserver to be queried, default is '-1' means all projects.
* `result_output_file` - (Optional) Used to save results.
* `tags` - (Optional) Tags of the GAAP proxy to be queried. Support up to 5, display the information as long as it matches one.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `realservers` - An information list of GAAP realserver. Each element contains the following attributes:
  * `domain` - Domain of the GAAP realserver.
  * `id` - ID of the GAAP realserver.
  * `ip` - IP of the GAAP realserver.
  * `name` - Name of the GAAP realserver.
  * `project_id` - ID of the project within the GAAP realserver.
  * `tags` - Tags of the GAAP realserver.


