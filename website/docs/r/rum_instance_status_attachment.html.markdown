---
subcategory: "Real User Monitoring(RUM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_rum_instance_status_attachment"
sidebar_current: "docs-tencentcloud-resource-rum_instance_status_attachment"
description: |-
  Provides a resource to create a rum instance_status_attachment
---

# tencentcloud_rum_instance_status_attachment

Provides a resource to create a rum instance_status_attachment

## Example Usage

```hcl
resource "tencentcloud_rum_instance_status_attachment" "instance_status_attachment" {
  instance_id = "rum-pasZKEI3RLgakj"
  operate     = "stop"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID.
* `operate` - (Required, String) `resume`, `stop`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `instance_status` - Instance status (`1`=creating, `2`=running, `3`=abnormal, `4`=restarting, `5`=stopping, `6`=stopped, `7`=deleted).


## Import

rum instance_status_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_rum_instance_status_attachment.instance_status_attachment instance_id
```

