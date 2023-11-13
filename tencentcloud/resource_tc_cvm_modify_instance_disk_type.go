/*
Provides a resource to create a cvm modify_instance_disk_type

Example Usage

```hcl
resource "tencentcloud_cvm_modify_instance_disk_type" "modify_instance_disk_type" {
  instance_id = "ins-r8hr2upy"
  data_disks {
		disk_size = 50
		disk_type = "CLOUD_BASIC"
		disk_id = "disk-hrsd0u81"
		delete_with_instance = true
		snapshot_id = "snap-r9unnd89"
		encrypt = false
		kms_key_id = "kms-abcd1234"
		throughput_performance = 2
		cdc_id = "cdc-b9pbd3px"

  }
  system_disk {
		disk_type = "CLOUD_PREMIUM"
		disk_id = "disk-1drr53sd"
		disk_size = 50
		cdc_id = "cdc-b9pbd3px"

  }
}
```

Import

cvm modify_instance_disk_type can be imported using the id, e.g.

```
terraform import tencentcloud_cvm_modify_instance_disk_type.modify_instance_disk_type modify_instance_disk_type_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudCvmModifyInstanceDiskType() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmModifyInstanceDiskTypeCreate,
		Read:   resourceTencentCloudCvmModifyInstanceDiskTypeRead,
		Delete: resourceTencentCloudCvmModifyInstanceDiskTypeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID. To obtain the instance IDs, you can call DescribeInstances and look for InstanceId in the response.",
			},

			"data_disks": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "For instance data disk configuration information, you only need to specify the media type of the target cloud disk to be converted, and specify the value of DiskType. Currently, only one data disk conversion is supported. The CdcId parameter is only supported for instances of the CDHPAID type.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_size": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Data disk size (in GB). The minimum adjustment increment is 10 GB. The value range varies by data disk type. The default value is 0, indicating that no data disk is purchased. For more information, see the product documentation.",
						},
						"disk_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Data disk type. Valid values：&amp;lt;br&amp;gt;&amp;lt;li&amp;gt;LOCAL_BASIC：local disk&amp;lt;br&amp;gt;&amp;lt;li&amp;gt;LOCAL_SSD：local SSD disk&amp;lt;br&amp;gt;&amp;lt;li&amp;gt;LOCAL_NVME：local NVME disk, specified in the InstanceType&amp;lt;br&amp;gt;&amp;lt;li&amp;gt;local HDD disk, specified in the InstanceType&amp;lt;br&amp;gt;&amp;lt;li&amp;gt;CLOUD_BASIC：HDD cloud disk&amp;lt;br&amp;gt;&amp;lt;li&amp;gt;CLOUD_PREMIUM：Premium Cloud Storage&amp;lt;br&amp;gt;&amp;lt;li&amp;gt;SSD&amp;lt;br&amp;gt;&amp;lt;li&amp;gt;Enhanced SSD&amp;lt;br&amp;gt;&amp;lt;li&amp;gt;Tremendous SSD&amp;lt;br&amp;gt;&amp;lt;li&amp;gt;CLOUD_BSSD：Balanced SSD&amp;lt;br&amp;gt;&amp;lt;br&amp;gt;Default value：LOCAL_BASIC。&amp;lt;br&amp;gt;&amp;lt;br&amp;gt;This parameter is invalid for the `ResizeInstanceDisk` API.",
						},
						"disk_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Data disk ID. Note that it’s not available for LOCAL_BASIC and LOCAL_SSD disks.It is only used as a response parameter for APIs such as DescribeInstances, and cannot be used as a request parameter for APIs such as RunInstances.",
						},
						"delete_with_instance": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to terminate the data disk when its CVM is terminated. Valid values: &amp;lt;li&amp;gt;TRUE: terminate the data disk when its CVM is terminated. This value only supports pay-as-you-go cloud disks billed on an hourly basis. &amp;lt;li&amp;gt;FALSE: retain the data disk when its CVM is terminated.&amp;lt;br&amp;gt; Default value: TRUE&amp;lt;br&amp;gt; Currently this parameter is only used in the RunInstances API.Note: This field may return null, indicating that no valid value is found.",
						},
						"snapshot_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Data disk snapshot ID. The size of the selected data disk snapshot must be smaller than that of the data disk. Note: This field may return null, indicating that no valid value is found.",
						},
						"encrypt": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Specifies whether the data disk is encrypted. Valid values: &amp;lt;li&amp;gt;TRUE：encrypted &amp;lt;li&amp;gt;FALSE：not encrypted&amp;lt;br&amp;gt; Default value: FALSE&amp;lt;br&amp;gt; This parameter is only used with RunInstances.Note: this field may return null, indicating that no valid value is obtained.",
						},
						"kms_key_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of the custom CMK in the format of UUID or “kms-abcd1234”. This parameter is used to encrypt cloud disks. Currently, this parameter is only used in the RunInstances API.Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"throughput_performance": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Cloud disk performance, in MB/s Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"cdc_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of the dedicated cluster to which the instance belongs. Note: this field may return null, indicating that no valid values can be obtained.",
						},
					},
				},
			},

			"system_disk": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "For instance system disk configuration information, you only need to specify the nature type of the target cloud disk to be converted, and specify the value of DiskType. Only CDHPAID type instances are supported to specify Cd.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "System disk type. Valid values:&amp;lt;br&amp;gt;&amp;lt;li&amp;gt;LOCAL_BASIC: local disk&amp;lt;br&amp;gt;&amp;lt;li&amp;gt;LOCAL_SSD：local SSD disk&amp;lt;br&amp;gt;&amp;lt;li&amp;gt;CLOUD_BASIC：HDD cloud disk&amp;lt;br&amp;gt;&amp;lt;li&amp;gt;CLOUD_SSD：SSD cloud disk&amp;lt;br&amp;gt;&amp;lt;li&amp;gt;CLOUD_PREMIUM：Premium cloud storage&amp;lt;br&amp;gt;&amp;lt;li&amp;gt;CLOUD_BSSD：Balanced SSD&amp;lt;br&amp;gt;&amp;lt;br&amp;gt;The disk currently in stock will be used by default.",
						},
						"disk_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "|System disk ID. System disks whose type is LOCAL_BASIC or LOCAL_SSD do not have an ID and do not support this parameter. It is only used as a response parameter for APIs such as DescribeInstances, and cannot be used as a request parameter for APIs such as RunInstances.",
						},
						"disk_size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "System disk size; unit: GB; default value: 50 GB.",
						},
						"cdc_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of the dedicated cluster to which the instance belongs.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCvmModifyInstanceDiskTypeCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_modify_instance_disk_type.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = cvm.NewModifyInstanceDiskTypeRequest()
		response   = cvm.NewModifyInstanceDiskTypeResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("data_disks"); ok {
		for _, item := range v.([]interface{}) {
			dataDisk := cvm.DataDisk{}
			if v, ok := dMap["disk_size"]; ok {
				dataDisk.DiskSize = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["disk_type"]; ok {
				dataDisk.DiskType = helper.String(v.(string))
			}
			if v, ok := dMap["disk_id"]; ok {
				dataDisk.DiskId = helper.String(v.(string))
			}
			if v, ok := dMap["delete_with_instance"]; ok {
				dataDisk.DeleteWithInstance = helper.Bool(v.(bool))
			}
			if v, ok := dMap["snapshot_id"]; ok {
				dataDisk.SnapshotId = helper.String(v.(string))
			}
			if v, ok := dMap["encrypt"]; ok {
				dataDisk.Encrypt = helper.Bool(v.(bool))
			}
			if v, ok := dMap["kms_key_id"]; ok {
				dataDisk.KmsKeyId = helper.String(v.(string))
			}
			if v, ok := dMap["throughput_performance"]; ok {
				dataDisk.ThroughputPerformance = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["cdc_id"]; ok {
				dataDisk.CdcId = helper.String(v.(string))
			}
			request.DataDisks = append(request.DataDisks, &dataDisk)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "system_disk"); ok {
		systemDisk := cvm.SystemDisk{}
		if v, ok := dMap["disk_type"]; ok {
			systemDisk.DiskType = helper.String(v.(string))
		}
		if v, ok := dMap["disk_id"]; ok {
			systemDisk.DiskId = helper.String(v.(string))
		}
		if v, ok := dMap["disk_size"]; ok {
			systemDisk.DiskSize = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["cdc_id"]; ok {
			systemDisk.CdcId = helper.String(v.(string))
		}
		request.SystemDisk = &systemDisk
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCvmClient().ModifyInstanceDiskType(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cvm modifyInstanceDiskType failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudCvmModifyInstanceDiskTypeRead(d, meta)
}

func resourceTencentCloudCvmModifyInstanceDiskTypeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_modify_instance_disk_type.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCvmModifyInstanceDiskTypeDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_modify_instance_disk_type.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
