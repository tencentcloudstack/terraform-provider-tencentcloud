---
subcategory: "Data Transmission Service(DTS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dts_migrate_job_resume_operation"
sidebar_current: "docs-tencentcloud-resource-dts_migrate_job_resume_operation"
description: |-
  Provides a resource to create a DTS migrate job resume operation
---

# tencentcloud_dts_migrate_job_resume_operation

Provides a resource to create a DTS migrate job resume operation

## Example Usage

```hcl
resource "tencentcloud_dts_migrate_job_resume_operation" "example" {
  job_id        = "dts-puwyj5uy"
  resume_option = "normal"
}
```

## Argument Reference

The following arguments are supported:

* `job_id` - (Required, String, ForceNew) job id.
* `resume_option` - (Required, String, ForceNew) resume mode: 1.clearData-Clear target data; 2.overwrite-The task is executed in overwrite mode; 3.normal-No extra action. Note that clearData and overwrite are valid only for redis links, normal is valid only for non-Redis links.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DescribeMigrationDetail](https://www.tencentcloud.com/document/product/571/51813) - This API is used to query the details of a migration task.
* Create: [ResumeMigrateJob](https://www.tencentcloud.com/document/product/571/51803) - This API is used to retry an abnormal or failed Redis data migration task. Note that this operation will skip the check stage and directly start the task, just like by calling `StartMigrationJob`.

