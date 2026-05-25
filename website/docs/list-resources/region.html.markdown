---
subcategory: "Provider Meta"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_region"
sidebar_current: "docs-tencentcloud-list_resource-region"
description: |-
  Provides a list resource that enumerates the static set of TencentCloud
regions. List resources require a companion managed resource with the
identical type name (`tencentcloud_region`) and an `IdentitySchema`. As
of the current release this list reference is shipped as an L0
placeholder: the helper data is wired up but the framework
`list.ListResource` interface is not yet implemented and the type is not
registered with the provider. Once the companion managed resource is
introduced, this document will be picked up automatically by `make doc`.
---

# tencentcloud_region

Provides a list resource that enumerates the static set of TencentCloud
regions. List resources require a companion managed resource with the
identical type name (`tencentcloud_region`) and an `IdentitySchema`. As
of the current release this list reference is shipped as an L0
placeholder: the helper data is wired up but the framework
`list.ListResource` interface is not yet implemented and the type is not
registered with the provider. Once the companion managed resource is
introduced, this document will be picked up automatically by `make doc`.

## Example Usage

```hcl
list "tencentcloud_region" "all" {
  config {}
}

output "first_region_id" {
  value = list.tencentcloud_region.all.data[0].id
}
```

## Argument Reference

The following arguments are supported:




