---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_function_replica"
sidebar_current: "docs-tencentcloud-resource-teo_function_replica"
description: |-
  Provides a resource to create a TEO edge function replica
---

# tencentcloud_teo_function_replica

Provides a resource to create a TEO edge function replica

## Example Usage

```hcl
resource "tencentcloud_teo_function_replica" "example" {
  zone_id      = "zone-2qtuhspy7cr6"
  function_id  = "ef-2qlxy8s7o96e"
  replica_name = "replica-example"
  content      = "addEventListener('fetch', event => { event.respondWith(new Response('hello world')) })"
  remark       = "example replica"
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required, String) Edge function replica content. Currently only supports JavaScript code, maximum 5MB.
* `function_id` - (Required, String, ForceNew) Function ID.
* `replica_name` - (Required, String, ForceNew) Edge function replica name. Limited to 1-50 characters, allowed characters are a-z, 0-9, -, and - cannot be used alone or consecutively, nor at the beginning or end. Replica names must be unique under the same FunctionId.
* `zone_id` - (Required, String, ForceNew) Zone ID.
* `remark` - (Optional, String) Edge function replica description. Maximum 50 characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

TEO function replica can be imported using the zone_id#function_id#replica_name, e.g.

```
terraform import tencentcloud_teo_function_replica.example zone-2qtuhspy7cr6#ef-2qlxy8s7o96e#replica-example
```

