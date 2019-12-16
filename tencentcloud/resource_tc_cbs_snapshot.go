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

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
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

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	storageId := d.Get("storage_id").(string)
	snapshotName := d.Get("snapshot_name").(string)

	cbsService := CbsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	snapshotId := ""
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		var e error
		snapshotId, e = cbsService.CreateSnapshot(ctx, storageId, snapshotName)
		if e != nil {
			return retryError(e)
		}
		d.SetId(snapshotId)
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cbs snapshot failed, reason:%s\n ", logId, err.Error())
		return err
	}

	err = resource.Retry(20*readRetryTimeout, func() *resource.RetryError {
		snapshot, e := cbsService.DescribeSnapshotById(ctx, snapshotId)
		if e != nil {
			return retryError(e)
		}
		if snapshot == nil {
			return resource.RetryableError(fmt.Errorf("cbs snapshot is nil"))
		}
		if *snapshot.SnapshotState == CBS_SNAPSHOT_STATUS_CREATING {
			return resource.RetryableError(fmt.Errorf("cbs snapshot status is still %s", *snapshot.SnapshotState))
		}
		if *snapshot.SnapshotState == CBS_SNAPSHOT_STATUS_NORMAL {
			return nil
		}
		return resource.NonRetryableError(fmt.Errorf("cbs snapshot status is %s, we won't wait for it finish.", *snapshot.SnapshotState))
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cbs snapshot failed, reason:%s\n ", logId, err.Error())
		return err
	}

	return resourceTencentCloudCbsSnapshotRead(d, meta)
}

func resourceTencentCloudCbsSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cbs_snapshot.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	snapshotId := d.Id()
	cbsService := CbsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	var snapshot *cbs.Snapshot
	var e error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		snapshot, e = cbsService.DescribeSnapshotById(ctx, snapshotId)
		if e != nil {
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read cbs snapshot failed, reason:%s\n ", logId, err.Error())
		return err
	}
	if snapshot == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("disk_type", snapshot.DiskUsage)
	_ = d.Set("percent", snapshot.Percent)
	_ = d.Set("storage_size", snapshot.DiskSize)
	_ = d.Set("storage_id", snapshot.DiskId)
	_ = d.Set("snapshot_name", snapshot.SnapshotName)
	_ = d.Set("snapshot_status", snapshot.SnapshotState)

	return nil
}

func resourceTencentCloudCbsSnapshotUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cbs_snapshot.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	snapshotId := d.Id()

	if d.HasChange("snapshot_name") {
		snapshotName := d.Get("snapshot_name").(string)
		cbsService := CbsService{
			client: meta.(*TencentCloudClient).apiV3Conn,
		}
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := cbsService.ModifySnapshotName(ctx, snapshotId, snapshotName)
			if e != nil {
				return retryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update cbs snapshot failed, reason:%s\n ", logId, err.Error())
			return err
		}
	}

	return nil
}

func resourceTencentCloudCbsSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cbs_snapshot.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	snapshotId := d.Id()
	cbsService := CbsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := cbsService.DeleteSnapshot(ctx, snapshotId)
		if e != nil {
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete cbs snapshot failed, reason:%s\n ", logId, err.Error())
		return err
	}

	return nil
}
