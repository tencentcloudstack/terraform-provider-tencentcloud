package cfs

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfs/v20190719"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCfsAutoSnapshotPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCfsAutoSnapshotPolicyAttachmentCreate,
		Read:   resourceTencentCloudCfsAutoSnapshotPolicyAttachmentRead,
		Delete: resourceTencentCloudCfsAutoSnapshotPolicyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"auto_snapshot_policy_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "ID of the snapshot to be unbound.",
			},

			"file_system_ids": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "List of IDs of the file systems to be unbound, separated by comma.",
			},
		},
	}
}

func resourceTencentCloudCfsAutoSnapshotPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfs_auto_snapshot_policy_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request              = cfs.NewBindAutoSnapshotPolicyRequest()
		autoSnapshotPolicyId string
		fileSystemIds        string
	)
	if v, ok := d.GetOk("auto_snapshot_policy_id"); ok {
		autoSnapshotPolicyId = v.(string)
		request.AutoSnapshotPolicyId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("file_system_ids"); ok {
		fileSystemIds = v.(string)
		request.FileSystemIds = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfsClient().BindAutoSnapshotPolicy(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cfs autoSnapshotPolicyAttachment failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(autoSnapshotPolicyId + tccommon.FILED_SP + fileSystemIds)

	return resourceTencentCloudCfsAutoSnapshotPolicyAttachmentRead(d, meta)
}

func resourceTencentCloudCfsAutoSnapshotPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfs_auto_snapshot_policy_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CfsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	autoSnapshotPolicyId := idSplit[0]
	fileSystemIds := idSplit[1]

	autoSnapshotPolicyAttachment, err := service.DescribeCfsAutoSnapshotPolicyAttachmentById(ctx, autoSnapshotPolicyId, fileSystemIds)
	if err != nil {
		return err
	}

	if autoSnapshotPolicyAttachment == nil {
		d.SetId("")
		return fmt.Errorf("resource `tencentcloud_cfs_auto_snapshot_policy_attachment` %s does not exist", d.Id())
	}

	if autoSnapshotPolicyAttachment.AutoSnapshotPolicyId != nil {
		_ = d.Set("auto_snapshot_policy_id", autoSnapshotPolicyId)
	}

	if autoSnapshotPolicyAttachment.FileSystems != nil {
		_ = d.Set("file_system_ids", fileSystemIds)
	}

	return nil
}

func resourceTencentCloudCfsAutoSnapshotPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfs_auto_snapshot_policy_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CfsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	autoSnapshotPolicyId := idSplit[0]
	fileSystemIds := idSplit[1]

	if err := service.DeleteCfsAutoSnapshotPolicyAttachmentById(ctx, autoSnapshotPolicyId, fileSystemIds); err != nil {
		return err
	}

	return nil
}
