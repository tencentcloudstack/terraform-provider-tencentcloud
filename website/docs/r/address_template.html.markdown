---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_address_template"
sidebar_current: "docs-tencentcloud-resource-address_template"
description: |-
  Provides a resource to manage address template.
---

# tencentcloud_address_template

Provides a resource to manage address template.

~> **NOTE:** It can be replaced by `tencentcloud_address_extra_template`.

## Example Usage

```hcl
resource "tencentcloud_address_template" "foo" {
  name      = "cam-user-test"
  addresses = ["10.0.0.1", "10.0.1.0/24", "10.0.0.1-10.0.0.100"]
}
```

## Argument Reference

The following arguments are supported:

* `addresses` - (Required, Set: [`String`]) Address list. IP(`10.0.0.1`), CIDR(`10.0.1.0/24`), IP range(`10.0.0.1-10.0.0.100`) format are supported.
* `name` - (Required, String) Name of the address template.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateAddressTemplate](https://cloud.tencent.com/document/api/215/16708)
* Read: [DescribeAddressTemplates](https://cloud.tencent.com/document/api/215/16717)
* Update: [ModifyAddressTemplateAttribute](https://cloud.tencent.com/document/api/215/16720)
* Delete: [DeleteAddressTemplate](https://cloud.tencent.com/document/api/215/16712)
* Delete: [DescribeAddressTemplates](https://cloud.tencent.com/document/api/215/16717)

## Import

Address template can be imported using the id, e.g.

```
$ terraform import tencentcloud_address_template.foo ipm-makf7k9e"
```

