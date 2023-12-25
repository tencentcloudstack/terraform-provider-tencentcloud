package cvm

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCvmModifyInstanceDiskType() *schema.Resource {
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
							Type:     schema.TypeString,
							Optional: true,
							Description: "Data disk type. Valid values:\n" +
								"- LOCAL_BASIC: local hard disk;\n" +
								"- LOCAL_SSD: local SSD hard disk;\n" +
								"- LOCAL_NVME: local NVME hard disk, which is strongly related to InstanceType and cannot be specified;\n" +
								"- LOCAL_PRO: local HDD hard disk, which is strongly related to InstanceType and cannot be specified;\n" +
								"- CLOUD_BASIC: ordinary cloud disk;\n" +
								"- CLOUD_PREMIUM: high-performance cloud disk;\n" +
								"- CLOUD_SSD:SSD cloud disk;\n" +
								"- CLOUD_HSSD: enhanced SSD cloud disk;\n" +
								"- CLOUD_TSSD: extremely fast SSD cloud disk;\n" +
								"- CLOUD_BSSD: general-purpose SSD cloud disk;\n" +
								"Default value: LOCAL_BASIC.",
						},
						"disk_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Data disk ID. Note that it's not available for LOCAL_BASIC and LOCAL_SSD disks.",
						},
						"delete_with_instance": {
							Type:     schema.TypeBool,
							Optional: true,
							Description: "Whether to terminate the data disk when its CVM is terminated. Valid values:\n" +
								"- TRUE: terminate the data disk when its CVM is terminated. This value only supports pay-as-you-go cloud disks billed on an hourly basis.\n" +
								"- FALSE: retain the data disk when its CVM is terminated.\n" +
								"Default value: TRUE.",
						},
						"snapshot_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Data disk snapshot ID. The size of the selected data disk snapshot must be smaller than that of the data disk.",
						},
						"encrypt": {
							Type:     schema.TypeBool,
							Optional: true,
							Description: "Specifies whether the data disk is encrypted. Valid values:\n" +
								"- TRUE: encrypted\n" +
								"- FALSE: not encrypted\n" +
								"Default value: FALSE.",
						},
						"kms_key_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of the custom CMK in the format of UUID or “kms-abcd1234”. This parameter is used to encrypt cloud disks.",
						},
						"throughput_performance": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Cloud disk performance, in MB/s.",
						},
						"cdc_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of the dedicated cluster to which the instance belongs.",
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
							Type:     schema.TypeString,
							Optional: true,
							Description: "System disk type. Valid values:" +
								"- LOCAL_BASIC: local disk\n" +
								"- LOCAL_SSD: local SSD disk\n" +
								"- CLOUD_BASIC: ordinary cloud disk\n" +
								"- CLOUD_SSD: SSD cloud disk\n" +
								"- CLOUD_PREMIUM: Premium cloud storage\n" +
								"- CLOUD_BSSD: Balanced SSD\n" +
								"The disk currently in stock will be used by default.",
						},
						"disk_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "System disk ID. System disks whose type is LOCAL_BASIC or LOCAL_SSD do not have an ID and do not support this parameter.",
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
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_modify_instance_disk_type.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = cvm.NewModifyInstanceDiskTypeRequest()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId := v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	if v, ok := d.GetOk("data_disks"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().ModifyInstanceDiskType(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cvm modifyInstanceDiskType failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudCvmModifyInstanceDiskTypeRead(d, meta)
}

func resourceTencentCloudCvmModifyInstanceDiskTypeRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_modify_instance_disk_type.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCvmModifyInstanceDiskTypeDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_modify_instance_disk_type.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
