---
subcategory: "Cfw"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cfw_address_template"
sidebar_current: "docs-tencentcloud-resource-cfw_address_template"
description: |-
  Provides a resource to create a cfw address_template
---

# tencentcloud_cfw_address_template

Provides a resource to create a cfw address_template

## Example Usage

### If type is 1

```hcl
resource "tencentcloud_cfw_address_template" "example" {
  name      = "tf_example"
  detail    = "test template"
  ip_string = "1.1.1.1,2.2.2.2"
  type      = 1
}
```

### If type is 5

```hcl
resource "tencentcloud_cfw_address_template" "example" {
  name      = "tf_example"
  detail    = "test template"
  ip_string = "www.qq.com,www.tencent.com"
  type      = 5
}
```

## Argument Reference

The following arguments are supported:

* `detail` - (Required, String) Template Detail.
* `ip_string` - (Required, String) Type is 1, ip template eg: 1.1.1.1,2.2.2.2; Type is 5, domain name template eg: www.qq.com, www.tencent.com.
* `name` - (Required, String) Template name.
* `type` - (Required, Int) 1: ip template; 5: domain name templates.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cfw address_template can be imported using the id, e.g.

```
terraform import tencentcloud_cfw_address_template.example mb_1300846651_1695611353900
```

