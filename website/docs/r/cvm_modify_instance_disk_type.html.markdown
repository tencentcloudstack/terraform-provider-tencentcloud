---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cvm_modify_instance_disk_type"
sidebar_current: "docs-tencentcloud-resource-cvm_modify_instance_disk_type"
description: |-
  Provides a resource to create a cvm modify_instance_disk_type
---

# tencentcloud_cvm_modify_instance_disk_type

Provides a resource to create a cvm modify_instance_disk_type

## Example Usage

```hcl
resource "tencentcloud_cvm_modify_instance_disk_type" "modify_instance_disk_type" {
  instance_id = "ins-r8hr2upy"
  data_disks {
    disk_size              = 50
    disk_type              = "CLOUD_BASIC"
    disk_id                = "disk-hrsd0u81"
    delete_with_instance   = true
    snapshot_id            = "snap-r9unnd89"
    encrypt                = false
    kms_key_id             = "kms-abcd1234"
    throughput_performance = 2
    cdc_id                 = "cdc-b9pbd3px"

  }
  system_disk {
    disk_type = "CLOUD_PREMIUM"
    disk_id   = "disk-1drr53sd"
    disk_size = 50
    cdc_id    = "cdc-b9pbd3px"

  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID. To obtain the instance IDs, you can call DescribeInstances and look for InstanceId in the response.
* `data_disks` - (Optional, List, ForceNew) For instance data disk configuration information, you only need to specify the media type of the target cloud disk to be converted, and specify the value of DiskType. Currently, only one data disk conversion is supported. The CdcId parameter is only supported for instances of the CDHPAID type.
* `system_disk` - (Optional, List, ForceNew) For instance system disk configuration information, you only need to specify the nature type of the target cloud disk to be converted, and specify the value of DiskType. Only CDHPAID type instances are supported to specify Cd.

The `data_disks` object supports the following:

* `disk_size` - (Required, Int) Data disk size (in GB). The minimum adjustment increment is 10 GB. The value range varies by data disk type. The default value is 0, indicating that no data disk is purchased. For more information, see the product documentation.
* `cdc_id` - (Optional, String) ID of the dedicated cluster to which the instance belongs.
* `delete_with_instance` - (Optional, Bool) Whether to terminate the data disk when its CVM is terminated. Valid values:
- TRUE: terminate the data disk when its CVM is terminated. This value only supports pay-as-you-go cloud disks billed on an hourly basis.
- FALSE: retain the data disk when its CVM is terminated.
Default value: TRUE.
* `disk_id` - (Optional, String) Data disk ID. Note that it's not available for LOCAL_BASIC and LOCAL_SSD disks.
* `disk_type` - (Optional, String) Data disk type. Valid values:
- LOCAL_BASIC: local hard disk;
- LOCAL_SSD: local SSD hard disk;
- LOCAL_NVME: local NVME hard disk, which is strongly related to InstanceType and cannot be specified;
- LOCAL_PRO: local HDD hard disk, which is strongly related to InstanceType and cannot be specified;
- CLOUD_BASIC: ordinary cloud disk;
- CLOUD_PREMIUM: high-performance cloud disk;
- CLOUD_SSD:SSD cloud disk;
- CLOUD_HSSD: enhanced SSD cloud disk;
- CLOUD_TSSD: extremely fast SSD cloud disk;
- CLOUD_BSSD: general-purpose SSD cloud disk;
Default value: LOCAL_BASIC.
* `encrypt` - (Optional, Bool) Specifies whether the data disk is encrypted. Valid values:
- TRUE: encrypted
- FALSE: not encrypted
Default value: FALSE.
* `kms_key_id` - (Optional, String) ID of the custom CMK in the format of UUID or `kms-abcd1234`. This parameter is used to encrypt cloud disks.
* `snapshot_id` - (Optional, String) Data disk snapshot ID. The size of the selected data disk snapshot must be smaller than that of the data disk.
* `throughput_performance` - (Optional, Int) Cloud disk performance, in MB/s.

The `system_disk` object supports the following:

* `cdc_id` - (Optional, String) ID of the dedicated cluster to which the instance belongs.
* `disk_id` - (Optional, String) System disk ID. System disks whose type is LOCAL_BASIC or LOCAL_SSD do not have an ID and do not support this parameter.
* `disk_size` - (Optional, Int) System disk size; unit: GB; default value: 50 GB.
* `disk_type` - (Optional, String) System disk type. Valid values:- LOCAL_BASIC: local disk
- LOCAL_SSD: local SSD disk
- CLOUD_BASIC: ordinary cloud disk
- CLOUD_SSD: SSD cloud disk
- CLOUD_PREMIUM: Premium cloud storage
- CLOUD_BSSD: Balanced SSD
The disk currently in stock will be used by default.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [ModifyInstanceDiskType](https://www.tencentcloud.com/search/?wd=ModifyInstanceDiskType&type=document&lang=en)

