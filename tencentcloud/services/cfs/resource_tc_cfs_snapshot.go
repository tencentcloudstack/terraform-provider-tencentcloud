package cfs

import (
	"context"
	"fmt"
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfs/v20190719"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCfsSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCfsSnapshotCreate,
		Read:   resourceTencentCloudCfsSnapshotRead,
		Update: resourceTencentCloudCfsSnapshotUpdate,
		Delete: resourceTencentCloudCfsSnapshotDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"file_system_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Id of file system.",
			},

			"snapshot_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Name of snapshot.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudCfsSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfs_snapshot.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = cfs.NewCreateCfsSnapshotRequest()
		response   = cfs.NewCreateCfsSnapshotResponse()
		snapshotId string
	)
	if v, ok := d.GetOk("file_system_id"); ok {
		request.FileSystemId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("snapshot_name"); ok {
		request.SnapshotName = helper.String(v.(string))
	}

	if v := helper.GetTags(d, "tags"); len(v) > 0 {
		for tagKey, tagValue := range v {
			tag := cfs.TagInfo{
				TagKey:   helper.String(tagKey),
				TagValue: helper.String(tagValue),
			}
			request.ResourceTags = append(request.ResourceTags, &tag)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfsClient().CreateCfsSnapshot(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cfs snapshot failed, reason:%+v", logId, err)
		return err
	}

	snapshotId = *response.Response.SnapshotId
	d.SetId(snapshotId)

	service := CfsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"available"}, 2*tccommon.ReadRetryTimeout, time.Second, service.CfsSnapshotStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := fmt.Sprintf("qcs::cfs:%s:uin/:snap/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudCfsSnapshotRead(d, meta)
}

func resourceTencentCloudCfsSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfs_snapshot.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CfsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	snapshotId := d.Id()

	snapshot, err := service.DescribeCfsSnapshotById(ctx, snapshotId)
	if err != nil {
		return err
	}

	if snapshot == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CfsSnapshot` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if snapshot.FileSystemId != nil {
		_ = d.Set("file_system_id", snapshot.FileSystemId)
	}

	if snapshot.SnapshotName != nil {
		_ = d.Set("snapshot_name", snapshot.SnapshotName)
	}

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, err := tagService.DescribeResourceTags(ctx, "cfs", "snap", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudCfsSnapshotUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfs_snapshot.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	needChange := false

	request := cfs.NewUpdateCfsSnapshotAttributeRequest()

	snapshotId := d.Id()

	request.SnapshotId = &snapshotId

	if d.HasChange("snapshot_name") {
		needChange = true
		if v, ok := d.GetOk("snapshot_name"); ok {
			request.SnapshotName = helper.String(v.(string))
		}
	}

	if needChange {
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfsClient().UpdateCfsSnapshotAttribute(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update cfs snapshot failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := tccommon.BuildTagResourceName("cfs", "snap", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudCfsSnapshotRead(d, meta)
}

func resourceTencentCloudCfsSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfs_snapshot.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CfsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	snapshotId := d.Id()

	if err := service.DeleteCfsSnapshotById(ctx, snapshotId); err != nil {
		return err
	}

	err := resource.Retry(2*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, errRet := service.DescribeCfsSnapshotById(ctx, snapshotId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if instance == nil {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("cfs snapshot status is %s, retry...", *instance.Status))
	})
	if err != nil {
		return err
	}

	return nil
}
