package cbs

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCbsStorageSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCbsStorageSetCreate,
		Read:   resourceTencentCloudCbsStorageSetRead,
		Update: resourceTencentCloudCbsStorageSetUpdate,
		Delete: resourceTencentCloudCbsStorageSetDelete,

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
			"disk_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "The number of disks to be purchased. Default 1.",
			},
			"charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      CBS_CHARGE_TYPE_POSTPAID,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CBS_CHARGE_TYPE),
				Description:  "The charge type of CBS instance. Only support `POSTPAID_BY_HOUR`.",
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
			"throughput_performance": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Add extra performance to the data disk. Only works when disk type is `CLOUD_TSSD` or `CLOUD_HSSD`.",
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
			"disk_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "disk id list.",
			},
		},
	}
}

func resourceTencentCloudCbsStorageSetCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_storage_set.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var diskCount int

	request := cbs.NewCreateDisksRequest()
	request.DiskName = helper.String(d.Get("storage_name").(string))
	request.DiskType = helper.String(d.Get("storage_type").(string))
	request.DiskSize = helper.IntUint64(d.Get("storage_size").(int))
	if v, ok := d.GetOk("disk_count"); ok {
		diskCount = v.(int)
		request.DiskCount = helper.Uint64(uint64(diskCount))
	}
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

	storageIds := make([]*string, 0)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCbsClient().CreateDisks(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())

			ee, ok := e.(*sdkErrors.TencentCloudSDKError)
			if ok && tccommon.IsContains(CVM_RETRYABLE_ERROR, ee.Code) {
				time.Sleep(1 * time.Second) // 需要重试的话，等待1s进行重试
				return resource.RetryableError(fmt.Errorf("cbs create error: %s, retrying", e.Error()))
			}
			return resource.NonRetryableError(e)
		}

		if len(response.Response.DiskIdSet) < diskCount {
			err := fmt.Errorf("number of instances is less than %s", strconv.Itoa(diskCount))
			return resource.NonRetryableError(err)
		}

		storageIds = response.Response.DiskIdSet
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cbs failed, reason:%s\n ", logId, err.Error())
		return err
	}

	_ = d.Set("disk_ids", storageIds)
	d.SetId(helper.StrListToStr(storageIds))

	return nil
}

func resourceTencentCloudCbsStorageSetRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_storage_set.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	storageId := d.Id()
	cbsService := CbsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	var storageSet []*cbs.Disk
	var errRet error

	storageSet, errRet = cbsService.DescribeDiskSetByIds(ctx, storageId)
	if errRet != nil {
		return errRet
	}
	if storageSet == nil {
		d.SetId("")
		return nil
	}

	storage := storageSet[0]

	_ = d.Set("disk_count", len(storageSet))
	_ = d.Set("storage_type", storage.DiskType)
	_ = d.Set("storage_size", storage.DiskSize)
	_ = d.Set("availability_zone", storage.Placement.Zone)
	_ = d.Set("storage_name", d.Get("storage_name"))
	_ = d.Set("project_id", storage.Placement.ProjectId)
	_ = d.Set("encrypt", storage.Encrypt)
	_ = d.Set("storage_status", storage.DiskState)
	_ = d.Set("attached", storage.Attached)
	_ = d.Set("charge_type", storage.DiskChargeType)
	_ = d.Set("throughput_performance", storage.ThroughputPerformance)

	return nil
}

func resourceTencentCloudCbsStorageSetUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_storage_set.update")()

	return fmt.Errorf("`tencentcloud_cbs_storage_set` do not support change now.")
}

func resourceTencentCloudCbsStorageSetDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_storage_set.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	storageId := d.Id()

	cbsService := CbsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := cbsService.DeleteDiskSetByIds(ctx, storageId)
		if e != nil {
			log.Printf("[CRITAL][first delete]%s api[%s] fail, reason[%s]\n",
				logId, "delete", e.Error())
			ee, ok := e.(*sdkErrors.TencentCloudSDKError)
			if ok && tccommon.IsContains(CVM_RETRYABLE_ERROR, ee.Code) {
				time.Sleep(1 * time.Second) // 需要重试的话，等待1s进行重试
				return resource.RetryableError(fmt.Errorf("[first delete]cvm delete error: %s, retrying", ee.Error()))
			}
			return resource.NonRetryableError(e)
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete cbs failed, reason:%s\n ", logId, err.Error())
		return err
	}

	return nil
}
