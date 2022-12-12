---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_param_template"
sidebar_current: "docs-tencentcloud-resource-redis_param_template"
description: |-
  Provides a resource to create a redis parameter template
---

# tencentcloud_redis_param_template

Provides a resource to create a redis parameter template

## Example Usage

```hcl
resource "tencentcloud_redis_param_template" "param_template" {
  name         = "example-template"
  description  = "This is an example redis param template."
  product_type = 6
  params_override {
    key   = "timeout"
    value = "7200"
  }
}
```

Copy from another template

```hcl
resource "tencentcloud_redis_param_template" "param_template" {
  name        = "example-copied"
  description = "This is an copied redis param template from xxx."
  template_id = "xxx"
  params_override {
    key   = "timeout"
    value = "7200"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Parameter template name.
* `description` - (Optional, String) Parameter template description.
* `params_override` - (Optional, List) Specify override parameter list, NOTE: Do not remove override params once set, removing will not take effects to current value.
* `product_type` - (Optional, Int) Specify product type. Valid values: 1 (Redis 2.8 Memory Edition in cluster architecture), 2 (Redis 2.8 Memory Edition in standard architecture), 3 (CKV 3.2 Memory Edition in standard architecture), 4 (CKV 3.2 Memory Edition in cluster architecture), 5 (Redis 2.8 Memory Edition in standalone architecture), 6 (Redis 4.0 Memory Edition in standard architecture), 7 (Redis 4.0 Memory Edition in cluster architecture), 8 (Redis 5.0 Memory Edition in standard architecture), 9 (Redis 5.0 Memory Edition in cluster architecture). If `template_id` is specified, this parameter can be left blank; otherwise, it is required.
* `template_id` - (Optional, String) Specify which existed template import from.

The `params_override` object supports the following:

* `key` - (Required, String) Parameter key e.g. `timeout`, check https://www.tencentcloud.com/document/product/239/39796 for more reference.
* `value` - (Required, String) Parameter value, check https://www.tencentcloud.com/document/product/239/39796 for more reference.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `param_details` - Readonly full parameter list details.
  * `current_value` - Current value.
  * `default` - Default value.
  * `description` - Parameter description.
  * `enum_value` - Enum values.
  * `max` - Maximum value.
  * `min` - Minimum value.
  * `name` - Parameter key name.
  * `need_reboot` - Indicates whether to reboot redis instance if modified.
  * `param_type` - Parameter type.


## Import

redis param_template can be imported using the id, e.g.
```
$ terraform import tencentcloud_redis_param_template.param_template param_template_id
```

