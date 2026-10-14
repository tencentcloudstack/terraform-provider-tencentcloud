---
subcategory: "TDSQL for MySQL(DCDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dcdb_activate_hour_instance_operation"
sidebar_current: "docs-tencentcloud-resource-dcdb_activate_hour_instance_operation"
description: |-
  Provides a resource to create a dcdb activate_hour_instance_operation
---

# tencentcloud_dcdb_activate_hour_instance_operation

Provides a resource to create a dcdb activate_hour_instance_operation

## Example Usage

```hcl
resource "tencentcloud_dcdb_activate_hour_instance_operation" "activate_hour_instance_operation" {
  instance_id = local.dcdb_id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) instance ID in the format of dcdbt-ow728lmc, which can be obtained through the `DescribeDCDBInstances` API.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [ActiveHourDCDBInstance](https://cloud.tencent.com/document/api/557/70421) - This API is used to remove a pay-as-you-go TDSQL instance from isolation.

