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

* `dst_database_type` - (Required, String, ForceNew) Destination database type, such as `mysql`, `mariadb`, `percona`, `cynosdbmysql` (TDSQL-C MySQL), `tdpg` (TDSQL for PostgreSQL), `tdsqlmysql`, `kafka`, `tdstore` (TDSQL TDStore), etc.
* `dst_region` - (Required, String, ForceNew) The region where the destination database resides, such as `ap-guangzhou`.
* `pay_mode` - (Required, String, ForceNew) Billing type. Valid values: `PrePay` (subscription, monthly/yearly billing), `PostPay` (pay-as-you-go).
* `src_database_type` - (Required, String, ForceNew) Source database type, such as `mysql`, `mariadb`, `percona`, `postgresql`, `cynosdbmysql` (TDSQL-C MySQL), `tdpg` (TDSQL for PostgreSQL), `tdsqlmysql`, `tdstore` (TDSQL TDStore), etc.
* `src_region` - (Required, String, ForceNew) The region where the source database resides, such as `ap-guangzhou`.
* `auto_renew` - (Optional, Int, ForceNew) Auto-renewal flag. Only takes effect when `pay_mode` is `PrePay`. Valid values: `1` (enable auto-renewal), `0` (disable auto-renewal, default).
* `existed_job_id` - (Optional, String, ForceNew) The existing sync job ID used to create a similar job.
* `instance_class` - (Optional, String) Sync link specification, such as `micro`, `small`, `medium`, `large`. Default is `medium`.
* `job_name` - (Optional, String, ForceNew) Sync job name.
* `specification` - (Optional, String, ForceNew) Sync job specification. `Standard` indicates the standard edition; currently only `Standard` is supported.
* `tags` - (Optional, List, ForceNew) Tag information.

The `tags` object supports the following:

* `tag_key` - (Optional, String) Tag key.
* `tag_value` - (Optional, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `job_id` - Sync job ID.


## Import

DTS sync job can be imported using the id, e.g.

```
terraform import tencentcloud_dts_sync_job.example sync-hpb214ua
```

