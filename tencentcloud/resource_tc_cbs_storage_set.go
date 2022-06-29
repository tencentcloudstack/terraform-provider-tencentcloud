/*
Provides a resource to create CBS set.

Example Usage

```hcl
resource "tencentcloud_cbs_storage_set" "storage" {
        disk_count 		  = 10
        storage_name      = "mystorage"
        storage_type      = "CLOUD_SSD"
        storage_size      = 100
        availability_zone = "ap-guangzhou-3"
        project_id        = 0
        encrypt           = false
}
```

*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCbsStorageSet() *schema.Resource {
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
				Description: "Type of CBS medium. Valid values: CLOUD_PREMIUM, CLOUD_SSD, CLOUD_TSSD and CLOUD_HSSD.",
			},
			"storage_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerInRange(10, 16000),
				Description:  "Volume of CBS, and unit is GB. If storage type is `CLOUD_SSD`, the size range is [100, 16000], and the others are [10-16000].",
			},
			"disk_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The number of disks to be purchased. Default 1.",
			},
			"charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      CBS_CHARGE_TYPE_POSTPAID,
				ValidateFunc: validateAllowedStringValue(CBS_CHARGE_TYPE),
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
				ValidateFunc: validateStringLengthInRange(2, 60),
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
	defer logElapsed("resource.tencentcloud_cbs_storage_set.create")()

	logId := getLogId(contextNil)

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
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		response, e := meta.(*TencentCloudClient).apiV3Conn.UseCbsClient().CreateDisks(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())

			ee, ok := e.(*sdkErrors.TencentCloudSDKError)
			if ok && IsContains(CVM_RETRYABLE_ERROR, ee.Code) {
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

	d.Set("disk_ids", storageIds)
	d.SetId(helper.StrListToStr(storageIds))

	return nil
}

func resourceTencentCloudCbsStorageSetRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cbs_storage_set.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	storageId := d.Id()
	cbsService := CbsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
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
	defer logElapsed("resource.tencentcloud_cbs_storage_set.update")()

	return fmt.Errorf("`tencentcloud_cbs_storage_set` do not support change now.")
}

func resourceTencentCloudCbsStorageSetDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cbs_storage_set.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	storageId := d.Id()

	cbsService := CbsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := cbsService.DeleteDiskSetByIds(ctx, storageId)
		if e != nil {
			log.Printf("[CRITAL][first delete]%s api[%s] fail, reason[%s]\n",
				logId, "delete", e.Error())
			ee, ok := e.(*sdkErrors.TencentCloudSDKError)
			if ok && IsContains(CVM_RETRYABLE_ERROR, ee.Code) {
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
