---
subcategory: "Serverless Cloud Function(SCF)"
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
resource "tencentcloud_scf_namespace" "foo" {
  namespace = "ci-test-scf"
}

data "tencentcloud_scf_namespaces" "foo" {
  namespace = tencentcloud_scf_namespace.foo.id
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional, String) Description of the SCF namespace to be queried.
* `namespace` - (Optional, String) Name of the SCF namespace to be queried.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `namespaces` - An information list of namespace. Each element contains the following attributes:
  * `create_time` - Create time of the SCF namespace.
  * `description` - Description of the SCF namespace.
  * `modify_time` - Modify time of the SCF namespace.
  * `namespace` - Name of the SCF namespace.
  * `type` - Type of the SCF namespace.


