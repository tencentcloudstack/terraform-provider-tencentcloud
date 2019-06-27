/*
Provide a resource to create a CBS.

Example Usage

```hcl
resource "tencentcloud_cbs_storage" "storage" {
        storage_name      = "mystorage"
        storage_type      = "CLOUD_SSD"
        storage_size      = "50"
        availability_zone = "ap-guangzhou-3"
		project_id        = 0
		encrypt = false
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
	"github.com/hashicorp/terraform/helper/resource"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
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
				Description:  "Type of CBS medium, and available values include CLOUD_BASIC, CLOUD_PREMIUM and CLOUD_SSD.",
			},
			"storage_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerInRange(10, 16000),
				Description:  "Volume of CBS.",
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntegerInRange(1, 36),
				Description:  "The purchased usage period of CBS, and value range [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36].",
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

			// computed
			"storage_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of CBS, and available values include UNATTACHED, ATTACHING, ATTACHED, DETACHING, EXPANDING, ROLLBACKING, TORECYCLE and DUMPING.",
			},
			"attached": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Indicates whether the CBS is mounted the CVM.",
			},
		},
	}
}

func resourceTencentCloudCbsStorageCreate(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	request := cbs.NewCreateDisksRequest()

	request.DiskName = stringToPointer(d.Get("storage_name").(string))
	request.DiskType = stringToPointer(d.Get("storage_type").(string))
	request.DiskSize = intToPointer(d.Get("storage_size").(int))
	request.Placement = &cbs.Placement{
		Zone: stringToPointer(d.Get("availability_zone").(string)),
	}
	if v, ok := d.GetOk("project_id"); ok {
		request.Placement.ProjectId = intToPointer(v.(int))
	}
	if v, ok := d.GetOk("snapshot_id"); ok {
		request.SnapshotId = stringToPointer(v.(string))
	}
	if _, ok := d.GetOk("encrypt"); ok {
		request.Encrypt = stringToPointer("ENCRYPT")
	}
	if v, ok := d.GetOk("tags"); ok {
		tags := v.(map[string]interface{})
		request.Tags = make([]*cbs.Tag, 0, len(tags))
		for key, value := range tags {
			tag := cbs.Tag{
				Key:   &key,
				Value: stringToPointer(value.(string)),
			}
			request.Tags = append(request.Tags, &tag)
		}
	}
	request.DiskChargeType = stringToPointer("POSTPAID_BY_HOUR")

	response, err := meta.(*TencentCloudClient).apiV3Conn.UseCbsClient().CreateDisks(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	} else {
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	}
	if len(response.Response.DiskIdSet) < 1 {
		return fmt.Errorf("storage id is nil")
	}
	d.SetId(*response.Response.DiskIdSet[0])

	// must wait for finishing creating disk
	time.Sleep(3 * time.Second)

	return resourceTencentCloudCbsStorageRead(d, meta)
}

func resourceTencentCloudCbsStorageRead(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	storageId := d.Id()
	cbsService := CbsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	storage, err := cbsService.DescribeDiskById(ctx, storageId)
	if err != nil {
		return err
	}

	d.Set("storage_type", storage.DiskType)
	d.Set("storage_size", storage.DiskSize)
	d.Set("availability_zone", storage.Placement.Zone)
	d.Set("storage_name", storage.DiskName)
	d.Set("project_id", storage.Placement.ProjectId)
	d.Set("encrypt", storage.Encrypt)
	d.Set("tags", flattenCbsTagsMapping(storage.Tags))
	d.Set("storage_status", storage.DiskState)
	d.Set("attached", storage.Attached)

	return nil
}

func resourceTencentCloudCbsStorageUpdate(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	d.Partial(true)
	cbsService := CbsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
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
		err := cbsService.ModifyDiskAttributes(ctx, storageId, storageName, projectId)
		if err != nil {
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
		old, new := d.GetChange("storage_size")
		oldValue := old.(int)
		newValue := new.(int)
		if oldValue > newValue {
			return fmt.Errorf("storage size must be greater than current storage size")
		}

		err := cbsService.ResizeDisk(ctx, storageId, newValue)
		if err != nil {
			return err
		}
		err = resource.Retry(3*time.Minute, func() *resource.RetryError {
			storage, e := cbsService.DescribeDiskById(ctx, storageId)
			if e != nil {
				return resource.NonRetryableError(e)
			}
			if *storage.DiskState == CBS_STORAGE_STATUS_EXPANDING {
				return resource.RetryableError(fmt.Errorf("cbs storage status is %s", *storage.DiskState))
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s cbs storage create failed, reason:%s\n ", logId, err.Error())
			return err
		}
		d.SetPartial("storage_size")
	}

	if d.HasChange("snapshot_id") {
		snapshotId := d.Get("snapshot_id").(string)
		err := cbsService.ApplySnapshot(ctx, storageId, snapshotId)
		if err != nil {
			return err
		}
		err = resource.Retry(3*time.Minute, func() *resource.RetryError {
			storage, e := cbsService.DescribeDiskById(ctx, storageId)
			if e != nil {
				return resource.NonRetryableError(e)
			}
			if *storage.DiskState == CBS_STORAGE_STATUS_ROLLBACKING {
				return resource.RetryableError(fmt.Errorf("cbs storage status is %s", *storage.DiskState))
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s cbs storage create failed, reason:%s\n ", logId, err.Error())
			return err
		}
		d.SetPartial("snapshot_id")
	}

	d.Partial(false)

	return nil
}

func resourceTencentCloudCbsStorageDelete(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	storageId := d.Id()
	cbsService := CbsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := cbsService.DeleteDiskById(ctx, storageId)
	if err != nil {
		return err
	}
	return nil
}
