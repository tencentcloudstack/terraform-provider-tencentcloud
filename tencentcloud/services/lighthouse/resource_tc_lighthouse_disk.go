package lighthouse

import (
	"context"
	"fmt"
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudLighthouseDisk() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLighthouseDiskCreate,
		Read:   resourceTencentCloudLighthouseDiskRead,
		Update: resourceTencentCloudLighthouseDiskUpdate,
		Delete: resourceTencentCloudLighthouseDiskDelete,
		Schema: map[string]*schema.Schema{
			"zone": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Availability zone.",
			},

			"disk_size": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Disk size, unit: GB.",
			},

			"disk_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Disk type. Value:CLOUD_PREMIUM, CLOUD_SSD.",
			},

			"disk_charge_prepaid": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Disk subscription related parameter settings.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"period": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "new purchase cycle.",
						},
						"renew_flag": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Automatic renewal flag. Value: `NOTIFY_AND_AUTO_RENEW`: Notice expires and auto-renews. `NOTIFY_AND_MANUAL_RENEW`: Notification expires without automatic renewal, users need to manually renew. `DISABLE_NOTIFY_AND_AUTO_RENEW`: No automatic renewal and no notification. Default: `NOTIFY_AND_MANUAL_RENEW`. If this parameter is specified as `NOTIFY_AND_AUTO_RENEW`, the disk will be automatically renewed monthly when the account balance is sufficient.",
						},
						"time_unit": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "newly purchased unit. Default: m.",
						},
					},
				},
			},

			"disk_name": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Disk name. Maximum length 60.",
			},

			"disk_count": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Disk count. Values: [1, 30]. Default: 1.",
			},

			"disk_backup_quota": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Specify the disk backup quota. If not uploaded, the default is no backup quota. Currently, only one disk backup quota is supported.",
			},

			"auto_voucher": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether to automatically use the voucher. Not used by default.",
			},

			"auto_mount_configuration": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Automatically mount and initialize data disks.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Instance ID to be mounted. The specified instance must be in the Running state.",
						},
						"mount_point": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The mount point within the instance. Only instances of the Linux operating system can pass in this parameter, and if it is not passed, it will be mounted under the /data/disk path by default.",
						},
						"file_system_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The file system type. Value: ext4, xfs. Only instances of the Linux operating system can pass in this parameter, and if it is not passed, it defaults to ext4.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudLighthouseDiskCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_disk.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = lighthouse.NewCreateDisksRequest()
		response = lighthouse.NewCreateDisksResponse()
	)
	if v, ok := d.GetOk("zone"); ok {
		request.Zone = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("disk_size"); ok {
		request.DiskSize = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("disk_type"); ok {
		request.DiskType = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "disk_charge_prepaid"); ok {
		diskChargePrepaid := lighthouse.DiskChargePrepaid{}
		if v, ok := dMap["period"]; ok {
			diskChargePrepaid.Period = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["renew_flag"]; ok {
			diskChargePrepaid.RenewFlag = helper.String(v.(string))
		}
		if v, ok := dMap["time_unit"]; ok {
			diskChargePrepaid.TimeUnit = helper.String(v.(string))
		}
		request.DiskChargePrepaid = &diskChargePrepaid
	}

	if v, ok := d.GetOk("disk_name"); ok {
		request.DiskName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("disk_count"); ok {
		request.DiskCount = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("disk_backup_quota"); ok {
		request.DiskBackupQuota = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("auto_voucher"); ok {
		request.AutoVoucher = helper.Bool(v.(bool))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "auto_mount_configuration"); ok {
		autoMountConfiguration := lighthouse.AutoMountConfiguration{}
		if v, ok := dMap["instance_id"]; ok {
			autoMountConfiguration.InstanceId = helper.String(v.(string))
		}
		if v, ok := dMap["mount_point"]; ok {
			autoMountConfiguration.MountPoint = helper.String(v.(string))
		}
		if v, ok := dMap["file_system_type"]; ok {
			autoMountConfiguration.FileSystemType = helper.String(v.(string))
		}
		request.AutoMountConfiguration = &autoMountConfiguration
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseLighthouseClient().CreateDisks(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create lighthouse disk failed, reason:%+v", logId, err)
		return err
	}

	if response == nil || response.Response == nil || len(response.Response.DiskIdSet) == 0 {
		return fmt.Errorf("Response is null")
	}
	diskIds := make([]string, 0, len(response.Response.DiskIdSet))
	for _, diskId := range response.Response.DiskIdSet {
		diskIds = append(diskIds, *diskId)
	}

	d.SetId(helper.IdFormat(diskIds...))

	service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"UNATTACHED", "ATTACHED"}, 20*tccommon.ReadRetryTimeout, time.Second, service.LighthouseDiskStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudLighthouseDiskRead(d, meta)
}

func resourceTencentCloudLighthouseDiskRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_disk.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	diskIds := helper.IdParse(d.Id())
	diskId := diskIds[0]
	disk, err := service.DescribeLighthouseDiskById(ctx, diskId)
	if err != nil {
		return err
	}

	if disk == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `LighthouseDisk` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if disk.Zone != nil {
		_ = d.Set("zone", disk.Zone)
	}

	if disk.DiskSize != nil {
		_ = d.Set("disk_size", disk.DiskSize)
	}

	if disk.DiskType != nil {
		_ = d.Set("disk_type", disk.DiskType)
	}

	if disk.DiskName != nil {
		_ = d.Set("disk_name", disk.DiskName)
	}

	if disk.DiskBackupQuota != nil {
		_ = d.Set("disk_backup_quota", disk.DiskBackupQuota)
	}
	return nil
}

func resourceTencentCloudLighthouseDiskUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_disk.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	diskId := d.Id()
	request := lighthouse.NewModifyDisksAttributeRequest()
	request.DiskIds = []*string{&diskId}
	changeAttribute := false

	if d.HasChange("disk_name") {
		if v, ok := d.GetOk("disk_name"); ok {
			request.DiskName = helper.String(v.(string))
			changeAttribute = true
		}
	}
	if changeAttribute {
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseLighthouseClient().ModifyDisksAttribute(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update lighthouse disk failed, reason:%+v", logId, err)
			return err
		}
	}

	renewFlagRequest := lighthouse.NewModifyDisksRenewFlagRequest()
	renewFlagRequest.DiskIds = []*string{&diskId}
	changeRenewFlag := false
	if d.HasChange("disk_charge_prepaid.0.renew_flag") {
		if v, ok := d.GetOk("disk_charge_prepaid.0.renew_flag"); ok {
			renewFlagRequest.RenewFlag = helper.String(v.(string))
			changeRenewFlag = true
		}
	}
	if changeRenewFlag {
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseLighthouseClient().ModifyDisksRenewFlag(renewFlagRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update lighthouse disk failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudLighthouseDiskRead(d, meta)
}

func resourceTencentCloudLighthouseDiskDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_disk.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	diskId := d.Id()

	if err := service.IsolateLighthouseDiskById(ctx, diskId); err != nil {
		return err
	}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"SUCCESS"}, 20*tccommon.ReadRetryTimeout, time.Second, service.LighthouseDiskIsolateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	if err := service.TerminateLighthouseDiskById(ctx, diskId); err != nil {
		return err
	}

	return nil
}
