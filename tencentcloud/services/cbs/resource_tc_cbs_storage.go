package cbs

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCbsStorage() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCbsStorageCreate,
		Read:   resourceTencentCloudCbsStorageRead,
		Update: resourceTencentCloudCbsStorageUpdate,
		Delete: resourceTencentCloudCbsStorageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"storage_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Type of CBS medium. Valid values: CLOUD_BASIC: HDD cloud disk, CLOUD_PREMIUM: Premium Cloud Storage, CLOUD_BSSD: General Purpose SSD, CLOUD_SSD: SSD, CLOUD_HSSD: Enhanced SSD, CLOUD_TSSD: Tremendous SSD.",
			},
			"storage_size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Volume of CBS, and unit is GB.",
			},
			"period": {
				Deprecated:   "It has been deprecated from version 1.33.0. Set `prepaid_period` instead.",
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(1, 36),
				Description:  "The purchased usage period of CBS. Valid values: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36].",
			},
			"charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      CBS_CHARGE_TYPE_POSTPAID,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CBS_CHARGE_TYPE),
				Description:  "The charge type of CBS instance. Valid values are `PREPAID` and `POSTPAID_BY_HOUR`. The default is `POSTPAID_BY_HOUR`.",
			},
			"prepaid_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue(CBS_PREPAID_PERIOD),
				Description:  "The tenancy (time unit is month) of the prepaid instance, NOTE: it only works when charge_type is set to `PREPAID`. Valid values are 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36.",
			},
			"prepaid_renew_flag": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CBS_PREPAID_RENEW_FLAG),
				Description:  "Auto Renewal flag. Value range: `NOTIFY_AND_AUTO_RENEW`: Notify expiry and renew automatically, `NOTIFY_AND_MANUAL_RENEW`: Notify expiry but do not renew automatically, `DISABLE_NOTIFY_AND_MANUAL_RENEW`: Neither notify expiry nor renew automatically. Default value range: `NOTIFY_AND_MANUAL_RENEW`: Notify expiry but do not renew automatically. NOTE: it only works when charge_type is set to `PREPAID`.",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The available zone that the CBS instance locates at.",
			},
			"storage_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(2, 60),
				Description:  "Name of CBS. The maximum length can not exceed 60 bytes.",
			},
			"snapshot_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the snapshot. If specified, created the CBS by this snapshot.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "ID of the project to which the instance belongs.",
			},
			"encrypt": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Indicates whether CBS is encrypted.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The available tags within this CBS.",
			},
			"force_delete": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicate whether to delete CBS instance directly or not. Default is false. If set true, the instance will be deleted instead of staying recycle bin.",
			},
			"throughput_performance": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Add extra performance to the data disk. Only works when disk type is `CLOUD_TSSD` or `CLOUD_HSSD`.",
			},
			"disk_backup_quota": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The quota of backup points of cloud disk.",
			},
			// computed
			"storage_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of CBS. Valid values: UNATTACHED, ATTACHING, ATTACHED, DETACHING, EXPANDING, ROLLBACKING, TORECYCLE and DUMPING.",
			},
			"attached": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the CBS is mounted the CVM.",
			},
		},
	}
}

func resourceTencentCloudCbsStorageCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_storage.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cbsService := CbsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	request := cbs.NewCreateDisksRequest()
	request.DiskName = helper.String(d.Get("storage_name").(string))
	request.DiskType = helper.String(d.Get("storage_type").(string))
	request.DiskSize = helper.IntUint64(d.Get("storage_size").(int))
	request.Placement = &cbs.Placement{
		Zone: helper.String(d.Get("availability_zone").(string)),
	}
	if v, ok := d.GetOk("project_id"); ok {
		request.Placement.ProjectId = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOk("snapshot_id"); ok {
		request.SnapshotId = helper.String(v.(string))
	}
	if _, ok := d.GetOk("encrypt"); ok {
		request.Encrypt = helper.String("ENCRYPT")
	}

	if v, ok := d.GetOk("throughput_performance"); ok {
		request.ThroughputPerformance = helper.IntUint64(v.(int))
	}

	chargeType := d.Get("charge_type").(string)

	request.DiskChargeType = &chargeType

	if chargeType == CBS_CHARGE_TYPE_PREPAID {
		request.DiskChargePrepaid = &cbs.DiskChargePrepaid{}

		if period, ok := d.GetOk("prepaid_period"); ok {
			periodInt64 := uint64(period.(int))
			request.DiskChargePrepaid.Period = &periodInt64
		} else {
			return fmt.Errorf("CBS instance charge type prepaid period can not be empty when charge type is %s",
				chargeType)
		}
		if renewFlag, ok := d.GetOk("prepaid_renew_flag"); ok {
			request.DiskChargePrepaid.RenewFlag = helper.String(renewFlag.(string))
		}
	}

	if v := helper.GetTags(d, "tags"); len(v) > 0 {
		for tagKey, tagValue := range v {
			tag := cbs.Tag{
				Key:   helper.String(tagKey),
				Value: helper.String(tagValue),
			}
			request.Tags = append(request.Tags, &tag)
		}
	}

	storageId := ""
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCbsClient().CreateDisks(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return tccommon.RetryError(e, tccommon.InternalError)
		}

		if len(response.Response.DiskIdSet) < 1 {
			return resource.NonRetryableError(fmt.Errorf("storage id is nil"))
		}

		storageId = *response.Response.DiskIdSet[0]
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cbs failed, reason:%s\n ", logId, err.Error())
		return err
	}
	d.SetId(storageId)

	if v, ok := d.GetOk("disk_backup_quota"); ok {
		err = cbsService.ModifyDiskBackupQuota(ctx, storageId, v.(int))
		if err != nil {
			return err
		}
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := &TagService{client: tcClient}
		resourceName := tccommon.BuildTagResourceName("cvm", "volume", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	// must wait for finishing creating disk
	err = resource.Retry(10*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		storage, e := cbsService.DescribeDiskById(ctx, storageId)
		if e != nil {
			return tccommon.RetryError(e, tccommon.InternalError)
		}
		if storage == nil {
			return resource.RetryableError(fmt.Errorf("storage is still creating..."))
		}
		return nil
	})
	if err != nil {
		return err
	}

	return resourceTencentCloudCbsStorageRead(d, meta)
}

func resourceTencentCloudCbsStorageRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_storage.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	storageId := d.Id()
	cbsService := CbsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	var storage *cbs.Disk
	var e error
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		storage, e = cbsService.DescribeDiskById(ctx, storageId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read cbs failed, reason:%s\n ", logId, err.Error())
		return err
	}
	if storage == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("storage_type", storage.DiskType)
	_ = d.Set("storage_size", storage.DiskSize)
	_ = d.Set("availability_zone", storage.Placement.Zone)
	_ = d.Set("storage_name", storage.DiskName)
	_ = d.Set("project_id", storage.Placement.ProjectId)
	_ = d.Set("encrypt", storage.Encrypt)
	_ = d.Set("storage_status", storage.DiskState)
	_ = d.Set("attached", storage.Attached)
	_ = d.Set("charge_type", storage.DiskChargeType)
	_ = d.Set("prepaid_renew_flag", storage.RenewFlag)
	_ = d.Set("throughput_performance", storage.ThroughputPerformance)

	if *storage.DiskChargeType == CBS_CHARGE_TYPE_PREPAID {
		_ = d.Set("prepaid_renew_flag", storage.RenewFlag)
	}

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "cvm", "volume", tcClient.Region, d.Id())
	if err != nil {
		return err
	}

	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudCbsStorageUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_storage.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	//only support update prepaid_period when upgrade chargeType
	if d.HasChange("prepaid_period") && (!d.HasChange("charge_type") && d.Get("charge_type").(string) == CBS_CHARGE_TYPE_PREPAID) {
		return fmt.Errorf("tencentcloud_cbs_storage renew is not support yet")
	}
	if d.HasChange("charge_type") && d.Get("charge_type").(string) != CBS_CHARGE_TYPE_PREPAID {
		return fmt.Errorf("tencentcloud_cbs_storage do not support downgrade instance")
	}

	cbsService := CbsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	d.Partial(true)
	storageId := d.Id()
	storageName := ""
	projectId := -1
	changed := false

	if d.HasChange("storage_name") {
		changed = true
		storageName = d.Get("storage_name").(string)
	}

	if d.HasChange("project_id") {
		changed = true
		projectId = d.Get("project_id").(int)
	}

	if changed {
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			e := cbsService.ModifyDiskAttributes(ctx, storageId, storageName, projectId)
			if e != nil {
				return tccommon.RetryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update cbs failed, reason:%s\n ", logId, err.Error())
			return err
		}
	}

	if d.HasChange("storage_size") {
		oldInterface, newInterface := d.GetChange("storage_size")
		oldValue := oldInterface.(int)
		newValue := newInterface.(int)
		if oldValue > newValue {
			return fmt.Errorf("storage size must be greater than current storage size")
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			e := cbsService.ResizeDisk(ctx, storageId, newValue)
			if e != nil {
				return tccommon.RetryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update cbs failed, reason:%s\n ", logId, err.Error())
			return err
		}

		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			storage, e := cbsService.DescribeDiskById(ctx, storageId)
			if e != nil {
				return tccommon.RetryError(e)
			}

			if *storage.DiskState == CBS_STORAGE_STATUS_EXPANDING {
				return resource.RetryableError(fmt.Errorf("cbs storage status is %s", *storage.DiskState))
			}
			if *storage.DiskSize != uint64(newValue) {
				return resource.RetryableError(fmt.Errorf("waiting for cbs size changed to %d, now %d", newValue, *storage.DiskSize))
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update cbs failed, reason:%s\n ", logId, err.Error())
			return err
		}

	}

	if d.HasChange("snapshot_id") {
		snapshotId := d.Get("snapshot_id").(string)
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			e := cbsService.ApplySnapshot(ctx, storageId, snapshotId)
			if e != nil {
				return tccommon.RetryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update cbs failed, reason:%s\n ", logId, err.Error())
			return err
		}

		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			storage, e := cbsService.DescribeDiskById(ctx, storageId)
			if e != nil {
				return tccommon.RetryError(e)
			}
			if storage != nil && *storage.DiskState == CBS_STORAGE_STATUS_ROLLBACKING {
				return resource.RetryableError(fmt.Errorf("cbs storage status is %s", *storage.DiskState))
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update cbs failed, reason:%s\n ", logId, err.Error())
			return err
		}

	}

	if d.HasChange("throughput_performance") {
		throughputPerformance := d.Get("throughput_performance").(int)
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			e := cbsService.ModifyThroughputPerformance(ctx, storageId, throughputPerformance)
			if e != nil {
				return tccommon.RetryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update cbs failed, reason:%s\n ", logId, err.Error())
			return err
		}

	}

	if d.HasChange("disk_backup_quota") {
		diskBackupQuota := d.Get("disk_backup_quota").(int)
		err := cbsService.ModifyDiskBackupQuota(ctx, storageId, diskBackupQuota)
		if err != nil {
			return err
		}

	}
	if d.HasChange("tags") {

		oldValue, newValue := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))

		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := &TagService{client: tcClient}
		resourceName := tccommon.BuildTagResourceName("cvm", "volume", tcClient.Region, d.Id())
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}

	}
	//charge type
	//not support renew
	chargeType := d.Get("charge_type").(string)
	renewFlag := d.Get("prepaid_renew_flag").(string)
	period := d.Get("prepaid_period").(int)

	if d.HasChange("charge_type") {
		//only support postpaid to prepaid
		err := cbsService.ModifyDiskChargeType(ctx, storageId, chargeType, renewFlag, period)
		if err != nil {
			return err
		}

		//check charge Type
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			storage, e := cbsService.DescribeDiskById(ctx, storageId)
			if e != nil {
				return tccommon.RetryError(e)
			}
			if storage != nil && (*storage.DiskState == CBS_STORAGE_STATUS_EXPANDING || *storage.DiskChargeType != CBS_CHARGE_TYPE_PREPAID) {
				return resource.RetryableError(fmt.Errorf("cbs storage status is %s, charget type is %s", *storage.DiskState, *storage.DiskChargeType))
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update cbs failed, reason:%s\n ", logId, err.Error())
			return err
		}

	} else {

		//only renew and change flag
		if d.HasChange("prepaid_renew_flag") {
			//check

			if chargeType != CBS_CHARGE_TYPE_PREPAID {
				return fmt.Errorf("tencentcloud_cbs_storage update on prepaid_period or prepaid_renew_flag is only supported with charge type PREPAID")
			}

			//renew api
			err := cbsService.ModifyDisksRenewFlag(ctx, storageId, renewFlag)
			if err != nil {
				return err
			}
			//check renew flag
			err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				storage, e := cbsService.DescribeDiskById(ctx, storageId)
				if e != nil {
					return tccommon.RetryError(e)
				}
				if storage != nil && (*storage.DiskState == CBS_STORAGE_STATUS_EXPANDING || *storage.RenewFlag != renewFlag) {
					return resource.RetryableError(fmt.Errorf("cbs storage status is %s, renew flag is %s", *storage.DiskState, *storage.RenewFlag))
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s update cbs failed, reason:%s\n ", logId, err.Error())
				return err
			}

		}
	}

	d.Partial(false)

	return nil
}

func resourceTencentCloudCbsStorageDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_storage.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	storageId := d.Id()
	//check is force delete or not
	forceDelete := d.Get("force_delete").(bool)
	notExist := false
	cbsService := CbsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := cbsService.DeleteDiskById(ctx, storageId)
		if e != nil {
			return tccommon.RetryError(e, tccommon.InternalError)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete cbs failed, reason:%s\n ", logId, err.Error())
		return err
	}

	//check exist
	err = resource.Retry(2*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		storage, errRet := cbsService.DescribeDiskById(ctx, storageId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if storage == nil {
			notExist = true
			return nil
		}
		if *storage.DiskState == CBS_STORAGE_STATUS_TORECYCLE {
			//in recycling
			return nil
		}
		return resource.RetryableError(fmt.Errorf("cbs instance status is %s, retry...", *storage.DiskState))
	})

	if err != nil {
		return err
	}

	if notExist || !forceDelete {
		return nil
	}

	//exist in recycle

	//delete again
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		errRet := cbsService.DeleteDiskById(ctx, storageId)
		//when state is terminating, do not delete but check exist
		if errRet != nil {
			return tccommon.RetryError(errRet)
		}

		return nil
	})
	if err != nil {
		return err
	}

	//describe and check not exist
	err = resource.Retry(2*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		storage, errRet := cbsService.DescribeDiskById(ctx, storageId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if storage == nil {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("cbs instance status is %s, retry...", *storage.DiskState))
	})

	if err != nil {
		return err
	}
	return nil
}
