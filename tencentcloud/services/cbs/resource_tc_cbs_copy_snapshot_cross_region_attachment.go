package cbs

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCbsCopySnapshotCrossRegionAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCbsCopySnapshotCrossRegionAttachmentCreate,
		Read:   resourceTencentCloudCbsCopySnapshotCrossRegionAttachmentRead,
		Update: resourceTencentCloudCbsCopySnapshotCrossRegionAttachmentUpdate,
		Delete: resourceTencentCloudCbsCopySnapshotCrossRegionAttachmentDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"snapshot_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Source snapshot ID for cross-region copy.",
			},

			"destination_region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Target region name for cross-region copy.",
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
				Default:     false,
				Description: "Whether to force-delete images associated with the snapshots when deleting.",
			},
		},
	}
}

// cbsClientWithRegion builds a CBS client bound to the given region using the
// provider's credential and client profile. The snapshots copied by
// CopySnapshotCrossRegions reside in the destination region, so describing and
// deleting them must use a client pointing at that region instead of the
// provider's default region.
func cbsClientWithRegion(meta interface{}, region string) *cbs.Client {
	conn := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	cpf := conn.NewClientProfile(300)
	client, _ := cbs.NewClient(conn.Credential, region, cpf)
	// Attach the LogRoundTripper so that requests issued by this region-specific client
	// (e.g. DescribeSnapshots / DeleteSnapshots) are printed in the SDK debug log,
	// consistent with clients created via UseCbsClient().
	client.WithHttpTransport(&connectivity.LogRoundTripper{})
	return client
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

	if v, ok := d.GetOk("destination_region"); ok {
		request.DestinationRegions = helper.Strings([]string{v.(string)})
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

	// Set composite ID using copied_snapshot_id (returned by CopySnapshotCrossRegions)
	// + FILED_SP + destination_region (using the first copied snapshot result).
	copiedSnapshotId := *snapshotCopyResultSet[0].SnapshotId
	destinationRegion := d.Get("destination_region").(string)
	if snapshotCopyResultSet[0].DestinationRegion != nil {
		destinationRegion = *snapshotCopyResultSet[0].DestinationRegion
	}
	d.SetId(strings.Join([]string{copiedSnapshotId, destinationRegion}, tccommon.FILED_SP))

	// Async polling: CopySnapshotCrossRegions is an asynchronous API. For each copied
	// snapshot, poll DescribeSnapshots in its DestinationRegion (using a region-specific
	// client) until SnapshotState becomes NORMAL, which indicates the copy succeeded.
	for _, result := range snapshotCopyResultSet {
		if result.SnapshotId == nil || result.DestinationRegion == nil || result.Code == nil || *result.Code != "Success" {
			continue
		}

		copiedId := *result.SnapshotId
		region := *result.DestinationRegion
		regionCbsClient := cbsClientWithRegion(meta, region)

		err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			describeRequest := cbs.NewDescribeSnapshotsRequest()
			describeRequest.SnapshotIds = []*string{helper.String(copiedId)}
			describeResponse, e := regionCbsClient.DescribeSnapshotsWithContext(ctx, describeRequest)
			if e != nil {
				return tccommon.RetryError(e, tccommon.InternalError)
			}

			if describeResponse == nil || describeResponse.Response == nil || len(describeResponse.Response.SnapshotSet) == 0 {
				return resource.RetryableError(fmt.Errorf("cbs_copy_snapshot_cross_region copied snapshot %s in region %s not ready yet", copiedId, region))
			}

			snapshot := describeResponse.Response.SnapshotSet[0]
			if snapshot.SnapshotState == nil {
				return resource.RetryableError(fmt.Errorf("cbs_copy_snapshot_cross_region copied snapshot %s in region %s state unknown", copiedId, region))
			}

			if *snapshot.SnapshotState != CBS_SNAPSHOT_STATUS_NORMAL {
				return resource.RetryableError(fmt.Errorf("cbs_copy_snapshot_cross_region copied snapshot %s in region %s state is %s, waiting for NORMAL", copiedId, region, *snapshot.SnapshotState))
			}

			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s cbs_copy_snapshot_cross_region async polling failed for snapshot %s, reason:%+v", logId, copiedId, err)
			return err
		}
	}

	return resourceTencentCloudCbsCopySnapshotCrossRegionAttachmentRead(d, meta)
}

