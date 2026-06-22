package cbs

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCbsCopySnapshotCrossRegionAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCbsCopySnapshotCrossRegionAttachmentCreate,
		Read:   resourceTencentCloudCbsCopySnapshotCrossRegionAttachmentRead,
		Delete: resourceTencentCloudCbsCopySnapshotCrossRegionAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"snapshot_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Source snapshot ID for cross-region copy.",
			},

			"destination_regions": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "Target region names for cross-region copy.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"snapshot_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Name of the copied snapshot.",
			},

			"delete_bind_images": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Whether to force-delete images associated with the snapshots when deleting.",
			},

			"snapshot_copy_result_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Cross-region copy results.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"snapshot_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "New snapshot ID in the destination region.",
						},
						"code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Error code, Success on success.",
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Error message, empty string on success.",
						},
						"destination_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Destination region for cross-region copy.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCbsCopySnapshotCrossRegionAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_copy_snapshot_cross_region.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = cbs.NewCopySnapshotCrossRegionsRequest()
	)

	if v, ok := d.GetOk("snapshot_id"); ok {
		request.SnapshotId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("destination_regions"); ok {
		regions := make([]string, 0, len(v.([]interface{})))
		for _, item := range v.([]interface{}) {
			regions = append(regions, item.(string))
		}
		request.DestinationRegions = helper.Strings(regions)
	}

	if v, ok := d.GetOk("snapshot_name"); ok {
		request.SnapshotName = helper.String(v.(string))
	}

	var response *cbs.CopySnapshotCrossRegionsResponse
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCbsClient().CopySnapshotCrossRegionsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			log.Printf("[CRITAL]%s logId=%s d.Id()=%s CopySnapshotCrossRegions return nil Response", logId, logId, d.Id())
			return resource.NonRetryableError(fmt.Errorf("CopySnapshotCrossRegions return nil Response"))
		}

		if result.Response.SnapshotCopyResultSet == nil || len(result.Response.SnapshotCopyResultSet) == 0 {
			log.Printf("[CRITAL]%s logId=%s d.Id()=%s CopySnapshotCrossRegions return empty SnapshotCopyResultSet", logId, logId, d.Id())
			return resource.NonRetryableError(fmt.Errorf("CopySnapshotCrossRegions return empty SnapshotCopyResultSet"))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create cbs_copy_snapshot_cross_region failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	snapshotCopyResultSet := response.Response.SnapshotCopyResultSet

	// Check each SnapshotCopyResult entry's Code field
	for _, result := range snapshotCopyResultSet {
		if result.Code != nil && *result.Code != "Success" {
			return fmt.Errorf("cbs_copy_snapshot_cross_region copy failed, Code: %s, Message: %s", *result.Code, *result.Message)
		}
	}

	// Set composite ID using snapshot_id + FILED_SP + copied_snapshot_id (using first copied snapshot ID)
	snapshotId := d.Get("snapshot_id").(string)
	copiedSnapshotId := *snapshotCopyResultSet[0].SnapshotId
	d.SetId(strings.Join([]string{snapshotId, copiedSnapshotId}, tccommon.FILED_SP))

	// Save snapshot_copy_result_set to state
	snapshotCopyResultList := make([]map[string]interface{}, 0, len(snapshotCopyResultSet))
	for _, result := range snapshotCopyResultSet {
		resultMap := map[string]interface{}{}
		if result.SnapshotId != nil {
			resultMap["snapshot_id"] = *result.SnapshotId
		}
		if result.Code != nil {
			resultMap["code"] = *result.Code
		}
		if result.Message != nil {
			resultMap["message"] = *result.Message
		}
		if result.DestinationRegion != nil {
			resultMap["destination_region"] = *result.DestinationRegion
		}
		snapshotCopyResultList = append(snapshotCopyResultList, resultMap)
	}
	_ = d.Set("snapshot_copy_result_set", snapshotCopyResultList)

	// Async polling: for each copied snapshot_id, poll DescribeSnapshots until SnapshotState is NORMAL
	cbsService := CbsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	for _, result := range snapshotCopyResultSet {
		if result.SnapshotId != nil && *result.Code == "Success" {
			copiedId := *result.SnapshotId
			err := resource.Retry(20*tccommon.ReadRetryTimeout, func() *resource.RetryError {
				snapshot, e := cbsService.DescribeSnapshotById(ctx, copiedId)
				if e != nil {
					return tccommon.RetryError(e, tccommon.InternalError)
				}
				if snapshot == nil || snapshot.SnapshotState == nil {
					return resource.RetryableError(fmt.Errorf("cbs_copy_snapshot_cross_region copied snapshot %s not ready yet", copiedId))
				}
				if *snapshot.SnapshotState != CBS_SNAPSHOT_STATUS_NORMAL {
					return resource.RetryableError(fmt.Errorf("cbs_copy_snapshot_cross_region copied snapshot %s state is %s, waiting for NORMAL", copiedId, *snapshot.SnapshotState))
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s cbs_copy_snapshot_cross_region async polling failed for snapshot %s, reason:%+v", logId, copiedId, err)
				return err
			}
		}
	}

	return resourceTencentCloudCbsCopySnapshotCrossRegionAttachmentRead(d, meta)
}

func resourceTencentCloudCbsCopySnapshotCrossRegionAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_copy_snapshot_cross_region.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		cbsService = CbsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("cbs_copy_snapshot_cross_region id is illegal: %s", d.Id())
	}

	snapshotId := idSplit[0]
	copiedSnapshotId := idSplit[1]

	var snapshot *cbs.Snapshot
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := cbsService.DescribeSnapshotById(ctx, copiedSnapshotId)
		if e != nil {
			return tccommon.RetryError(e, tccommon.InternalError)
		}
		snapshot = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read cbs_copy_snapshot_cross_region failed, reason:%+v", logId, err)
		return err
	}

	if snapshot == nil {
		log.Printf("[CRUD] cbs_copy_snapshot_cross_region id=%s", d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("snapshot_id", snapshotId)

	if snapshot.SnapshotName != nil {
		_ = d.Set("snapshot_name", snapshot.SnapshotName)
	}

	// Populate snapshot_copy_result_set if not already set (e.g., after import)
	// During normal Create-then-Read flow, snapshot_copy_result_set is preserved from state.
	// After import, we need to reconstruct it from the composite ID and API response.
	existingResults := d.Get("snapshot_copy_result_set").([]interface{})
	if len(existingResults) == 0 {
		resultMap := map[string]interface{}{
			"snapshot_id": copiedSnapshotId,
			"code":        "Success",
			"message":     "",
		}
		if snapshot.Placement != nil && snapshot.Placement.Zone != nil {
			resultMap["destination_region"] = *snapshot.Placement.Zone
		}
		_ = d.Set("snapshot_copy_result_set", []map[string]interface{}{resultMap})
	}

	return nil
}

func resourceTencentCloudCbsCopySnapshotCrossRegionAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_copy_snapshot_cross_region.delete")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = cbs.NewDeleteSnapshotsRequest()
	)

	// Read all copied snapshot IDs from state's snapshot_copy_result_set field
	snapshotCopyResultSet := d.Get("snapshot_copy_result_set").([]interface{})
	snapshotIds := make([]*string, 0, len(snapshotCopyResultSet))
	for _, item := range snapshotCopyResultSet {
		resultMap := item.(map[string]interface{})
		if v, ok := resultMap["snapshot_id"].(string); ok && v != "" {
			snapshotIds = append(snapshotIds, helper.String(v))
		}
	}

	// Fallback: if snapshot_copy_result_set is empty (e.g., after import),
	// parse the composite ID to get the copied_snapshot_id
	if len(snapshotIds) == 0 {
		idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
		if len(idSplit) == 2 && idSplit[1] != "" {
			snapshotIds = append(snapshotIds, helper.String(idSplit[1]))
		}
	}

	if len(snapshotIds) == 0 {
		log.Printf("[WARN]%s cbs_copy_snapshot_cross_region no copied snapshot IDs found, skip delete", logId)
		return nil
	}

	request.SnapshotIds = snapshotIds

	if v, ok := d.GetOk("delete_bind_images"); ok {
		request.DeleteBindImages = helper.Bool(v.(bool))
	} else {
		request.DeleteBindImages = helper.Bool(false)
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCbsClient().DeleteSnapshotsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			if result == nil || result.Response == nil {
				log.Printf("[CRITAL]%s logId=%s d.Id()=%s DeleteSnapshots return nil Response", logId, logId, d.Id())
				return resource.NonRetryableError(fmt.Errorf("DeleteSnapshots return nil Response"))
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete cbs_copy_snapshot_cross_region failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
