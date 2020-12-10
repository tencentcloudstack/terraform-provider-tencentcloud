/*
Provides a resource to create a CBS.

Example Usage

```hcl
resource "tencentcloud_cbs_storage" "storage" {
  storage_name      = "mystorage"
  storage_type      = "CLOUD_SSD"
  storage_size      = 100
  availability_zone = "ap-guangzhou-3"
  project_id        = 0
  encrypt           = false

  tags = {
    test = "tf"
  }
}
```

Import

CBS storage can be imported using the id, e.g.

```
$ terraform import tencentcloud_cbs_storage.storage disk-41s6jwy4
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCbsStorage() *schema.Resource {
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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(CBS_STORAGE_TYPE),
				Description:  "Type of CBS medium. Valid values: CLOUD_BASIC, CLOUD_PREMIUM and CLOUD_SSD.",
			},
			"storage_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerInRange(10, 16000),
				Description:  "Volume of CBS, and unit is GB. If storage type is `CLOUD_SSD`, the size range is [100, 16000], and the others are [10-16000].",
			},
			"period": {
				Deprecated:   "It has been deprecated from version 1.33.0. Set `prepaid_period` instead.",
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntegerInRange(1, 36),
				Description:  "The purchased usage period of CBS. Valid values: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36].",
			},
			"charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      CBS_CHARGE_TYPE_POSTPAID,
				ValidateFunc: validateAllowedStringValue(CBS_CHARGE_TYPE),
				Description:  "The charge type of CBS instance. Valid values are `PREPAID` and `POSTPAID_BY_HOUR`. The default is `POSTPAID_BY_HOUR`.",
			},
			"prepaid_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateAllowedIntValue(CBS_PREPAID_PERIOD),
				Description:  "The tenancy (time unit is month) of the prepaid instance, NOTE: it only works when charge_type is set to `PREPAID`. Valid values are 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36.",
			},
			"prepaid_renew_flag": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateAllowedStringValue(CBS_PREPAID_RENEW_FLAG),
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
	defer logElapsed("resource.tencentcloud_cbs_storage.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cbsService := CbsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
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

	storageId := ""
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		response, e := meta.(*TencentCloudClient).apiV3Conn.UseCbsClient().CreateDisks(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return retryError(e, InternalError)
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

	// must wait for finishing creating disk
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		storage, e := cbsService.DescribeDiskById(ctx, storageId)
		if e != nil {
			return retryError(e, InternalError)
		}
		if storage == nil {
			return resource.RetryableError(fmt.Errorf("storage is still creating..."))
		}
		return nil
	})
	if err != nil {
		return err
	}
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		resourceName := BuildTagResourceName("cvm", "volume", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudCbsStorageRead(d, meta)
}

func resourceTencentCloudCbsStorageRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cbs_storage.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	storageId := d.Id()
	cbsService := CbsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	var storage *cbs.Disk
	var e error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		storage, e = cbsService.DescribeDiskById(ctx, storageId)
		if e != nil {
			return retryError(e)
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

	if *storage.DiskChargeType == CBS_CHARGE_TYPE_PREPAID {
		_ = d.Set("prepaid_renew_flag", storage.RenewFlag)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "cvm", "volume", tcClient.Region, d.Id())
	if err != nil {
		return err
	}

	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudCbsStorageUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cbs_storage.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	//only support update prepaid_period when upgrade chargeType
	if d.HasChange("prepaid_period") && (!d.HasChange("charge_type") && d.Get("charge_type").(string) == CBS_CHARGE_TYPE_PREPAID) {
		return fmt.Errorf("tencentcloud_cbs_storage renew is not support yet")
	}
	if d.HasChange("charge_type") && d.Get("charge_type").(string) != CBS_CHARGE_TYPE_PREPAID {
		return fmt.Errorf("tencentcloud_cbs_storage do not support downgrade instance")
	}

	cbsService := CbsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
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
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := cbsService.ModifyDiskAttributes(ctx, storageId, storageName, projectId)
			if e != nil {
				return retryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update cbs failed, reason:%s\n ", logId, err.Error())
			return err
		}
		if d.HasChange("storage_name") {
			d.SetPartial("storage_name")
		}
		if d.HasChange("project_id") {
			d.SetPartial("project_id")
		}
	}

	if d.HasChange("storage_size") {
		oldInterface, newInterface := d.GetChange("storage_size")
		oldValue := oldInterface.(int)
		newValue := newInterface.(int)
		if oldValue > newValue {
			return fmt.Errorf("storage size must be greater than current storage size")
		}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := cbsService.ResizeDisk(ctx, storageId, newValue)
			if e != nil {
				return retryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update cbs failed, reason:%s\n ", logId, err.Error())
			return err
		}

		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			storage, e := cbsService.DescribeDiskById(ctx, storageId)
			if e != nil {
				return retryError(e)
			}
			if storage != nil && *storage.DiskState == CBS_STORAGE_STATUS_EXPANDING {
				return resource.RetryableError(fmt.Errorf("cbs storage status is %s", *storage.DiskState))
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update cbs failed, reason:%s\n ", logId, err.Error())
			return err
		}

		d.SetPartial("storage_size")
	}

	if d.HasChange("snapshot_id") {
		snapshotId := d.Get("snapshot_id").(string)
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := cbsService.ApplySnapshot(ctx, storageId, snapshotId)
			if e != nil {
				return retryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update cbs failed, reason:%s\n ", logId, err.Error())
			return err
		}

		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			storage, e := cbsService.DescribeDiskById(ctx, storageId)
			if e != nil {
				return retryError(e)
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

		d.SetPartial("snapshot_id")
	}

	if d.HasChange("tags") {

		oldValue, newValue := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))

		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		resourceName := BuildTagResourceName("cvm", "volume", tcClient.Region, d.Id())
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
		d.SetPartial("tags")
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
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			storage, e := cbsService.DescribeDiskById(ctx, storageId)
			if e != nil {
				return retryError(e)
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

		d.SetPartial("charge_type")
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
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				storage, e := cbsService.DescribeDiskById(ctx, storageId)
				if e != nil {
					return retryError(e)
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
			d.SetPartial("prepaid_renew_flag")

		}
	}

	d.Partial(false)

	return nil
}

func resourceTencentCloudCbsStorageDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cbs_storage.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	storageId := d.Id()
	//check is force delete or not
	forceDelete := d.Get("force_delete").(bool)
	notExist := false
	cbsService := CbsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := cbsService.DeleteDiskById(ctx, storageId)
		if e != nil {
			return retryError(e, InternalError)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete cbs failed, reason:%s\n ", logId, err.Error())
		return err
	}

	//check exist
	err = resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
		storage, errRet := cbsService.DescribeDiskById(ctx, storageId)
		if errRet != nil {
			return retryError(errRet, InternalError)
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
	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		errRet := cbsService.DeleteDiskById(ctx, storageId)
		//when state is terminating, do not delete but check exist
		if errRet != nil {
			return retryError(errRet)
		}

		return nil
	})
	if err != nil {
		return err
	}

	//describe and check not exist
	err = resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
		storage, errRet := cbsService.DescribeDiskById(ctx, storageId)
		if errRet != nil {
			return retryError(errRet, InternalError)
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
