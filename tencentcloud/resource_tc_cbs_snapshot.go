package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTencentCloudCbsSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCbsSnapshotCreate,
		Read:   resourceTencentCloudCbsSnapshotRead,
		Update: resourceTencentCloudCbsSnapshotUpdate,
		Delete: resourceTencentCloudCbsSnapshotDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"snapshot_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(2, 60),
			},
			"storage_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"storage_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"snapshot_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pecent": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTencentCloudCbsSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	storageId := d.Get("storage_id").(string)
	snapshotName := d.Get("snapshot_name").(string)
	cbsService := CbsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	snapshotId, err := cbsService.CreateSnapshot(ctx, storageId, snapshotName)
	if err != nil {
		return nil
	}
	d.SetId(snapshotId)

	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		snapshot, e := cbsService.DescribeSnapshotById(ctx, snapshotId)
		if e != nil {
			return resource.NonRetryableError(e)
		}
		if *snapshot.SnapshotState == CBS_SNAPSHOT_STATUS_CREATING {
			return resource.RetryableError(fmt.Errorf("cbs snapshot status is %s", *snapshot.SnapshotState))
		}
		if *snapshot.SnapshotState == CBS_SNAPSHOT_STATUS_NORMAL {
			return nil
		}
		e = fmt.Errorf("cbs snapshot status is %s, we won't wait for it finish.", *snapshot.SnapshotState)
		return resource.NonRetryableError(e)
	})
	if err != nil {
		log.Printf("[CRITAL]%s cbs snapshot attachment failed, reason:%s\n ", logId, err.Error())
		return err
	}

	return resourceTencentCloudCbsSnapshotRead(d, meta)
}

func resourceTencentCloudCbsSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	snapshotId := d.Id()
	cbsService := CbsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	snapshot, err := cbsService.DescribeSnapshotById(ctx, snapshotId)
	if err != nil {
		return err
	}

	d.Set("disk_type", snapshot.DiskUsage)
	d.Set("percent", snapshot.Percent)
	d.Set("storage_size", snapshot.DiskSize)
	d.Set("storage_id", snapshot.DiskId)
	d.Set("snapshot_name", snapshot.SnapshotName)
	d.Set("snapshot_status", snapshot.SnapshotState)
	return nil
}

func resourceTencentCloudCbsSnapshotUpdate(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)
	snapshotId := d.Id()

	if d.HasChange("snapshot_name") {
		snapshotName := d.Get("snapshot_name").(string)
		cbsService := CbsService{
			client: meta.(*TencentCloudClient).apiV3Conn,
		}
		err := cbsService.ModifySnapshotName(ctx, snapshotId, snapshotName)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceTencentCloudCbsSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	snapshotId := d.Id()
	cbsService := CbsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := cbsService.DeleteSnapshot(ctx, snapshotId)
	if err != nil {
		return err
	}
	return nil
}
