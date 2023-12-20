package cbs

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
)

func ResourceTencentCloudCbsSnapshotPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCbsSnapshotPolicyAttachmentCreate,
		Read:   resourceTencentCloudCbsSnapshotPolicyAttachmentRead,
		Delete: resourceTencentCloudCbsSnapshotPolicyAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"storage_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of CBS.",
			},
			"snapshot_policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of CBS snapshot policy.",
			},
		},
	}
}

func resourceTencentCloudCbsSnapshotPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_snapshot_policy_attachment.create")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	storageId := d.Get("storage_id").(string)
	policyId := d.Get("snapshot_policy_id").(string)
	cbsService := CbsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		errRet := cbsService.AttachSnapshotPolicy(ctx, storageId, policyId)
		if errRet != nil {
			return tccommon.RetryError(errRet)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s cbs storage policy attach failed, reason:%s\n ", logId, err.Error())
		return err
	}

	d.SetId(storageId + tccommon.FILED_SP + policyId)
	return resourceTencentCloudCbsSnapshotPolicyAttachmentRead(d, meta)
}

func resourceTencentCloudCbsSnapshotPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_snapshot_policy_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()
	idSplit := strings.Split(id, tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("tencentcloud_cbs_snapshot_policy_attachment id is illegal: %s", id)
	}
	storageId := idSplit[0]
	policyId := idSplit[1]
	cbsService := CbsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	var policy *cbs.AutoSnapshotPolicy
	var errRet error
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		policy, errRet = cbsService.DescribeAttachedSnapshotPolicy(ctx, storageId, policyId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s cbs storage policy attach failed, reason:%s\n ", logId, err.Error())
		return err
	}
	if policy == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("storage_id", storageId)
	_ = d.Set("snapshot_policy_id", policyId)

	return nil
}

func resourceTencentCloudCbsSnapshotPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_snapshot_policy_attachment.delete")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()
	idSplit := strings.Split(id, tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("tencentcloud_cbs_snapshot_policy_attachment id is illegal: %s", id)
	}
	storageId := idSplit[0]
	policyId := idSplit[1]
	cbsService := CbsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		errRet := cbsService.UnattachSnapshotPolicy(ctx, storageId, policyId)
		if errRet != nil {
			return tccommon.RetryError(errRet)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s cbs storage policy unattach failed, reason:%s\n ", logId, err.Error())
		return err
	}

	return nil
}
