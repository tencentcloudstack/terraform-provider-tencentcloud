---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_isolate_instance"
sidebar_current: "docs-tencentcloud-resource-mysql_isolate_instance"
description: |-
  Provides a resource to create a mysql isolate_instance
---

# tencentcloud_mysql_isolate_instance

Provides a resource to create a mysql isolate_instance

## Example Usage

```hcl
resource "tencentcloud_mysql_isolate_instance" "isolate_instance" {
  instance_id = "cdb-1tru99al"
  operate     = "recover"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID, the format is: cdb-c1nl9rpv, which is the same as the instance ID displayed on the cloud database console page, and you can use the [query instance list] (https://cloud.tencent.com/document/api/236/15872) interface Gets the value of the field InstanceId in the output parameter.
* `operate` - (Required, String) Manipulate instance, `isolate` - isolate instance, `recover`- recover isolated instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `status` - Instance status.


