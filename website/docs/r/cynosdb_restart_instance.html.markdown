---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_restart_instance"
sidebar_current: "docs-tencentcloud-resource-cynosdb_restart_instance"
description: |-
  Provides a resource to create a cynosdb restart_instance
---

# tencentcloud_cynosdb_restart_instance

Provides a resource to create a cynosdb restart_instance

## Example Usage

```hcl
resource "tencentcloud_cynosdb_restart_instance" "restart_instance" {
  instance_id = "cynosdbmysql-ins-afqx1hy0"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `status` - instance state.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DescribeFlow](https://cloud.tencent.com/document/api/1003/86162) - This API is used to query task flow information.
* Create: [RestartInstance](https://cloud.tencent.com/document/api/1003/86028) - This API is used to reboot an instance.
* Read: [DescribeInstanceDetail](https://cloud.tencent.com/document/api/1003/48250) - This API is used to query instance details.
* Read: [DescribeInstances](https://cloud.tencent.com/document/api/1003/48334) - This API is used to query the list of instances.