func resourceTencentCloudCbsCopySnapshotCrossRegionAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_copy_snapshot_cross_region.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is illegal: %s", d.Id())
	}

	copiedSnapshotId := idSplit[0]
	region := idSplit[1]

	// The copied snapshot resides in the destination region, so describe it with a
	// region-specific client based on the region parsed from the resource ID.
	regionCbsClient := cbsClientWithRegion(meta, region)

	var snapshot *cbs.Snapshot
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		describeRequest := cbs.NewDescribeSnapshotsRequest()
		describeRequest.SnapshotIds = []*string{helper.String(copiedSnapshotId)}
		describeResponse, e := regionCbsClient.DescribeSnapshotsWithContext(ctx, describeRequest)
		if e != nil {
			return tccommon.RetryError(e, tccommon.InternalError)
		}

		if describeResponse != nil && describeResponse.Response != nil && len(describeResponse.Response.SnapshotSet) > 0 {
			snapshot = describeResponse.Response.SnapshotSet[0]
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read cbs_copy_snapshot_cross_region failed, reason:%+v", logId, err)
		return err
	}

	if snapshot == nil {
		log.Printf("[CRUD] tencentcloud_cbs_copy_snapshot_cross_region id=%s", d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("destination_region", region)

	if snapshot.SnapshotName != nil {
		_ = d.Set("snapshot_name", snapshot.SnapshotName)
	}

	return nil
}

func resourceTencentCloudCbsCopySnapshotCrossRegionAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_copy_snapshot_cross_region.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	// delete_bind_images is only consumed during deletion. Updating it just persists
	// the new value into state, no API call is required here.
	return resourceTencentCloudCbsCopySnapshotCrossRegionAttachmentRead(d, meta)
}

func resourceTencentCloudCbsCopySnapshotCrossRegionAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_copy_snapshot_cross_region.delete")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = cbs.NewDeleteSnapshotsRequest()
	)

	// Parse the composite ID (copied_snapshot_id + FILED_SP + destination_region).
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("cbs_copy_snapshot_cross_region id is illegal: %s", d.Id())
	}

	copiedSnapshotId := idSplit[0]
	region := idSplit[1]

	snapshotIds := make([]*string, 0, 1)
	if copiedSnapshotId != "" {
		snapshotIds = append(snapshotIds, helper.String(copiedSnapshotId))
	}

	if len(snapshotIds) == 0 {
		log.Printf("[WARN]%s tencentcloud_cbs_copy_snapshot_cross_region no copied snapshot IDs found, skip delete", logId)
		return nil
	}

	request.SnapshotIds = snapshotIds

	if v, ok := d.GetOk("delete_bind_images"); ok {
		request.DeleteBindImages = helper.Bool(v.(bool))
	} else {
		request.DeleteBindImages = helper.Bool(false)
	}

	// The copied snapshots reside in the destination region, so delete them with a
	// region-specific client based on the region parsed from the resource ID.
	regionCbsClient := cbsClientWithRegion(meta, region)

	reqErr := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		result, e := regionCbsClient.DeleteSnapshotsWithContext(ctx, request)
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

	// DeleteSnapshots is asynchronous. Poll DescribeSnapshots until SnapshotSet is an
	// empty list, which indicates the snapshot has been fully deleted.
	pollErr := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		describeRequest := cbs.NewDescribeSnapshotsRequest()
		describeRequest.SnapshotIds = []*string{helper.String(copiedSnapshotId)}
		describeResponse, e := regionCbsClient.DescribeSnapshotsWithContext(ctx, describeRequest)
		if e != nil {
			return tccommon.RetryError(e, tccommon.InternalError)
		}

		if describeResponse == nil || describeResponse.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("DescribeSnapshots return nil Response"))
		}

		if len(describeResponse.Response.SnapshotSet) > 0 {
			return resource.RetryableError(fmt.Errorf("cbs_copy_snapshot_cross_region snapshot %s in region %s is still deleting", copiedSnapshotId, region))
		}

		return nil
	})

	if pollErr != nil {
		log.Printf("[CRITAL]%s delete cbs_copy_snapshot_cross_region polling failed, reason:%+v", logId, pollErr)
		return pollErr
	}

	return nil
}
