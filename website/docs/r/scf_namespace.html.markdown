---
subcategory: "Serverless Cloud Function(SCF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_scf_namespace"
sidebar_current: "docs-tencentcloud-resource-scf_namespace"
description: |-
  Provide a resource to create a SCF namespace.
---

# tencentcloud_scf_namespace

Provide a resource to create a SCF namespace.

## Example Usage

```hcl
resource "tencentcloud_scf_namespace" "foo" {
  namespace = "ci-test-scf"
}
```

## Argument Reference

The following arguments are supported:

* `namespace` - (Required, String, ForceNew) Name of the SCF namespace.
* `description` - (Optional, String) Description of the SCF namespace.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - SCF namespace creation time.
* `modify_time` - SCF namespace last modified time.
* `type` - SCF namespace type.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateNamespace](https://www.tencentcloud.com/document/api/583/34413)
* Read: [ListNamespaces](https://www.tencentcloud.com/document/api/583/34411)
* Update: [UpdateNamespace](https://www.tencentcloud.com/document/api/583/38339)
* Delete: [DeleteNamespace](https://www.tencentcloud.com/document/api/583/34412)

## Import

SCF namespace can be imported, e.g.

```
$ terraform import tencentcloud_scf_function.test default
```

