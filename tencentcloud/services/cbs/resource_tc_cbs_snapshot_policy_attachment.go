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
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"storage_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"storage_ids"},
				Description:  "ID of CBS.",
			},

			"storage_ids": {
				Type:         schema.TypeSet,
				Optional:     true,
				ForceNew:     true,
				MinItems:     2,
				ExactlyOneOf: []string{"storage_id"},
				Description:  "IDs of CBS.",
				Elem:         &schema.Schema{Type: schema.TypeString},
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

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		cbsService = CbsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		storageId  string
		storageIds []string
		policyId   string
	)

	if v, ok := d.GetOk("storage_id"); ok {
		storageId = v.(string)
	}

	if v, ok := d.GetOk("storage_ids"); ok {
		for _, item := range v.(*schema.Set).List() {
			if storageId, ok := item.(string); ok && storageId != "" {
				storageIds = append(storageIds, storageId)
			}
		}
	}

	if v, ok := d.GetOk("snapshot_policy_id"); ok {
		policyId = v.(string)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		errRet := cbsService.AttachSnapshotPolicy(ctx, storageId, storageIds, policyId)
		if errRet != nil {
			return tccommon.RetryError(errRet)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s cbs storage policy attach failed, reason:%s\n ", logId, err.Error())
		return err
	}

	if storageId != "" {
		d.SetId(strings.Join([]string{storageId, policyId}, tccommon.FILED_SP))
	} else {
		storageIdsStr := strings.Join(storageIds, tccommon.COMMA_SP)
		d.SetId(strings.Join([]string{storageIdsStr, policyId}, tccommon.FILED_SP))
	}

	return resourceTencentCloudCbsSnapshotPolicyAttachmentRead(d, meta)
}

func resourceTencentCloudCbsSnapshotPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_snapshot_policy_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		cbsService = CbsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		id         = d.Id()
		storageId  string
		storageIds []string
	)

	idSplit := strings.Split(id, tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("tencentcloud_cbs_snapshot_policy_attachment id is illegal: %s", id)
	}

	storageIdObj := idSplit[0]
	policyId := idSplit[1]

	storageIdObjSplit := strings.Split(storageIdObj, tccommon.COMMA_SP)
	if len(storageIdObjSplit) == 1 {
		storageId = storageIdObjSplit[0]
	} else {
		storageIds = storageIdObjSplit
	}

	var (
		policy *cbs.AutoSnapshotPolicy
		errRet error
	)

	if storageId != "" {
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
	}

	if len(storageIds) > 0 {
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			policy, errRet = cbsService.DescribeAttachedSnapshotPolicyDisksById(ctx, policyId)
			if errRet != nil {
				return tccommon.RetryError(errRet, tccommon.InternalError)
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s cbs storage policy attach failed, reason:%s\n ", logId, err.Error())
			return err
		}

		if policy == nil || policy.DiskIdSet == nil || len(policy.DiskIdSet) < 1 {
			d.SetId("")
			return nil
		}

		tmpList := GetDiskIds(storageIds, policy.DiskIdSet)
		_ = d.Set("storage_ids", tmpList)
	}

	_ = d.Set("snapshot_policy_id", policyId)

	return nil
}

func resourceTencentCloudCbsSnapshotPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_snapshot_policy_attachment.delete")()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		cbsService = CbsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		id         = d.Id()
		storageId  string
		storageIds []string
	)

	idSplit := strings.Split(id, tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("tencentcloud_cbs_snapshot_policy_attachment id is illegal: %s", id)
	}

	storageIdObj := idSplit[0]
	policyId := idSplit[1]

	storageIdObjSplit := strings.Split(storageIdObj, tccommon.COMMA_SP)
	if len(storageIdObjSplit) == 1 {
		storageId = storageIdObjSplit[0]
	} else {
		storageIds = storageIdObjSplit
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		errRet := cbsService.UnattachSnapshotPolicy(ctx, storageId, storageIds, policyId)
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

func GetDiskIds(A []string, B []*string) []string {
	set := make(map[string]bool, len(B))
	for _, ptr := range B {
		if ptr != nil {
			set[*ptr] = true
		}
	}

	var tmpList []string
	for _, s := range A {
		if set[s] {
			tmpList = append(tmpList, s)
		}
	}

	return tmpList
}
