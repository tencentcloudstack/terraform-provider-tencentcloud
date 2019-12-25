---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_scf_namespaces"
sidebar_current: "docs-tencentcloud-datasource-scf_namespaces"
description: |-
  Use this data source to query SCF namespaces.
---

# tencentcloud_scf_namespaces

Use this data source to query SCF namespaces.

## Example Usage

```hcl
resource "tencentcloud_scf_function" "foo" {
  name    = "ci-test-function"
  handler = "main.do_it"
  runtime = "Python3.6"

  cos_bucket_name   = "scf-code-1234567890"
  cos_object_name   = "code.zip"
  cos_bucket_region = "ap-guangzhou"
}

data "tencentcloud_scf_functions" "foo" {
  name = tencentcloud_scf_function.foo.name
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Description of the SCF namespace to be queried.
* `namespace` - (Optional) Name of the SCF namespace to be queried.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `namespaces` - An information list of namespace. Each element contains the following attributes:
  * `create_time` - Create time of the SCF namespace.
  * `description` - Description of the SCF namespace.
  * `modify_time` - Modify time of the SCF namespace.
  * `namespace` - Name of the SCF namespace.
  * `type` - Type of the SCF namespace.


