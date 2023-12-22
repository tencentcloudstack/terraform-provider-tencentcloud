package dnspod

import (
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDnspodDownloadSnapshotOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDnspodDownloadSnapshotOperationCreate,
		Read:   resourceTencentCloudDnspodDownloadSnapshotOperationRead,
		Delete: resourceTencentCloudDnspodDownloadSnapshotOperationDelete,
		Schema: map[string]*schema.Schema{
			"domain": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Domain.",
			},

			"snapshot_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Snapshot ID.",
			},

			"cos_url": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Snapshot download url.",
			},
		},
	}
}

func resourceTencentCloudDnspodDownloadSnapshotOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_download_snapshot_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = dnspod.NewDownloadSnapshotRequest()
		response   = dnspod.NewDownloadSnapshotResponse()
		domain     string
		snapshotId string
		cosUrl     string
	)
	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
		request.Domain = helper.String(v.(string))
	}

	if v, ok := d.GetOk("snapshot_id"); ok {
		snapshotId = v.(string)
		request.SnapshotId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDnsPodClient().DownloadSnapshot(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dnspod download_snapshot failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.CosUrl != nil {
		cosUrl = *response.Response.CosUrl
		_ = d.Set("cos_url", cosUrl)
	}

	d.SetId(strings.Join([]string{domain, snapshotId}, tccommon.FILED_SP))

	return resourceTencentCloudDnspodDownloadSnapshotOperationRead(d, meta)
}

func resourceTencentCloudDnspodDownloadSnapshotOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_download_snapshot_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDnspodDownloadSnapshotOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_download_snapshot_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
