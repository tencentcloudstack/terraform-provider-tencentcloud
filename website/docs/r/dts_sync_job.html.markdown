---
subcategory: "Data Transmission Service(DTS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dts_sync_job"
sidebar_current: "docs-tencentcloud-resource-dts_sync_job"
description: |-
  Provides a resource to create a DTS sync job
---

# tencentcloud_dts_sync_job

Provides a resource to create a DTS sync job

~> **NOTE:** Import function does not support field `existed_job_id`.

## Example Usage

```hcl
resource "tencentcloud_dts_sync_job" "example" {
  pay_mode          = "PostPay"
  src_database_type = "mysql"
  src_region        = "ap-guangzhou"
  dst_database_type = "cynosdbmysql"
  dst_region        = "ap-guangzhou"
  auto_renew        = 0
  instance_class    = "micro"
  tags {
    tag_key   = "createBy"
    tag_value = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `dst_database_type` - (Required, String, ForceNew) destination database type.
* `dst_region` - (Required, String, ForceNew) destination region.
* `pay_mode` - (Required, String, ForceNew) pay mode, optional value is PrePay or PostPay.
* `src_database_type` - (Required, String, ForceNew) source database type.
* `src_region` - (Required, String, ForceNew) source region.
* `auto_renew` - (Optional, Int, ForceNew) auto renew.
* `existed_job_id` - (Optional, String, ForceNew) existed job id.
* `instance_class` - (Optional, String, ForceNew) instance class.
* `job_name` - (Optional, String, ForceNew) job name.
* `specification` - (Optional, String, ForceNew) specification.
* `tags` - (Optional, List, ForceNew) tags.

The `tags` object supports the following:

* `tag_key` - (Optional, String) tag key.
* `tag_value` - (Optional, String) tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `job_id` - job id.


