package cbs

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
)

func ResourceTencentCloudCbsSnapshot() *schema.Resource {
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
				ValidateFunc: tccommon.ValidateStringLengthInRange(2, 60),
				Description:  "Name of the snapshot.",
			},
			"storage_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the the CBS which this snapshot created from.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Deprecated:  "cbs snapshot do not support tag now.",
				Description: "The available tags within this CBS Snapshot.",
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
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_snapshot.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	storageId := d.Get("storage_id").(string)
	snapshotName := d.Get("snapshot_name").(string)

	var tags map[string]string

	if temp := helper.GetTags(d, "tags"); len(temp) > 0 {
		tags = temp
	}
	cbsService := CbsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	snapshotId := ""
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		var e error
		snapshotId, e = cbsService.CreateSnapshot(ctx, storageId, snapshotName, tags)
		if e != nil {
			return tccommon.RetryError(e)
		}
		d.SetId(snapshotId)
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cbs snapshot failed, reason:%s\n ", logId, err.Error())
		return err
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := &TagService{client: tcClient}
		resourceName := tccommon.BuildTagResourceName("cvm", "volume", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	err = resource.Retry(20*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		snapshot, e := cbsService.DescribeSnapshotById(ctx, snapshotId)
		if e != nil {
			return tccommon.RetryError(e)
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
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_snapshot.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	snapshotId := d.Id()
	cbsService := CbsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	var snapshot *cbs.Snapshot
	var e error
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		snapshot, e = cbsService.DescribeSnapshotById(ctx, snapshotId)
		if e != nil {
			return tccommon.RetryError(e)
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

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "cvm", "volume", tcClient.Region, d.Id())
	if err != nil {
		return err
	}

	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudCbsSnapshotUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_snapshot.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	snapshotId := d.Id()

	if d.HasChange("snapshot_name") {
		snapshotName := d.Get("snapshot_name").(string)
		cbsService := CbsService{
			client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
		}
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			e := cbsService.ModifySnapshotName(ctx, snapshotId, snapshotName)
			if e != nil {
				return tccommon.RetryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update cbs snapshot failed, reason:%s\n ", logId, err.Error())
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

	return nil
}

func resourceTencentCloudCbsSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_snapshot.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	snapshotId := d.Id()
	cbsService := CbsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := cbsService.DeleteSnapshot(ctx, snapshotId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete cbs snapshot failed, reason:%s\n ", logId, err.Error())
		return err
	}

	return nil
}
