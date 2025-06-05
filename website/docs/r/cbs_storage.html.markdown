---
subcategory: "Cloud Block Storage(CBS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cbs_storage"
sidebar_current: "docs-tencentcloud-resource-cbs_storage"
description: |-
  Provides a resource to create a CBS storage.
---

# tencentcloud_cbs_storage

Provides a resource to create a CBS storage.

-> **NOTE:** When creating an encrypted disk, if `kms_key_id` is not entered, the product side will generate a key by default.

-> **NOTE:** When using CBS encrypted disk, it is necessary to add `CVM_QcsRole` role and `QcloudKMSAccessForCVMRole` strategy to the account.

## Example Usage

### Create a standard CBS storage

```hcl
resource "tencentcloud_cbs_storage" "example" {
  storage_name      = "tf-example"
  storage_type      = "CLOUD_SSD"
  storage_size      = 100
  availability_zone = "ap-guangzhou-3"
  project_id        = 0
  encrypt           = false

  tags = {
    createBy = "Terraform"
  }
}
```

### Create an encrypted CBS storage with customize kms_key_id

```hcl
resource "tencentcloud_cbs_storage" "example" {
  storage_name      = "tf-example"
  storage_type      = "CLOUD_SSD"
  storage_size      = 100
  availability_zone = "ap-guangzhou-3"
  project_id        = 0
  kms_key_id        = "2e860789-7ef0-11ef-8d1c-5254001955d1"
  encrypt           = true

  tags = {
    createBy = "Terraform"
  }
}
```

### Create an encrypted CBS storage with default generated kms_key_id

```hcl
resource "tencentcloud_cbs_storage" "example" {
  storage_name      = "tf-example"
  storage_type      = "CLOUD_SSD"
  storage_size      = 100
  availability_zone = "ap-guangzhou-3"
  project_id        = 0
  encrypt           = true

  tags = {
    createBy = "Terraform"
  }
}
```

### Create an encrypted CBS storage with encrypt_type

```hcl
resource "tencentcloud_cbs_storage" "example" {
  storage_name      = "tf-example"
  storage_type      = "CLOUD_SSD"
  storage_size      = 100
  availability_zone = "ap-guangzhou-3"
  project_id        = 0
  encrypt           = true
  encrypt_type      = "ENCRYPT_V2"

  tags = {
    createBy = "Terraform"
  }
}
```

### Create a dedicated cluster CBS storage

```hcl
resource "tencentcloud_cbs_storage" "example" {
  storage_name         = "tf-example"
  storage_type         = "CLOUD_SSD"
  storage_size         = 100
  availability_zone    = "ap-guangzhou-4"
  dedicated_cluster_id = "cluster-262n63e8"
  charge_type          = "DEDICATED_CLUSTER_PAID"
  project_id           = 0
  encrypt              = false

  tags = {
    createBy = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required, String, ForceNew) The available zone that the CBS instance locates at.
* `storage_name` - (Required, String) Name of CBS. The maximum length can not exceed 60 bytes.
* `storage_size` - (Required, Int) Volume of CBS, and unit is GB.
* `storage_type` - (Required, String, ForceNew) Type of CBS medium. Valid values: CLOUD_BASIC: HDD cloud disk, CLOUD_PREMIUM: Premium Cloud Storage, CLOUD_BSSD: General Purpose SSD, CLOUD_SSD: SSD, CLOUD_HSSD: Enhanced SSD, CLOUD_TSSD: Tremendous SSD.
* `burst_performance` - (Optional, Bool) Whether to enable performance burst when creating a cloud disk.
* `charge_type` - (Optional, String) The charge type of CBS instance. Valid values are `PREPAID`, `POSTPAID_BY_HOUR`, `CDCPAID` and `DEDICATED_CLUSTER_PAID`. The default is `POSTPAID_BY_HOUR`.
* `dedicated_cluster_id` - (Optional, String, ForceNew) Exclusive cluster id.
* `disk_backup_quota` - (Optional, Int) The quota of backup points of cloud disk.
* `encrypt_type` - (Optional, String, ForceNew) Specifies the cloud disk encryption type. The values are `ENCRYPT_V1` and `ENCRYPT_V2`, which represent the first-generation and second-generation encryption technologies respectively. The two encryption technologies are incompatible with each other. It is recommended to use the second-generation encryption technology `ENCRYPT_V2` first. The first-generation encryption technology is only supported on some older models. This parameter is only valid when creating an encrypted cloud disk.
* `encrypt` - (Optional, Bool, ForceNew) Pass in this parameter to create an encrypted cloud disk.
* `force_delete` - (Optional, Bool) Indicate whether to delete CBS instance directly or not. Default is false. If set true, the instance will be deleted instead of staying recycle bin.
* `kms_key_id` - (Optional, String, ForceNew) Optional parameters. When purchasing an encryption disk, customize the key. When this parameter is passed in, the `encrypt` parameter need be set.
* `period` - (Optional, Int, **Deprecated**) It has been deprecated from version 1.33.0. Set `prepaid_period` instead. The purchased usage period of CBS. Valid values: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36].
* `prepaid_period` - (Optional, Int) The tenancy (time unit is month) of the prepaid instance, NOTE: it only works when charge_type is set to `PREPAID`. Valid values are 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36.
* `prepaid_renew_flag` - (Optional, String) Auto Renewal flag. Value range: `NOTIFY_AND_AUTO_RENEW`: Notify expiry and renew automatically, `NOTIFY_AND_MANUAL_RENEW`: Notify expiry but do not renew automatically, `DISABLE_NOTIFY_AND_MANUAL_RENEW`: Neither notify expiry nor renew automatically. Default value range: `NOTIFY_AND_MANUAL_RENEW`: Notify expiry but do not renew automatically. NOTE: it only works when charge_type is set to `PREPAID`.
* `project_id` - (Optional, Int) ID of the project to which the instance belongs.
* `snapshot_id` - (Optional, String) ID of the snapshot. If specified, created the CBS by this snapshot.
* `tags` - (Optional, Map) The available tags within this CBS.
* `throughput_performance` - (Optional, Int) Add extra performance to the data disk. Only works when disk type is `CLOUD_TSSD` or `CLOUD_HSSD`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `attached` - Indicates whether the CBS is mounted the CVM.
* `storage_status` - Status of CBS. Valid values: UNATTACHED, ATTACHING, ATTACHED, DETACHING, EXPANDING, ROLLBACKING, TORECYCLE and DUMPING.


## Import

CBS storage can be imported using the id, e.g.

```
$ terraform import tencentcloud_cbs_storage.example disk-41s6jwy4
```

