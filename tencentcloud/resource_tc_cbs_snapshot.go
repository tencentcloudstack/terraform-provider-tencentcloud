/*
Provides a resource to create a CBS snapshot.

Example Usage

```hcl
resource "tencentcloud_cbs_snapshot" "snapshot" {
  snapshot_name = "unnamed"
  storage_id   = "disk-kdt0sq6m"
}
```

Import

CBS snapshot can be imported using the id, e.g.

```
$ terraform import tencentcloud_cbs_snapshot.snapshot snap-3sa3f39b
```
*/
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
				Description:  "Name of the snapshot.",
			},
			"storage_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the the CBS which this snapshot created from.",
			},
			"storage_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Volume of storage which this snapshot created from.",
			},
			"snapshot_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the snapshot.",
			},
			"disk_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Types of CBS which this snapshot created from.",
			},
			"percent": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Snapshot creation progress percentage. If the snapshot has created successfully, the constant value is 100.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of snapshot.",
			},
		},
	}
}

func resourceTencentCloudCbsSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cbs_snapshot.create")()

	logId := getLogId(nil)
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
	defer logElapsed("resource.tencentcloud_cbs_snapshot.read")()

	logId := getLogId(nil)
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
	defer logElapsed("resource.tencentcloud_cbs_snapshot.update")()

	logId := getLogId(nil)
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
	defer logElapsed("resource.tencentcloud_cbs_snapshot.delete")()

	logId := getLogId(nil)
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
